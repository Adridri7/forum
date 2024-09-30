package comments

import (
	"encoding/json"
	"forum/server"
	"forum/server/comments"
	"net/http"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var params map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	newComment, err := comments.CreateComment(server.Db, params)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content", "application/json")
	json.NewEncoder(w).Encode(newComment)
}
