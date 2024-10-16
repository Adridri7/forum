package user

import (
	"net/http"
	"time"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	}

	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
}
