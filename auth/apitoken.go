package auth

import (
	"fmt"
	"net/http"
)

type apiToken struct {
	token string
}

// NewApiToken returns Credentials for the REST MGMT API.
func NewApiToken(token string) Credentials {
	return &apiToken{token: token}
}

func (a *apiToken) Apply(req *http.Request) error {
	if a.token == "" {
		return fmt.Errorf("auth: API token is empty")
	}
	req.Header.Set("Authorization", "ApiToken "+a.token)
	return nil
}
