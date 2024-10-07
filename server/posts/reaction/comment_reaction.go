package reaction

import (
	"database/sql"
	"fmt"
	"forum/server"
)

func HandleLikeDislikeComment(db *sql.DB, CommentID string, userID string, action string) error {
	// Vérifier si l'utilisateur a déjà liké ou disliké ce post
	checkQuery := `SELECT action FROM comment_reactions WHERE post_uuid = ? AND user_uuid = ?`
	rows, err := server.RunQuery(checkQuery, CommentID, userID)
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification de la réaction: %v", err)
	}

	var updateQuery string
	var queryParams []interface{}

	if len(rows) == 0 {
		// L'utilisateur n'a pas encore réagi, insérer une nouvelle réaction
		updateQuery = `INSERT INTO comment_reactions (post_uuid, user_uuid, action) VALUES (?, ?, ?)`
		queryParams = []interface{}{CommentID, userID, action}
	} else if rows[0]["action"] != action {
		// L'utilisateur change son like
		updateQuery = `UPDATE comment_reactions SET action = ? WHERE post_uuid = ? AND user_uuid = ?`
		queryParams = []interface{}{action, CommentID, userID}
	} else {
		// L'utilisateur annule son like
		updateQuery = `DELETE FROM comment_reactions WHERE post_uuid = ? AND user_uuid = ?`
		queryParams = []interface{}{CommentID, userID}
	}

	_, err = server.RunQuery(updateQuery, queryParams...)
	if err != nil {
		return fmt.Errorf("erreur lors de la mise à jour de la réaction: %v", err)
	}

	// Met à jour le compte des likes/dislikes dans la table posts
	updateCountQuery := `
        UPDATE posts 
        SET likes = (SELECT COUNT(*) FROM comment_reactions WHERE post_uuid = ? AND action = 'like'),
            dislikes = (SELECT COUNT(*) FROM comment_reactions WHERE post_uuid = ? AND action = 'dislike')
        WHERE post_uuid = ?`
	_, err = server.RunQuery(updateCountQuery, CommentID, CommentID, CommentID)
	if err != nil {
		return fmt.Errorf("erreur lors de la mise à jour du compte des réactions: %v", err)
	}

	return nil
}
