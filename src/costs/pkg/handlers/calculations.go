package handlers

import (
	"costs/pkg/myjson"
	"costs/pkg/services"
	"net/http"

	"go.uber.org/zap"
)

type CalculationHandler interface {
	TotalBalance(w http.ResponseWriter, r *http.Request)
}

type CalcMainHandler struct {
	Logger        *zap.SugaredLogger
	CostsSource   services.CostService
	IncomesSource services.IncomeService
}

func (h CalcMainHandler) TotalBalance(w http.ResponseWriter, r *http.Request) {
	// lol := ps.ByName("id")
	elems_minus, err := h.CostsSource.Query(r.Context(), 0, 64)
	if err != nil {
		myjson.JSONError(w, http.StatusInternalServerError, "DB error: "+err.Error())
		return
	}

	elems_plus, err := h.IncomesSource.Query(r.Context(), 0, 64)
	if err != nil {
		myjson.JSONError(w, http.StatusInternalServerError, "DB error: "+err.Error())
		return
	}

	var sum float32 = 0.0
	for _, v := range elems_plus {
		sum += v.Amount
	}
	for _, v := range elems_minus {
		sum -= v.Amount
	}

	myjson.JSONResponce(w, http.StatusOK, sum)
}
