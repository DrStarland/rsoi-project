package handlers

import (
	"context"
	"gateway/pkg/models/cost"
	"gateway/pkg/myjson"
	"gateway/pkg/services"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type CostsHandler interface {
	List(w http.ResponseWriter, r *http.Request)
	Show(w http.ResponseWriter, r *http.Request)
	Add(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type CostService interface {
	Get(ctx context.Context, id int) (cost.Cost, error)
	Query(ctx context.Context, offset, limit int) ([]cost.Cost, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, input *cost.CostCreationRequest) (cost.Cost, error)
	Update(ctx context.Context, id int, input *cost.CostCreationRequest) (cost.Cost, error)
	Delete(ctx context.Context, id int) (cost.Cost, error)

	CreateFromRequest(ctx context.Context, r *http.Request) (cost.Cost, error)
	UpdateFromRequest(ctx context.Context, r *http.Request, id int) (cost.Cost, error)
	DeleteFromRequest(ctx context.Context, r *http.Request, id int) (cost.Cost, error)
}

type CostMainHandler struct {
	Logger  *zap.SugaredLogger
	Service CostService
}

func NewCostHandler(logger *zap.SugaredLogger) (h *CostMainHandler) {
	client := &http.Client{}

	ctrl := services.NewCostService(client, logger)
	h = &CostMainHandler{
		Logger:  logger,
		Service: ctrl,
	}
	return h
}

func (h CostMainHandler) List(w http.ResponseWriter, r *http.Request) {
	// lol := ps.ByName("id")
	elems, err := h.Service.Query(context.WithValue(r.Context(), "Authorization", r.Header.Get("Authorization")), 0, 64)
	if err != nil {
		myjson.JSONError(w, http.StatusInternalServerError, "DB error: "+err.Error())
		return
	}

	myjson.JSONResponce(w, http.StatusOK, elems)
}

func (h CostMainHandler) Show(w http.ResponseWriter, r *http.Request) {
	ps := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		myjson.JSONResponce(w, http.StatusBadRequest, errors.Wrap(err, "bad ID in URL"))
		return
	}

	cost, err := h.Service.Get(context.WithValue(r.Context(), "Authorization", r.Header.Get("Authorization")), id)
	if err != nil {
		myjson.JSONResponce(w, http.StatusInternalServerError, errors.Wrap(err, ""))
		return
	}

	myjson.JSONResponce(w, http.StatusOK, cost)
}

func (h *CostMainHandler) Add(w http.ResponseWriter, r *http.Request) {
	res, err := h.Service.CreateFromRequest(r.Context(), r)
	if err != nil {
		myjson.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	myjson.JSONResponce(w, http.StatusCreated, res)
}

func (h CostMainHandler) Update(w http.ResponseWriter, r *http.Request) {
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

func (h CostMainHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ps := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		myjson.JSONResponce(w, http.StatusBadRequest, errors.Wrap(err, "bad ID in URL"))
		return
	}

	cost, err := h.Service.DeleteFromRequest(r.Context(), r, id)

	// if err != nil {
	// 	switch err.Error() {
	// 	case "not exist":
	// 		myjson.JSONResponce(w, http.StatusNoContent, errors.Wrap(err, ""))
	// 		return
	// 	default:
	// 		myjson.JSONResponce(w, http.StatusInternalServerError, errors.Wrap(err, ""))
	// 	}
	// 	return
	// }

	myjson.JSONResponce(w, http.StatusNoContent, cost)
}
