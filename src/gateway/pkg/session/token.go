package session

// func RetrieveToken(w http.ResponseWriter, r *http.Request) *Token {
// 	reqToken := r.Header.Get("Authorization")
// 	if len(reqToken) == 0 {
// 		responses.TokenIsMissing(w)
// 		return nil
// 	}
// 	splitToken := strings.Split(reqToken, "Bearer ")
// 	tokenStr := splitToken[1]
// 	jwks := newJWKs(utils.Config.RawJWKS)
// 	tk := &Token{}

// 	token, err := jwt.ParseWithClaims(tokenStr, tk, jwks.Keyfunc)
// 	if err != nil || !token.Valid {
// 		responses.JwtAccessDenied(w)
// 		return nil
// 	}
// 	if time.Now().Unix()-tk.ExpiresAt > 0 {
// 		responses.TokenExpired(w)
// 		return nil
// 	}

// 	return tk
// }
