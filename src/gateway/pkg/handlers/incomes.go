package handlers

import (
	"context"
	"gateway/pkg/models/income"
	"gateway/pkg/myjson"
	"gateway/pkg/services"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type IncomesHandler interface {
	List(w http.ResponseWriter, r *http.Request)
	Show(w http.ResponseWriter, r *http.Request)
	Add(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type IncomeService interface {
	Get(ctx context.Context, id int) (income.Income, error)
	Query(ctx context.Context, offset, limit int) ([]income.Income, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, input *income.IncomeCreationRequest) (income.Income, error)
	Update(ctx context.Context, id int, input *income.IncomeCreationRequest) (income.Income, error)
	Delete(ctx context.Context, id int) (income.Income, error)

	CreateFromRequest(ctx context.Context, r *http.Request) (income.Income, error)
	UpdateFromRequest(ctx context.Context, r *http.Request, id int) (income.Income, error)
	DeleteFromRequest(ctx context.Context, r *http.Request, id int) (income.Income, error)
}

type IncomeMainHandler struct {
	Logger  *zap.SugaredLogger
	Service IncomeService
}

func NewIncomeHandler(logger *zap.SugaredLogger) (h *IncomeMainHandler) {
	client := &http.Client{}

	ctrl := services.NewIncomeService(client, logger)
	h = &IncomeMainHandler{
		Logger:  logger,
		Service: ctrl,
	}
	return h
}

func (h IncomeMainHandler) List(w http.ResponseWriter, r *http.Request) {
	// lol := ps.ByName("id")
	elems, err := h.Service.Query(context.WithValue(r.Context(), "Authorization", r.Header.Get("Authorization")), 0, 64)
	if err != nil {
		myjson.JSONError(w, http.StatusInternalServerError, "DB error: "+err.Error())
		return
	}

	myjson.JSONResponce(w, http.StatusOK, elems)
}

func (h IncomeMainHandler) Show(w http.ResponseWriter, r *http.Request) {
	ps := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		myjson.JSONResponce(w, http.StatusBadRequest, errors.Wrap(err, "bad ID in URL"))
		return
	}

	income, err := h.Service.Get(context.WithValue(r.Context(), "Authorization", r.Header.Get("Authorization")), id)
	if err != nil {
		myjson.JSONResponce(w, http.StatusInternalServerError, errors.Wrap(err, ""))
		return
	}

	myjson.JSONResponce(w, http.StatusOK, income)
}

func (h *IncomeMainHandler) Add(w http.ResponseWriter, r *http.Request) {
	res, err := h.Service.CreateFromRequest(r.Context(), r)
	if err != nil {
		myjson.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	myjson.JSONResponce(w, http.StatusCreated, res)
}

func (h IncomeMainHandler) Update(w http.ResponseWriter, r *http.Request) {
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

func (h IncomeMainHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ps := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		myjson.JSONResponce(w, http.StatusBadRequest, errors.Wrap(err, "bad ID in URL"))
		return
	}

	income, _ := h.Service.DeleteFromRequest(r.Context(), r, id)

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

	myjson.JSONResponce(w, http.StatusNoContent, income)
}
