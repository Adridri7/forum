package comments

import (
	"encoding/json"
	"forum/server"
	"forum/server/comments"
	"net/http"
)

func CommentDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var params map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	commentId, ok := params["comment_id"].(string)
	if !ok || commentId == "" {
		http.Error(w, "Missing or invalid postId", http.StatusBadRequest)
		return
	}

	postData, err := comments.FetchDetailsComment(server.Db, commentId)
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
