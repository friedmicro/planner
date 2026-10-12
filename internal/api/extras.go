package api

import (
	"database/sql"
	"errors"
	"net/http"

	"git.friedmicro-lab.org/friedmicro/ai-planner/internal/model"
)

func (s *Server) listExtras(w http.ResponseWriter, r *http.Request) {
	v, err := s.store.ListExtras(r.Context())
	if err != nil {
		s.internalErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}

func (s *Server) addExtra(w http.ResponseWriter, r *http.Request) {
	var x model.ExtraEvent
	if err := decode(r, &x); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	id, err := s.store.AddExtra(r.Context(), x)
	if err != nil {
		s.internalErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

func (s *Server) deleteExtra(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.store.DeleteExtra(r.Context(), id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeErr(w, http.StatusNotFound, "not found")
			return
		}
		s.internalErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
