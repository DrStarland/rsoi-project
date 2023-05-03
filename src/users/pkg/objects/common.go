package objects

import (
	_ "encoding/json"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
	IdToken     string `json:"id_token"`
}

type UserCreateRequest struct {
	Profile struct {
		Firstname   string `json:"firstName,omitempty"`
		Lastname    string `json:"lastName,omitempty"`
		Email       string `json:"email,omitempty"`
		Login       string `json:"login,omitempty"`
		Mobilephone string `json:"mobilePhone,omitempty"`
		UserType    string `json:"userType,omitempty"`
	} `json:"profile,omitempty"`
	Credentials struct {
		Password struct {
			Value string `json:"value"`
		} `json:"password"`
	} `json:"credentials"`
	GroupIds []string `json:"groupIds,omitempty"`
}
