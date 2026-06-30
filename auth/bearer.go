package auth

import (
	"fmt"
	"net/http"
)

type bearer struct {
	token string
}

// NewBearer returns Credentials for the SDL and GraphQL APIs.
func NewBearer(token string) Credentials {
	return &bearer{token: token}
}

func (b *bearer) Apply(req *http.Request) error {
	if b.token == "" {
		return fmt.Errorf("auth: bearer token is empty")
	}
	req.Header.Set("Authorization", "Bearer "+b.token)
	return nil
}
