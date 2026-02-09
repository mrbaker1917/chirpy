package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", errors.New("no authorization header found")
	}

	apikey, found := strings.CutPrefix(auth, "ApiKey")
	if !found {
		return "", errors.New("no ApiKey prefix found")
	}
	apikey = strings.TrimSpace(apikey)

	return apikey, nil

}
