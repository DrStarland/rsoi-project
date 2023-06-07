package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"tasks/pkg/myjson"
	"tasks/pkg/utils"
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
