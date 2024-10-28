package requests

import (
	"encoding/json"
	"forum/server"
	request "forum/server/admin_requests"
	authentification "forum/server/api/login"
	"net/http"
)

func FetchAdminRequestHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	role, err := authentification.GetUserInfoFromCookie(r)
	if err != nil {
		http.Error(w, "Error getting role from cookie", http.StatusInternalServerError)
	}

	if role != "admin" {
		http.Error(w, "Not allowed to get request", http.StatusUnauthorized)
		return
	}

	postData, err := request.FetchAdminRequest(server.Db)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content", "application/json")
	json.NewEncoder(w).Encode(postData)
}
