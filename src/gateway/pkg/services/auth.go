package services

import (
	"encoding/json"
	"fmt"
	"gateway/pkg/myjson"
	"gateway/pkg/utils"
	"net/http"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type AuthController struct {
	Client *http.Client
	Logger *zap.SugaredLogger
}

func NewAuthController(client *http.Client, logger *zap.SugaredLogger) *AuthController {
	return &AuthController{Client: client, Logger: logger}
}

type Token struct {
	jwt.StandardClaims
	Role string `json:"role,omitempty"`
	UID  string
}

func newJWKs(rawJWKS string, logger *zap.SugaredLogger) *keyfunc.JWKS {
	// Get the JWKS as JSON.
	/// logger.Infoln("GOvno: ", rawJWKS)
	jwksJSON := json.RawMessage(rawJWKS)
	// jwksJSON, _ := myjson.To(rawJWKS)

	// logger.Infoln("GOvno 2: ", jwksJSON)

	jwks, err := keyfunc.NewJSON(jwksJSON)

	if err != nil {
		panic(errors.Wrap(err, "inside jwk func generator"))
	}
	return jwks
}

func RetrieveToken(w http.ResponseWriter, r *http.Request, logger *zap.SugaredLogger) (*Token, error) {
	reqToken := r.Header.Get("Authorization")
	if len(reqToken) == 0 {
		myjson.JSONError(w, http.StatusUnauthorized, "Missing auth token")
		return nil, fmt.Errorf("TokenIsMissed")
	}
	// logger.Infoln("KAL1: ", reqToken)

	_tokenStr := strings.Split(reqToken, "Bearer ")
	var tokenStr string
	if len(_tokenStr) == 2 {
		tokenStr = _tokenStr[1]
	} else {
		tokenStr = _tokenStr[0]
	}

	//  logger.Infoln("KAL2: ", tokenStr)

	jwks := newJWKs(utils.Config.RawJWKS, logger)
	tk := &Token{}

	//  logger.Infoln("KAL4: ", jwks)
	token, err := jwt.ParseWithClaims(tokenStr, tk, jwks.Keyfunc)
	claims := jwt.MapClaims{}
	_, _ = jwt.ParseWithClaims(tokenStr, claims, jwks.Keyfunc)
	// ... error handling
	// do something with decoded claims
	tk.UID = claims["uid"].(string)

	//  logger.Infoln("KAL5: ", token)
	// token.Valid = true
	if err != nil {
		myjson.JSONError(w, http.StatusUnauthorized, errors.Wrap(err, "JwtAccessDenied").Error())
	}

	if !token.Valid {
		return nil, errors.Wrap(err, "JwtAccessDenied")
	}

	// проверка времени существования токена
	if time.Now().Unix()-tk.ExpiresAt > 0 {
		myjson.JSONError(w, http.StatusUnauthorized, "jwt-token expired")
		return nil, errors.New("token expired")
	}
	tk.UID = claims["uid"].(string)
	// logger.Infoln("KAL3: ", token)
	return tk, nil
}

// func (model *AuthController) Create(request *authorization.UserCreateRequest) error {
// 	request.GroupIds = []string{utils.Config.Okta.ClientGroup}
// 	request.Profile.UserType = utils.User.String()

// 	req_body, err := json.Marshal(request)
// 	if err != nil {
// 		return err
// 	}
// 	req, _ := http.NewRequest(
// 		"POST",
// 		fmt.Sprintf("%s/api/v1/gateway/?activate=true", utils.Config.Okta.Endpoint),
// 		bytes.NewBuffer(req_body),
// 	)
// 	req.Header.Add("Content-Type", "application/json")
// 	req.Header.Add("Authorization", fmt.Sprintf("SSWS %s", utils.Config.Okta.OktetoToken))

// 	resp, err := model.client.Do(req)
// 	if err != nil {
// 		return err
// 	} else if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("auth failed with code %d", resp.StatusCode)
// 	} else {
// 		return nil
// 	}
// }

// func (model *AuthController) Auth(username string, password string) (*authorization.AuthResponse, error) {
// 	authRequest := url.Values{}
// 	authRequest.Set("scope", "openid")
// 	authRequest.Set("grant_type", "password")
// 	authRequest.Set("username", username)
// 	authRequest.Set("password", password)
// 	authRequest.Set("client_id", utils.Config.Okta.ClientId)
// 	authRequest.Set("client_secret", utils.Config.Okta.ClientSecret)
// 	encodedData := authRequest.Encode()

// 	req, _ := http.NewRequest(
// 		"POST",
// 		fmt.Sprintf("%s/oauth2/default/v1/token", utils.Config.Okta.Endpoint),
// 		strings.NewReader(encodedData),
// 	)
// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	req.Header.Add("Content-Length", strconv.Itoa(len(authRequest.Encode())))

// 	resp, err := model.Client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
// 	data := &authorization.AuthResponse{}
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("auth failed, code: %d", resp.StatusCode)
// 	}

// 	err = json.Unmarshal(body, data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return data, nil
// }

// import (
// 	"encoding/json"
// 	"fmt"
// 	"gateway/pkg/controllers/responses"
// 	"gateway/pkg/objects"
// 	"gateway/pkg/utils"
// 	"io"
// 	"strings"

// 	"net/http"
// 	"time"

// 	"github.com/MicahParks/keyfunc"
// 	"github.com/golang-jwt/jwt/v4"
// 	"github.com/gorilla/mux"
// )

// type Token struct {
// 	jwt.StandardClaims
// }

// func newJWKs(rawJWKS string) *keyfunc.JWKS {
// 	jwksJSON := json.RawMessage(rawJWKS)
// 	jwks, err := keyfunc.NewJSON(jwksJSON)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return jwks
// }

// var JwtAuthentication = func(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if token := RetrieveToken(w, r); token != nil {
// 			r.Header.Set("X-User-Name", token.Subject)
// 			next.ServeHTTP(w, r)
// 		}
// 	})
// }

// type AuthCtrl struct {
// 	client *http.Client
// }
