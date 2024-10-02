package posts

import (
	"errors"
	"net/http"
	"strings"
)

func GetUserFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("UserLogged")
	if err != nil {
		return "", errors.New("user not logged in")
	}

	// Décoder la valeur du cookie
	// Supposons que la valeur du cookie soit au format "uuid:username:email"
	parts := strings.Split(cookie.Value, "|")
	if len(parts) < 3 {
		return "", errors.New("invalid cookie format")
	}

	return parts[0], nil // Retourne le user_uuid
}
