package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"users/pkg/models/authorization"
	"users/pkg/myjson"
	"users/pkg/services"
	"users/pkg/utils"

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
	token, err := services.RetrieveToken(w, r)
	if err != nil {
		ctrl.Logger.Errorln(errors.Wrap(err, "problem with token"))
		return
	}

	if token.Role != utils.Admin.String() {
		err_text := fmt.Sprintf("not allowed for %s role", token.Role)
		ctrl.Logger.Errorln(errors.Wrap(nil, err_text))
		myjson.JSONError(w, http.StatusForbidden, err_text)
		return
	}

	req_body := new(authorization.UserCreateRequest)
	err = json.NewDecoder(r.Body).Decode(req_body)
	if err != nil {
		ctrl.Logger.Errorln(errors.Wrap(err, ""))
		myjson.JSONError(w, http.StatusBadRequest, "validation error: "+err.Error())
		return
	}

	err = ctrl.Create(req_body)
	if err != nil {
		ctrl.Logger.Errorln(errors.Wrap(err, ""))
		myjson.JSONError(w, http.StatusBadRequest, "user creation failed: "+err.Error())
	} else {
		myjson.JSONResponce(w, http.StatusOK, nil)
	}
}

func (ctrl *AuthHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	req_body := new(authorization.AuthRequest)
	err := json.NewDecoder(r.Body).Decode(req_body)
	if err != nil {
		ctrl.Logger.Errorln(errors.Wrap(err, "!!!"))
		myjson.JSONError(w, http.StatusBadRequest, errors.Wrap(err, "bad credentials").Error())
		return
	}

	data, err := ctrl.Auth(req_body.Username, req_body.Password)
	if err != nil {
		ctrl.Logger.Errorln(errors.Wrap(err, "after auth method"))
		myjson.JSONError(w, http.StatusBadRequest, "Authorization failed: "+err.Error())
	} else {
		myjson.JSONResponce(w, http.StatusOK, data)
	}
}
