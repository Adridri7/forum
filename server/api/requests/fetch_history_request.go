package requests

import (
	"encoding/json"
	"forum/server"
	request "forum/server/admin_requests"
	"net/http"
)

func HistoryRequestHandler(w http.ResponseWriter, r *http.Request) {
	var params map[string]interface{}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user_uuid, ok := params["user_uuid"].(string)
	if user_uuid == "" || !ok {
		http.Error(w, "missing argument in request payload", http.StatusBadRequest)
		return
	}

	requestData, err := request.HistoryRequest(server.Db, user_uuid)
	if err != nil {
		http.Error(w, "Failed to fetch request for user", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(requestData); err != nil {
		http.Error(w, "Failed to encode request", http.StatusInternalServerError)
		return
	}
}
