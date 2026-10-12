package api

import (
	"database/sql"
	"errors"
	"net/http"

	"git.friedmicro-lab.org/friedmicro/ai-planner/internal/model"
)

func (s *Server) listEvents(w http.ResponseWriter, r *http.Request) {
	v, err := s.store.ListEvents(r.Context())
	if err != nil {
		s.internalErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}

func (s *Server) addEvent(w http.ResponseWriter, r *http.Request) {
	var e model.Event
	if err := decode(r, &e); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	id, err := s.store.AddEvent(r.Context(), e)
	if err != nil {
		s.internalErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

func (s *Server) updateEvent(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	var e model.Event
	if err := decode(r, &e); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	e.ID = int(id)
	if err := s.store.UpdateEvent(r.Context(), e); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeErr(w, http.StatusNotFound, "not found")
			return
		}
		s.internalErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, e)
}

func (s *Server) deleteEvent(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.store.DeleteEvent(r.Context(), id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeErr(w, http.StatusNotFound, "not found")
			return
		}
		s.internalErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
