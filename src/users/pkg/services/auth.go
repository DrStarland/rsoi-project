package services

// curl -X POST "http://127.0.0.1:8080/api/v1/register" -H  "accept: application/json" -H  "Content-Type: application/json" -H "Authorization: Bearer eyJraWQiOiJvRDdxMkQzLTExdEVGUWdaWGZvaWtqSFZtamNVRVBVLWlOR2lyR2FkTlVvIiwiYWxnIjoiUlMyNTYifQ.eyJzdWIiOiIwMHU3djh1aXNlM2tQbWF5ejVkNyIsInZlciI6MSwiaXNzIjoiaHR0cHM6Ly9kZXYtOTg1NDExNDIub2t0YS5jb20vb2F1dGgyL2RlZmF1bHQiLCJhdWQiOiIwb2E3djhyYWlyT1ViWUF2eTVkNyIsImlhdCI6MTY4MzIwNjc4MywiZXhwIjoxNjgzMjEwMzgzLCJqdGkiOiJJRC43Vy11ZzVVTEtYa1F4LXJyUnc4aDFQZjhGcmJZRWg3RmZqUFVsNDg0eGx3IiwiYW1yIjpbInB3ZCJdLCJpZHAiOiIwMG83djFuM2czdkphZTZLYTVkNyIsImF1dGhfdGltZSI6MTY4MzIwNjc4MywiYXRfaGFzaCI6IkFNWmIyRGx1QTdpSXZXSUpWRzMtR0EifQ.PXDWZzvie7IZi0iyFCkfH7nkW0HvGllT20kqeWcFnojC7sbE0hwSuPxcZ3rPFS5aZRd_duL7e3axv4CGRsjru6shohC-ldRhDFwa9Nfx4w8IX0bXyHoG94tG1_-B-P9keEmZMKPMzxEVrXqCoApN6dsK68gzWLfNHCQ44zk-IeWQXN431buTXVTO77WkIMYjxqvKqbdve1c8WeVL9aWCYM6j4BK7hIF_yg1DkmO6qx_FF7wSRDBo1l_0DGgsslgjm3MHKBx9aL3XuiQjIOEmZG30OHxkzk1X1jLy3XSYYQ5qLmBSkwxs_nQ1kVgEGQV71SufFj7pRJkxV9F0dC9b-w" -d "{  \"credentials\": {    \"password\": {      \"value\": \"228lohNemamont\$\"    }  },  \"groupIds\": [    \"string\"  ],  \"profile\": {    \"email\": \"string\",    \"firstName\": \"string\",    \"lastName\": \"string\",   \"login\": \"ojmutmhlacsxrhkgds@tmmcv.com\",    \"mobilePhone\": \"string\",    \"userType\": \"string\"  }}"

// {
// 	"password": "228lohNemamont$",
// 	"username": "ojmutmhlacsxrhkgds@tmmcv.com"
// }

// "values": [
// 		{
// 			"key": "serviceUrl",
// 			"value": "https://load-balancer-drstarland.cloud.okteto.net",
// 			"enabled": true
// 		},
// 		{
// 			"key": "identityProviderUrl",
// 			"value": "https://dev-98541142.okta.com",
// 			"enabled": true
// 		},
// 		{
// 			"key": "username",
// 			"value": "ojmutmhlacsxrhkgds@tmmcv.com",
// 			"enabled": true
// 		},
// 		{
// 			"key": "password",
// 			"value": "228lohNemamont$",
// 			"enabled": true
// 		},
// 		{
// 			"key": "clientId",
// 			"value": "0oa7v8rairOUbYAvy5d7",
// 			"enabled": true
// 		},
// 		{
// 			"key": "clientSecret",
// 			"value": "iQcihL2DDY6AXyYG3_0XuPsFdWQ9w9vk98xFKBIR",
// 			"enabled": true
// 		}
// 	],

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
	"users/pkg/models/authorization"
	"users/pkg/myjson"
	"users/pkg/utils"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

type AuthController struct {
	client *http.Client
	logger *zap.SugaredLogger
}

func NewAuthController(client *http.Client, logger *zap.SugaredLogger) *AuthController {
	return &AuthController{client: client, logger: logger}
}

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
		myjson.JSONError(w, http.StatusUnauthorized, "Missing auth token")
		return nil, fmt.Errorf("TokenIsMissed")
	}

	tokenStr := strings.Split(reqToken, "Bearer ")[1]
	jwks := newJWKs(utils.Config.RawJWKS)
	tk := &Token{}

	token, err := jwt.ParseWithClaims(tokenStr, tk, jwks.Keyfunc)
	if err != nil || !token.Valid {
		myjson.JSONError(w, http.StatusUnauthorized, "jwt-token is not valid")
		return nil, fmt.Errorf("JwtAccessDenied")
	}

	// проверка времени существования токена
	if time.Now().Unix()-tk.ExpiresAt > 0 {
		myjson.JSONError(w, http.StatusUnauthorized, "jwt-token expired")
		return nil, fmt.Errorf("token expired")
	}

	return tk, nil
}

func (model *AuthController) Create(request *authorization.UserCreateRequest) error {
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
	req.Header.Add("Authorization", fmt.Sprintf("SSWS %s", utils.Config.Okta.OktetoToken))

	resp, err := model.client.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth failed with code %d", resp.StatusCode)
	} else {
		return nil
	}
}

func (model *AuthController) Auth(username string, password string) (*authorization.AuthResponse, error) {
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
	data := &authorization.AuthResponse{}
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
