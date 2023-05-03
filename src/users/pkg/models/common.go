package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"users/pkg/objects"
	"users/pkg/utils"
)

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
	body, err := ioutil.ReadAll(resp.Body)
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
