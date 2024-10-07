package post

import (
	"encoding/json"
	"forum/server"
	"forum/server/posts"
	"net/http"
)

func FetchPostsMostLikedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer tous les posts
	postMostLiked, err := posts.FetchPostMostLiked(server.Db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(postMostLiked)
}
