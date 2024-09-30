package post

import (
	"encoding/json"
	"forum/server"
	"forum/server/posts"
	"net/http"
)

func FetchPostHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	queryParams := r.URL.Query()
	userUUID := queryParams.Get("user_uuid")
	postUUID := queryParams.Get("post_uuid")

	params := map[string]interface{}{}

	if userUUID != "" {
		params["user_uuid"] = userUUID
	} else if postUUID != "" {
		params["post_uuid"] = postUUID
	}

	postData, err := posts.FetchPost(server.Db, params)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content", "application/json")
	json.NewEncoder(w).Encode(postData)
}

// Pour l'instant on test sans user_uuid
