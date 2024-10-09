package comments

import (
	"encoding/json"
	"forum/server"
	"forum/server/posts/reaction"
	posts "forum/server/utils"

	"net/http"
)

type LikeDislikeRequest struct {
	CommentID string `json:"commentId"`
	Action    string `json:"action"`
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

func HandleLikeDislikeCommentAPI(w http.ResponseWriter, r *http.Request) {
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

	err = reaction.HandleLikeDislikeComment(server.Db, req.CommentID, userUUID, req.Action)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Récupérer les nouveaux compteurs
	getCountsQuery := `SELECT likes, dislikes FROM comments WHERE comment_id = ?`
	rows, err := server.RunQuery(getCountsQuery, req.CommentID)
	if err != nil || len(rows) == 0 {
		http.Error(w, "Erreur lors de la récupération des compteurs", http.StatusInternalServerError)
		return
	}

	// Vérifie si l'utilisateur a aimé ou non
	userReactionQuery := `SELECT 
		(SELECT COUNT(*) FROM comment_reactions WHERE comment_id = ? AND user_uuid = ? AND action = 'like') AS hasLiked,
		(SELECT COUNT(*) FROM comment_reactions WHERE comment_id = ? AND user_uuid = ? AND action = 'dislike') AS hasDisliked`

	var hasLiked, hasDisliked bool
	err = server.Db.QueryRow(userReactionQuery, req.CommentID, userUUID, req.CommentID, userUUID).Scan(&hasLiked, &hasDisliked)
	if err != nil {
		http.Error(w, "Erreur lors de la vérification de la réaction de l'utilisateur", http.StatusInternalServerError)
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
