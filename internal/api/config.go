package api

import (
	"net/http"

	"git.friedmicro-lab.org/friedmicro/ai-planner/internal/model"
)

func (s *Server) getConfig(w http.ResponseWriter, r *http.Request) {
	c, err := s.store.GetConfig(r.Context())
	if err != nil {
		s.internalErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (s *Server) putConfig(w http.ResponseWriter, r *http.Request) {
	var c model.Config
	if err := decode(r, &c); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	c.ID = 1
	if err := s.store.SaveConfig(r.Context(), c); err != nil {
		s.internalErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, c)
}
