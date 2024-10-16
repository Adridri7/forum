package authentification

import (
	"errors"
	"net/http"
)

func GetUserFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return "", errors.New("user not logged in")
	}

	sessionData, exists := Sessions[cookie.Value]
	if !exists {
		return "invalid session", err
	}
	return sessionData.UUID, nil
}
