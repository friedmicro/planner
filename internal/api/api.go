package api

import (
	"log"
	"net/http"

	"git.friedmicro-lab.org/friedmicro/ai-planner/internal/store"
)

type Server struct {
	store *store.Store
	log   *log.Logger
}

func New(st *store.Store, l *log.Logger) *Server {
	if l == nil {
		l = log.Default()
	}
	return &Server{store: st, log: l}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", s.health)

	mux.HandleFunc("GET /api/config", s.getConfig)
	mux.HandleFunc("PUT /api/config", s.putConfig)

	mux.HandleFunc("GET /api/sleep", s.listSleep)
	mux.HandleFunc("POST /api/sleep", s.addSleep)
	mux.HandleFunc("DELETE /api/sleep/{id}", s.deleteSleep)

	mux.HandleFunc("GET /api/blacklist", s.listBlacklist)
	mux.HandleFunc("POST /api/blacklist", s.addBlacklist)
	mux.HandleFunc("DELETE /api/blacklist/{id}", s.deleteBlacklist)

	mux.HandleFunc("GET /api/events", s.listEvents)
	mux.HandleFunc("POST /api/events", s.addEvent)
	mux.HandleFunc("PUT /api/events/{id}", s.updateEvent)
	mux.HandleFunc("DELETE /api/events/{id}", s.deleteEvent)

	mux.HandleFunc("GET /api/extras", s.listExtras)
	mux.HandleFunc("POST /api/extras", s.addExtra)
	mux.HandleFunc("DELETE /api/extras/{id}", s.deleteExtra)

	mux.HandleFunc("POST /api/solve", s.solve)
	mux.HandleFunc("GET /api/schedule", s.getSchedule)

	return cors(logRequests(s.log, mux))
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Access-Control-Allow-Origin", "*")
		h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		h.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		h.Set("Access-Control-Max-Age", "86400")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func logRequests(l *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		l.Printf("%s %s -> %d", r.Method, r.URL.Path, rec.status)
	})
}
