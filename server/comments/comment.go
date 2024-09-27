package comments

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/server"
	"time"
)

type Comment struct {
	Comment_id int       `json:"comment_id"`
	Post_uuid  string    `json:"post_uuid"`
	User_uuid  string    `json:"user_uuid"`
	Content    string    `json:"content"`
	Created_at time.Time `json:"created_at"`
}

// CreateComment crée un nouveau commentaire et l'insère dans la base de données
func CreateComment(db *sql.DB, params map[string]interface{}) (*Comment, error) {
	post_UUID, postOK := params["post_uuid"].(string)

	// !userOK
	user_UUID, _ := params["user_uuid"].(string)
	content, contentOK := params["content"].(string)

	//!userOK

	if !postOK || !contentOK {
		return nil, errors.New("informations manquantes pour le commentaire")
	}

	creationDate := time.Now()

	// Insertion du commentaire dans la table
	createCommentQuery := `INSERT INTO comments (post_uuid, user_uuid, content, created_at) VALUES (?, ?, ?, ?)`
	_, err := server.RunQuery(createCommentQuery, post_UUID, user_UUID, content, creationDate, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création du commentaire: %v", err)
	}

	// Récupération de l'ID généré automatiquement
	var lastID int

	err = db.QueryRow("SELECT last_insert_rowid()").Scan(&lastID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération de l'ID du commentaire: %v", err)
	}

	// Création de la structure Comment
	newComment := &Comment{
		Comment_id: lastID,
		Post_uuid:  post_UUID,
		User_uuid:  user_UUID,
		Content:    content,
		Created_at: creationDate,
	}

	return newComment, nil
}
