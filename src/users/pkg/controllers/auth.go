package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"users/pkg/utils"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
)

// type auhtCtrl struct {
// 	auth *models.AuthM
// }

// func InitAuth(r *mux.Router, auth *models.AuthM) {
// 	ctrl := &auhtCtrl{auth}
// 	r.HandleFunc("/register", ctrl.register).Methods("POST")
// 	r.HandleFunc("/authorize", ctrl.authorize).Methods("POST")
// }

type Token struct {
	jwt.StandardClaims
	Role string `json:"role,omitempty"`
}

func newJWKs(rawJWKS string) *keyfunc.JWKS {
	jwksJSON := json.RawMessage(rawJWKS)
	jwks, err := keyfunc.NewJSON(jwksJSON)
	if err != nil {
		panic(err)
	}
	return jwks
}

func RetrieveToken(w http.ResponseWriter, r *http.Request) (*Token, error) {
	reqToken := r.Header.Get("Authorization")
	if len(reqToken) == 0 {
		TokenIsMissing(w)
		return nil, fmt.Errorf("TokenIsMissing")
	}
	splitToken := strings.Split(reqToken, "Bearer ")
	tokenStr := splitToken[1]
	jwks := newJWKs(utils.Config.RawJWKS)
	tk := &Token{}

	token, err := jwt.ParseWithClaims(tokenStr, tk, jwks.Keyfunc)
	if err != nil || !token.Valid {
		JwtAccessDenied(w)
		return nil, fmt.Errorf("JwtAccessDenied")
	}
	if time.Now().Unix()-tk.ExpiresAt > 0 {
		TokenExpired(w)
		return nil, fmt.Errorf("TokenExpired")
	}

	return tk, nil
}

// func (ctrl *auhtCtrl) register(w http.ResponseWriter, r *http.Request) {
// 	token, err := RetrieveToken(w, r)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}

// 	if token.Role != utils.Admin.String() {
// 		Forbidden(w, fmt.Sprintf("not allowed for %s role", token.Role))
// 		return
// 	}

// 	req_body := new(objects.UserCreateRequest)
// 	err = json.NewDecoder(r.Body).Decode(req_body)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		ValidationErrorResponse(w, err.Error())
// 		return
// 	}

// 	err = ctrl.auth.Create(req_body)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		BadRequest(w, "user creation failed")
// 	} else {
// 		JsonSuccess(w, nil)
// 	}
// }

// func (ctrl *auhtCtrl) authorize(w http.ResponseWriter, r *http.Request) {
// 	req_body := new(objects.AuthRequest)
// 	err := json.NewDecoder(r.Body).Decode(req_body)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		ValidationErrorResponse(w, err.Error())
// 		return
// 	}

// 	data, err := ctrl.auth.Auth(req_body.Username, req_body.Password)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		BadRequest(w, "auth failed")
// 	} else {
// 		JsonSuccess(w, data)
// 	}
// }
