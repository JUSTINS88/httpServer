package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")
	if authorization == "" {
		return "", errors.New("authorization header not found")
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authorization, prefix) {
		return "", errors.New("invalid authorization scheme")
	}

	token := strings.TrimSpace(strings.TrimPrefix(authorization, prefix))
	if token == "" {
		return "", errors.New("bearer token is empty")
	}

	return token, nil
}
