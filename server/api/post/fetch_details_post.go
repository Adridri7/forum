package post

import (
	"encoding/json"
	"forum/server"
	"forum/server/posts"
	"net/http"
)

func PostDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var params map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	postId, ok := params["post_uuid"].(string)
	if !ok || postId == "" {
		http.Error(w, "Missing or invalid postId", http.StatusBadRequest)
		return
	}

	postData, err := posts.FetchDetailsPost(server.Db, postId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(postData); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
