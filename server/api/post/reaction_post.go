package post

import (
	"encoding/json"
	"forum/server"
	"forum/server/posts/reaction"
	posts "forum/server/utils"

	"net/http"
)

type LikeDislikeRequest struct {
	PostID string `json:"postId"`
	Action string `json:"action"`
}

type LikeDislikeResponse struct {
	Likes    int `json:"likes"`
	Dislikes int `json:"dislikes"`
}

func HandleLikeDislikeAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var req LikeDislikeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Erreur de décodage JSON", http.StatusBadRequest)
		return
	}

	userUUID, err := posts.GetUserFromCookie(r)
	if err != nil {
		http.Error(w, "Utilisateur non authentifié", http.StatusUnauthorized)
		return
	}

	err = reaction.HandleLikeDislike(server.Db, req.PostID, userUUID, req.Action)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Récupérer les nouveaux compteurs
	getCountsQuery := `SELECT likes, dislikes FROM posts WHERE post_uuid = ?`
	rows, err := server.RunQuery(getCountsQuery, req.PostID)
	if err != nil || len(rows) == 0 {
		http.Error(w, "Erreur lors de la récupération des compteurs", http.StatusInternalServerError)
		return
	}

	response := LikeDislikeResponse{
		Likes:    int(rows[0]["likes"].(int64)),
		Dislikes: int(rows[0]["dislikes"].(int64)),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
