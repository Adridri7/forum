package comments

import (
	"encoding/json"
	"forum/server"
	"forum/server/comments"
	"net/http"
)

func UpdateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var params map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Appeler la fonction pour mettre à jour le post
	if err := comments.UpdateComment(server.Db, params); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
