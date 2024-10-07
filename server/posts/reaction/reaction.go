package reaction

import (
	"database/sql"
	"fmt"
	"forum/server"
)

func HandleLikeDislike(db *sql.DB, postID string, userID string, action string) error {
	// Vérifier si l'utilisateur a déjà liké ou disliké ce post
	checkQuery := `SELECT action FROM post_reactions WHERE post_uuid = ? AND user_uuid = ?`
	rows, err := server.RunQuery(checkQuery, postID, userID)
	if err != nil {
		fmt.Println("Erreur lors de la verif :", err)
		return fmt.Errorf("erreur lors de la vérification de la réaction: %v", err)
	}

	var updateQuery string
	var queryParams []interface{}

	if len(rows) == 0 {
		// L'utilisateur n'a pas encore réagi, insérer une nouvelle réaction
		updateQuery = `INSERT INTO post_reactions (post_uuid, user_uuid, action) VALUES (?, ?, ?)`
		queryParams = []interface{}{postID, userID, action}
	} else if rows[0]["action"] != action {
		// L'utilisateur change son like
		updateQuery = `UPDATE post_reactions SET action = ? WHERE post_uuid = ? AND user_uuid = ?`
		queryParams = []interface{}{action, postID, userID}
	} else {
		// L'utilisateur annule son like
		updateQuery = `DELETE FROM post_reactions WHERE post_uuid = ? AND user_uuid = ?`
		queryParams = []interface{}{postID, userID}
	}

	_, err = server.RunQuery(updateQuery, queryParams...)
	if err != nil {
		fmt.Println("erreur dans la query :", err)
		return fmt.Errorf("erreur lors de la mise à jour de la réaction: %v", err)
	}

	// Mettre à jour le compte des likes/dislikes dans la table posts
	updateCountQuery := `
        UPDATE posts 
        SET likes = (SELECT COUNT(*) FROM post_reactions WHERE post_uuid = ? AND action = 'like'),
            dislikes = (SELECT COUNT(*) FROM post_reactions WHERE post_uuid = ? AND action = 'dislike')
        WHERE post_uuid = ?`
	_, err = server.RunQuery(updateCountQuery, postID, postID, postID)
	if err != nil {
		fmt.Println("Erreur lors de la maj des react :", err)
		return fmt.Errorf("erreur lors de la mise à jour du compte des réactions: %v", err)
	}

	return nil
}
