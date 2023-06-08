package handlers

import (
	"context"
	"gateway/pkg/myjson"
	"gateway/pkg/services"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type NotesHandler interface {
	List(w http.ResponseWriter, r *http.Request)
	Show(w http.ResponseWriter, r *http.Request)
	Add(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type NoteMainHandler struct {
	Logger  *zap.SugaredLogger
	Service services.NoteService
}

func NewNoteHandler(logger *zap.SugaredLogger) (h *NoteMainHandler) {
	client := &http.Client{}

	ctrl := services.NewNoteService(client, logger)
	h = &NoteMainHandler{
		Logger:  logger,
		Service: ctrl,
	}
	return h
}

func (h NoteMainHandler) List(w http.ResponseWriter, r *http.Request) {
	// lol := ps.ByName("id")
	elems, err := h.Service.Query(context.WithValue(r.Context(), "Authorization", r.Header.Get("Authorization")), 0, 64)
	if err != nil {
		myjson.JSONError(w, http.StatusInternalServerError, "DB error: "+err.Error())
		return
	}

	myjson.JSONResponce(w, http.StatusOK, elems)
}

func (h NoteMainHandler) Show(w http.ResponseWriter, r *http.Request) {
	ps := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		myjson.JSONResponce(w, http.StatusBadRequest, errors.Wrap(err, "bad ID in URL"))
		return
	}

	note, err := h.Service.Get(context.WithValue(r.Context(), "Authorization", r.Header.Get("Authorization")), id)
	if err != nil {
		myjson.JSONResponce(w, http.StatusInternalServerError, errors.Wrap(err, ""))
		return
	}

	myjson.JSONResponce(w, http.StatusOK, note)
}

func (h *NoteMainHandler) Add(w http.ResponseWriter, r *http.Request) {
	res, err := h.Service.CreateFromRequest(r.Context(), r)
	if err != nil {
		myjson.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	myjson.JSONResponce(w, http.StatusCreated, res)
}

func (h NoteMainHandler) Update(w http.ResponseWriter, r *http.Request) {
	ps := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		myjson.JSONResponce(w, http.StatusBadRequest, errors.Wrap(err, "bad ID in URL"))
	}

	res, err := h.Service.UpdateFromRequest(r.Context(), r, id)
	if err != nil {
		myjson.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	myjson.JSONResponce(w, http.StatusCreated, res)
}

func (h NoteMainHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ps := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		myjson.JSONResponce(w, http.StatusBadRequest, errors.Wrap(err, "bad ID in URL"))
		return
	}

	note, err := h.Service.DeleteFromRequest(r.Context(), r, id)

	if err != nil {
		switch err.Error() {
		case "not exist":
			myjson.JSONResponce(w, http.StatusNoContent, errors.Wrap(err, ""))
			return
		default:
			myjson.JSONResponce(w, http.StatusInternalServerError, errors.Wrap(err, ""))
		}
		return
	}

	myjson.JSONResponce(w, http.StatusNoContent, note)
}
