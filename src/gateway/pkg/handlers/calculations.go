package handlers

import (
	"fmt"
	"gateway/pkg/myjson"
	"gateway/pkg/utils"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type CalculationHandler interface {
	TotalBalance(w http.ResponseWriter, r *http.Request)
}

type CalcMainHandler struct {
	Logger *zap.SugaredLogger
	Client *http.Client
}

func NewCalcHandler(logger *zap.SugaredLogger) (h *CalcMainHandler) {
	client := &http.Client{}

	h = &CalcMainHandler{
		Logger: logger,
		Client: client,
	}
	return h
}

func (h CalcMainHandler) TotalBalance(w http.ResponseWriter, r *http.Request) {
	requestURL := fmt.Sprintf("%s/api/v1/balance", utils.Config.CostsEndpoint)

	URL, _ := url.Parse(requestURL)
	r.URL.Scheme = URL.Scheme
	r.URL.Host = URL.Host
	r.URL.Path = URL.Path
	r.RequestURI = ""

	res, err := h.Client.Do(r)
	if err != nil {
		myjson.JSONError(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		h.Logger.Errorln(err.Error())
	}
	res.Body.Close()

	var result float32
	if err = myjson.From(body, &result); err != nil {
		myjson.JSONError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to decode response").Error())
	}

	myjson.JSONResponce(w, http.StatusOK, result)
}
