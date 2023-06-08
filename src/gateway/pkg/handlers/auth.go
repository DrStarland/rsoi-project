package handlers

import (
	"encoding/json"
	"fmt"
	"gateway/pkg/myjson"
	"gateway/pkg/objects"
	"gateway/pkg/services"
	"gateway/pkg/utils"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type AuthHandler struct {
	Logger *zap.SugaredLogger
	*services.AuthController
}

func NewAuthHandler(logger *zap.SugaredLogger) (h *AuthHandler) {
	client := &http.Client{}

	auth_var := services.NewAuthController(client, logger)
	h = &AuthHandler{
		Logger:         logger,
		AuthController: auth_var,
	}
	return h
}

func (ctrl *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/api/v1/register", utils.Config.UsersEndpoint),
		r.Body,
	)
	req.Header.Add("Authorization", r.Header.Get("Authorization"))

	resp, err := ctrl.Client.Do(req)
	if err != nil {
		ctrl.Logger.Errorln(errors.Wrap(err, "handler"))
		myjson.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer resp.Body.Close()

	myjson.JSONResponce(w, resp.StatusCode, resp)
}

func (ctrl *AuthHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/authorize", utils.Config.UsersEndpoint), r.Body)

	resp, err := ctrl.Client.Do(req)
	if err != nil {
		ctrl.Logger.Errorln(errors.Wrap(err, "!!!"))
		myjson.JSONError(w, http.StatusBadRequest, errors.Wrap(err, "bad credentials").Error())
		return
	}
	if resp.StatusCode == http.StatusOK {
		data := &objects.AuthResponse{}
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, data)
		myjson.JSONResponce(w, http.StatusOK, data)
	} else {
		ctrl.Logger.Errorln(errors.Wrap(err, "after auth method"))
		myjson.JSONError(w, http.StatusBadRequest, "Authorization failed: "+err.Error())
	}
}
