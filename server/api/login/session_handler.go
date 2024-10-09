package authentification

import (
	"encoding/json"
	"net/http"

	users "forum/server/api/user"
)

var Sessions = map[string]users.User{}

func GetSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Session not found"))
		cookie.MaxAge = -1
		return
	}

	sessionData, exists := Sessions[cookie.Value]
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid session"))
		cookie.MaxAge = -1
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sessionData)
}
