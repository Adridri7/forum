package post

import (
	"encoding/json"
	"fmt"
	"forum/server"
	authentification "forum/server/api/login"
	"forum/server/posts/reaction"

	"net/http"
)

type LikeDislikeRequest struct {
	PostID string `json:"postId"`
	Action string `json:"action"`
}

type LikeDislikeResponse struct {
	Likes        int          `json:"likes"`
	Dislikes     int          `json:"dislikes"`
	UserReaction UserReaction `json:"userReaction"`
}

type UserReaction struct {
	UserUUID    string `json:"user_uuid"`
	HasLiked    bool   `json:"hasLiked"`
	HasDisliked bool   `json:"hasDisliked"` // Si tu gères aussi les dislikes
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

	userUUID, err := authentification.GetUserFromCookie(r)
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

	// Vérifie si l'utilisateur a aimé ou non
	userReactionQuery := `SELECT 
		(SELECT COUNT(*) FROM post_reactions WHERE post_uuid = ? AND user_uuid = ? AND action = 'like') AS hasLiked,
		(SELECT COUNT(*) FROM post_reactions WHERE post_uuid = ? AND user_uuid = ? AND action = 'dislike') AS hasDisliked`

	var hasLiked, hasDisliked bool
	err = server.Db.QueryRow(userReactionQuery, req.PostID, userUUID, req.PostID, userUUID).Scan(&hasLiked, &hasDisliked)
	if err != nil {
		http.Error(w, "Erreur lors de la vérification de la réaction de l'utilisateur", http.StatusInternalServerError)
		fmt.Println("Erreur lors de la vérif de la réac : ", err)
		return
	}

	response := LikeDislikeResponse{
		Likes:    int(rows[0]["likes"].(int64)),
		Dislikes: int(rows[0]["dislikes"].(int64)),
		UserReaction: UserReaction{
			UserUUID:    userUUID,
			HasLiked:    hasLiked,
			HasDisliked: hasDisliked,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
