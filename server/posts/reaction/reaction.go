package reaction

import (
	"database/sql"
	"fmt"
	"forum/server"
	"forum/server/api/notifications"
)

func HandleLikeDislike(db *sql.DB, postID string, userID string, action string) error {
	change := false
	// Vérifier si l'utilisateur a déjà liké ou disliké ce post
	checkQuery := `SELECT action FROM post_reactions WHERE post_uuid = ? AND user_uuid = ?`
	rows, err := server.RunQuery(checkQuery, postID, userID)
	if err != nil {
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
		change = true
		updateQuery = `DELETE FROM post_reactions WHERE post_uuid = ? AND user_uuid = ?`
		queryParams = []interface{}{postID, userID}
	}

	_, err = server.RunQuery(updateQuery, queryParams...)
	if err != nil {
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
		return fmt.Errorf("erreur lors de la mise à jour du compte des réactions: %v", err)
	}

	if !change {
		userIDQuery := `
		SELECT user_uuid FROM posts WHERE post_uuid = ?;
		`
		row, err := server.RunQuery(userIDQuery, postID)
		if err != nil {
			return fmt.Errorf("erreur lors de la mise à jour du compte des réactions: %v", err)
		}

		if len(row) == 0 {
			return fmt.Errorf("aucun résultat trouvé pour l'UUID de l'utilisateur")
		}

		user_UUID := ""

		if userIDValue, ok := row[0]["user_uuid"].(string); ok {
			user_UUID = userIDValue
		} else {
			return fmt.Errorf("erreur de conversion: user_uuid n'est pas une chaîne")
		}
		usernameQuery := `
		SELECT username FROM users
		WHERE user_uuid = ?`

		rows, err = server.RunQuery(usernameQuery, userID)
		if err != nil {
			return fmt.Errorf("database query failed: %v", err)
		}

		var username string

		for _, row := range rows {
			username = row["username"].(string)
		}

		if err = notifications.InsertNotification(db, user_UUID, action, "post", postID, username); err != nil {
			return fmt.Errorf("erreur lors de la création d'une notification: %v", err)
		}
	}

	return nil
}
