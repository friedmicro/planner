package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"git.friedmicro-lab.org/friedmicro/ai-planner/internal/model"
)

func (s *Server) listSleep(w http.ResponseWriter, r *http.Request) {
	v, err := s.store.ListSleep(r.Context())
	if err != nil {
		s.internalErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}

func (s *Server) addSleep(w http.ResponseWriter, r *http.Request) {
	var b model.BlockedInterval
	if err := decode(r, &b); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	id, err := s.store.AddSleep(r.Context(), b)
	if err != nil {
		s.internalErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

func (s *Server) deleteSleep(w http.ResponseWriter, r *http.Request) {
	s.deleteInterval(w, r, s.store.DeleteSleep)
}

func (s *Server) listBlacklist(w http.ResponseWriter, r *http.Request) {
	v, err := s.store.ListBlacklist(r.Context())
	if err != nil {
		s.internalErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}

func (s *Server) addBlacklist(w http.ResponseWriter, r *http.Request) {
	var b model.BlockedInterval
	if err := decode(r, &b); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	id, err := s.store.AddBlacklist(r.Context(), b)
	if err != nil {
		s.internalErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

func (s *Server) deleteBlacklist(w http.ResponseWriter, r *http.Request) {
	s.deleteInterval(w, r, s.store.DeleteBlacklist)
}

func (s *Server) deleteInterval(w http.ResponseWriter, r *http.Request, del func(context.Context, int64) error) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := del(r.Context(), id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeErr(w, http.StatusNotFound, "not found")
			return
		}
		s.internalErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
