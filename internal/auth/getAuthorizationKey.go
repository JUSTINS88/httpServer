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

func GetAPIKey(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")
	if authorization == "" {
		return "", errors.New("authorization header not found")
	}

	const prefix = "ApiKey "
	if !strings.HasPrefix(authorization, prefix) {
		return "", errors.New("invalid authorization scheme")
	}

	apikey := strings.TrimSpace(strings.TrimPrefix(authorization, prefix))
	if apikey == "" {
		return "", errors.New("ApiKey token is empty")
	}

	return apikey, nil
}
