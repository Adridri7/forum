package comments

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/server"
	posts "forum/server/utils"
	"time"
)

type Comment struct {
	Comment_id     string    `json:"comment_id"`
	Post_uuid      string    `json:"post_uuid"`
	User_uuid      string    `json:"user_uuid"`
	Content        string    `json:"content"`
	Username       string    `json:"username"`
	ProfilePicture string    `json:"profile_picture"`
	Created_at     time.Time `json:"created_at"`
	Likes          int64     `json:"likes"`
	Dislikes       int64     `json:"dislikes"`
	Updated_at     time.Time `json:"update_at"`
}

type Post struct {
	Post_uuid      string    `json:"post_uuid"`
	User_uuid      string    `json:"user_uuid"`
	Username       string    `json:"username"`
	ProfilePicture string    `json:"profile_picture"`
	Content        string    `json:"content"`
	Category       []string  `json:"categories"`
	Likes          int64     `json:"likes"`
	Dislikes       int64     `json:"dislikes"`
	Created_at     time.Time `json:"created_at"`
	Post_image     string    `json:"post_image"`
	IsUpdated      bool      `json:"isUpdated"`
}

type Response struct {
	Comments []Comment `json:"comments"`
	Posts    []Post    `json:"posts"`
}

// CreateComment crée un nouveau commentaire et l'insère dans la base de données
func CreateComment(db *sql.DB, params map[string]interface{}) (*Comment, error) {

	comment_UUID, _ := posts.GenerateUUID()

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
	createCommentQuery := `INSERT INTO comments (comment_id, post_uuid, user_uuid, content, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err := server.RunQuery(createCommentQuery, comment_UUID, post_UUID, user_UUID, content, creationDate, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création du commentaire: %v", err)
	}

	// Création de la structure Comment
	newComment := &Comment{
		Comment_id: comment_UUID,
		Post_uuid:  post_UUID,
		User_uuid:  user_UUID,
		Content:    content,
		Created_at: creationDate,
	}

	return newComment, nil
}

// FetchAllComments récupère tous les commentaires de la base de données
func FetchAllComments(db *sql.DB) ([]Comment, error) {
	results, err := server.RunQuery("SELECT comment_id, post_uuid, user_uuid, content, created_at, updated_at FROM comments")
	if err != nil {
		return nil, err
	}

	var comments []Comment

	// ce qu'on veut renvoyer
	for _, row := range results {
		comment := Comment{}

		if commentUUID, ok := row["comment_id"].(string); ok {
			comment.Comment_id = commentUUID
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
		if username, ok := row["username"].(string); ok {
			comment.Username = username
		}
		// if profilePicture, ok := row["profile_picture"].(string); ok {
		// 	comment.ProfilePicture = profilePicture
		// }
		if updated_at, ok := row["updated_at"].(time.Time); ok {
			comment.Updated_at = updated_at
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

		fetchCommentquery = `
    SELECT c.comment_id, c.content, c.post_uuid, u.user_uuid, u.username, u.profile_picture, c.created_at, likes, dislikes, c.updated_at
    FROM comments AS c
    JOIN users AS u ON c.user_uuid = u.user_uuid
    WHERE c.post_uuid = ?`
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

	results, err := server.RunQuery(fetchCommentquery, param)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération du formulaire: %v", err)
	}

	var comments []Comment

	// Lecture des résultats à partir du slice de maps

	for _, row := range results {
		comment := Comment{}

		if commentUUID, ok := row["comment_id"].(string); ok {
			comment.Comment_id = commentUUID
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
		if username, ok := row["username"].(string); ok {
			comment.Username = username
		}
		if profilePicture, ok := row["profile_picture"].(string); ok {
			comment.ProfilePicture = profilePicture
		}

		if like, ok := row["likes"].(int64); ok {
			comment.Likes = like
		}

		if dislike, ok := row["dislikes"].(int64); ok {
			comment.Dislikes = dislike
		}

		if updated_at, ok := row["updated_at"].(time.Time); ok {
			comment.Updated_at = updated_at
		}
		comments = append(comments, comment)
	}

	return comments, nil

}

func DeleteComment(db *sql.DB, params map[string]interface{}) error {
	comment_ID, comment_IDOK := params["comment_id"].(string)

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

func FetchUserComments(db *sql.DB, user_uuid string) ([]Comment, error) {
	// Préparer la requête SQL pour récupérer les commentaires de l'utilisateur spécifié
	fetchUserCommentsQuery := `
        SELECT c.comment_id, c.content, c.post_uuid, u.user_uuid, u.username, u.profile_picture, c.created_at, c.likes, c.dislikes
        FROM comments c
        JOIN users u ON c.user_uuid = u.user_uuid
        WHERE c.user_uuid = ?  -- Filtrer par user_uuid
        ORDER BY c.created_at DESC`

	// Exécuter la requête
	rows, err := server.RunQuery(fetchUserCommentsQuery, user_uuid)
	if err != nil {
		return nil, fmt.Errorf("database query failed: %v", err)
	}

	var comments []Comment
	// Parcourir les résultats et remplir la structure Comment
	for _, row := range rows {
		comment := Comment{
			Comment_id:     row["comment_id"].(string),
			Post_uuid:      row["post_uuid"].(string),
			User_uuid:      row["user_uuid"].(string),
			Content:        row["content"].(string),
			Created_at:     row["created_at"].(time.Time),
			Username:       row["username"].(string),
			ProfilePicture: row["profile_picture"].(string),
			Likes:          (row["likes"].(int64)),
			Dislikes:       (row["dislikes"].(int64)),
		}
		if updatedAt, ok := row["updated_at"].(time.Time); ok {
			comment.Updated_at = updatedAt
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func FetchUserReactions(db *sql.DB, user_uuid string) (Response, error) {
	// Requête pour récupérer les likes sur les posts
	fetchUserLikePost := `
    SELECT p.post_uuid, p.content, u.user_uuid, u.username, u.profile_picture, p.created_at, p.likes, p.dislikes, p.post_image
    FROM post_reactions pr
    JOIN posts p ON pr.post_uuid = p.post_uuid
    JOIN users u ON p.user_uuid = u.user_uuid
    WHERE pr.user_uuid = ?  -- Filtrer par user_uuid de la réaction
    ORDER BY p.created_at DESC`

	// Requête pour récupérer les likes sur les commentaires
	fetchUserCommentReactionsQuery := `
    SELECT c.comment_id, c.content, c.post_uuid, u.user_uuid, u.username, u.profile_picture, c.created_at, c.likes, c.dislikes
    FROM comment_reactions cr
    JOIN comments c ON cr.comment_id = c.comment_id
    JOIN users u ON c.user_uuid = u.user_uuid
    WHERE cr.user_uuid = ?  -- Filtrer par user_uuid de la réaction
    AND cr.action = 'like'  -- Ne prendre en compte que les likes
    ORDER BY c.created_at DESC`

	// Exécuter les requêtes
	postRows, err := server.RunQuery(fetchUserLikePost, user_uuid)
	if err != nil {
		return Response{Comments: nil, Posts: nil}, fmt.Errorf("database query failed for posts: %v", err)
	}

	commentRows, err := server.RunQuery(fetchUserCommentReactionsQuery, user_uuid)
	if err != nil {
		return Response{Comments: nil, Posts: nil}, fmt.Errorf("database query failed for comments: %v", err)
	}

	var comments []Comment
	var posts []Post

	// Parcourir les résultats des posts
	for _, row := range postRows {
		post := Post{}
		// Extraire et vérifier chaque champ
		if postUUID, ok := row["post_uuid"].(string); ok {
			post.Post_uuid = postUUID
		}
		if userUUID, ok := row["user_uuid"].(string); ok {
			post.User_uuid = userUUID
		}
		if content, ok := row["content"].(string); ok {
			post.Content = content
		}
		if createdAt, ok := row["created_at"].(time.Time); ok {
			post.Created_at = createdAt
		}
		if username, ok := row["username"].(string); ok {
			post.Username = username
		}
		if profilePicture, ok := row["profile_picture"].(string); ok {
			post.ProfilePicture = profilePicture
		}
		if likes, ok := row["likes"].(int64); ok {
			post.Likes = likes
		}
		if dislikes, ok := row["dislikes"].(int64); ok {
			post.Dislikes = dislikes
		}

		if post_image, ok := row["post_image"].(string); ok {
			post.Post_image = post_image
		}

		posts = append(posts, post)
	}

	// Parcourir les résultats des commentaires
	for _, row := range commentRows {
		comment := Comment{}
		// Extraire et vérifier chaque champ
		if commentID, ok := row["comment_id"].(string); ok {
			comment.Comment_id = commentID
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
		if username, ok := row["username"].(string); ok {
			comment.Username = username
		}
		if profilePicture, ok := row["profile_picture"].(string); ok {
			comment.ProfilePicture = profilePicture
		}
		if likes, ok := row["likes"].(int64); ok {
			comment.Likes = likes
		}
		if dislikes, ok := row["dislikes"].(int64); ok {
			comment.Dislikes = dislikes
		}

		comments = append(comments, comment)
	}

	response := Response{
		Comments: comments,
		Posts:    posts,
	}

	return response, nil
}

func UpdateComment(db *sql.DB, params map[string]interface{}) error {
	commentID, ok1 := params["comment_id"].(string)
	updatedMessage, ok2 := params["content"].(string)

	if !ok1 || !ok2 {
		return fmt.Errorf("invalid parameters")
	}

	updatePostQuery := `
		UPDATE comments 
		SET content = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE comment_id = ?`

	_, err := server.RunQuery(updatePostQuery, updatedMessage, commentID)
	if err != nil {
		return fmt.Errorf("failed to update post: %v", err)
	}

	return nil
}
