package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"users/pkg/models"
	"users/pkg/objects"
	"users/pkg/utils"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/julienschmidt/httprouter"
)

type AuthCtrl struct {
	Auth *models.AuthM
}

type Token struct {
	jwt.StandardClaims
	Role string `json:"role,omitempty"`
}

// func InitAuth(r *mux.Router, auth *models.AuthM) {
// 	ctrl := &AuthCtrl{auth}
// 	r.HandleFunc("/register", ctrl.Register).Methods("POST")
// 	r.HandleFunc("/authorize", ctrl.Authorize).Methods("POST")
// }

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

func (ctrl *AuthCtrl) Register(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	token, err := RetrieveToken(w, r)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if token.Role != utils.Admin.String() {
		Forbidden(w, fmt.Sprintf("not allowed for %s role", token.Role))
		return
	}

	req_body := new(objects.UserCreateRequest)
	err = json.NewDecoder(r.Body).Decode(req_body)
	if err != nil {
		fmt.Println(err.Error())
		ValidationErrorResponse(w, err.Error())
		return
	}

	err = ctrl.Auth.Create(req_body)
	if err != nil {
		fmt.Println(err.Error())
		BadRequest(w, "user creation failed")
	} else {
		JsonSuccess(w, nil)
	}
}

func (ctrl *AuthCtrl) Authorize(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req_body := new(objects.AuthRequest)
	err := json.NewDecoder(r.Body).Decode(req_body)
	if err != nil {
		fmt.Println(err.Error())
		ValidationErrorResponse(w, err.Error())
		return
	}

	data, err := ctrl.Auth.Auth(req_body.Username, req_body.Password)
	if err != nil {
		fmt.Println(err.Error())
		BadRequest(w, "auth failed")
	} else {
		JsonSuccess(w, data)
	}
}

type AuthM struct {
	client *http.Client
}

func NewAuthM(client *http.Client) *AuthM {
	return &AuthM{client: client}
}

type Models struct {
	Auth *AuthM
}

func InitModels() *Models {
	models := new(Models)
	client := &http.Client{}

	models.Auth = NewAuthM(client)
	return models
}

func (model *AuthM) Create(request *objects.UserCreateRequest) error {
	request.GroupIds = []string{utils.Config.Okta.ClientGroup}
	request.Profile.UserType = utils.User.String()

	req_body, err := json.Marshal(request)
	if err != nil {
		return err
	}
	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/api/v1/users/?activate=true", utils.Config.Okta.Endpoint),
		bytes.NewBuffer(req_body),
	)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("SSWS %s", utils.Config.Okta.SSWSToken))

	resp, err := model.client.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth failed with code %d", resp.StatusCode)
	} else {
		return nil
	}
}

func (model *AuthM) Auth(username string, password string) (*objects.AuthResponse, error) {
	authRequest := url.Values{}
	authRequest.Set("scope", "openid")
	authRequest.Set("grant_type", "password")
	authRequest.Set("username", username)
	authRequest.Set("password", password)
	authRequest.Set("client_id", utils.Config.Okta.ClientId)
	authRequest.Set("client_secret", utils.Config.Okta.ClientSecret)
	encodedData := authRequest.Encode()

	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/oauth2/default/v1/token", utils.Config.Okta.Endpoint),
		strings.NewReader(encodedData),
	)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(authRequest.Encode())))

	resp, err := model.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data := &objects.AuthResponse{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth failed, code: %d", resp.StatusCode)
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
