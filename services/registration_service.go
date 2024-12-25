package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xanagit/kotoquiz-api/config"
	"io"
	"net/http"
)

type RegistrationService interface {
	RegisterUser(username, email, password string) error
}

type RegistrationServiceImpl struct {
	KeycloakConfig *config.KeycloakConfig
}

// Make sure that RegistrationServiceImpl implements RegistrationService
var _ RegistrationService = (*RegistrationServiceImpl)(nil)

type KeycloakToken struct {
	AccessToken string `json:"access_token"`
}

type KeycloakUser struct {
	Username    string               `json:"username"`
	Email       string               `json:"email"`
	Enabled     bool                 `json:"enabled"`
	Credentials []KeycloakCredential `json:"credentials"`
}

type KeycloakCredential struct {
	Type      string `json:"type"`
	Value     string `json:"value"`
	Temporary bool   `json:"temporary"`
}

func (s *RegistrationServiceImpl) getAdminToken() (string, error) {
	url := fmt.Sprintf("%s/realms/master/protocol/openid-connect/token", s.KeycloakConfig.BaseUrl)

	data := fmt.Sprintf("client_id=%s&grant_type=password&username=%s&password=%s",
		s.KeycloakConfig.AdminCliClientId, s.KeycloakConfig.User, s.KeycloakConfig.Password)

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	if err != nil {
		return "", fmt.Errorf("error creating token request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error requesting token: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("error closing response body: ", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error getting token, status: %d", resp.StatusCode)
	}

	var token KeycloakToken
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", fmt.Errorf("error decoding token response: %v", err)
	}

	return token.AccessToken, nil
}

func (s *RegistrationServiceImpl) RegisterUser(username, email, password string) error {
	token, err := s.getAdminToken()
	if err != nil {
		return fmt.Errorf("failed to get admin token: %v", err)
	}

	user := KeycloakUser{
		Username: username,
		Email:    email,
		Enabled:  true,
		Credentials: []KeycloakCredential{
			{
				Type:      "password",
				Value:     password,
				Temporary: false,
			},
		},
	}

	userJSON, mErr := json.Marshal(user)
	if mErr != nil {
		return fmt.Errorf("error marshaling user data: %v", err)
	}

	url := fmt.Sprintf("%s/admin/realms/%s/users", s.KeycloakConfig.BaseUrl, s.KeycloakConfig.Realm)

	req, httpErr := http.NewRequest("POST", url, bytes.NewBuffer(userJSON))
	if httpErr != nil {
		return fmt.Errorf("error creating registration request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, reqErr := client.Do(req)
	if reqErr != nil {
		return fmt.Errorf("error making registration request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("error closing response body: ", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error registering user, status: %d", resp.StatusCode)
	}

	return nil
}
