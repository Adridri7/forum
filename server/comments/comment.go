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

// FetchAllComments récupère tous les commentaires de la base de données
func FetchAllComments(db *sql.DB) ([]Comment, error) {
	results, err := server.RunQuery("SELECT comment_id, post_uuid, user_uuid, content, created_at FROM comments")
	if err != nil {
		return nil, err
	}

	var comments []Comment
	for _, row := range results {
		comment := Comment{}

		if id, ok := row["comment_id"].(int64); ok {
			comment.Comment_id = int(id)
		}
		if postUUID, ok := row["post_uuid"].(string); ok {
			comment.Post_uuid = postUUID
		}
		if userUUID, ok := row["user_uuid"].(string); ok {
			comment.User_uuid = userUUID
		}
		if content, ok := row["content"].(string); ok {
			comment.Content = content
		}
		if createdAt, ok := row["created_at"].(time.Time); ok {
			comment.Created_at = createdAt
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

func FetchComment(db *sql.DB, params map[string]interface{}) ([]Comment, error) {

	comment_ID, comment_IDOK := params["comment_id"].(string)
	post_UUID, post_UUIDOK := params["post_uuid"].(string)
	user_UUID, user_UUIDOK := params["user_uuid"].(string)

	var fetchCommentquery string

	// Param car on souhaite récuperer de la data depuis la db
	var param string

	// Préparer les requetes SQL et fetch si user_uuid ou post_uuid
	if post_UUIDOK {
		fetchCommentquery = `SELECT * FROM comments WHERE post_uuid = ?`
		param = post_UUID
	} else if user_UUIDOK {
		fetchCommentquery = `SELECT * FROM comments WHERE user_uuid = ?`
		param = user_UUID
	} else if comment_IDOK {
		fetchCommentquery = `SELECT * FROM comments WHERE comment_id = ?`
		param = comment_ID
	} else {
		return nil, errors.New("informations manquantes")
	}

	rows, err := server.RunQuery(fetchCommentquery, param)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération du formulaire: %v", err)
	}

	var comments []Comment

	// Lecture des résultats à partir du slice de maps
	for _, row := range rows {

		comment := Comment{
			Comment_id: row["comment_id"].(int),
			Post_uuid:  row["post_uuid"].(string),
			User_uuid:  row["user_uuid"].(string),
			Content:    row["content"].(string),
			Created_at: row["created_at"].(time.Time),
		}

		/*
			Si une colonne retournée par la base de données est NULL,
			alors l'accès direct avec row["clé"].(type) va provoquer un panic.
			Pour éviter cela, on s'assure que toutes les colonnes sont non-null
			(NOT NULL) ou on ajoute des vérifications supplémentaires, par exemple :
		*/

		if val, ok := row["comment_id"].(int); ok {
			comment.Comment_id = val
		}
		comments = append(comments, comment)
	}

	return comments, nil

}

func DeleteComment(db *sql.DB, params map[string]interface{}) error {
	comment_ID, comment_IDOK := params["comment_id"].(int)

	if !comment_IDOK {
		return errors.New("informations manquantes")
	}

	// Exécution de la requête de suppression
	deleteCommentQuery := `DELETE FROM comments WHERE comment_id = ?`
	_, err := server.RunQuery(deleteCommentQuery, comment_ID)
	if err != nil {
		return fmt.Errorf("erreur lors de la suppression du commentaire: %v", err)
	}

	return nil
}
