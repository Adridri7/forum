package posts

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/server"
	posts "forum/server/utils"
	"time"
)

type Post struct {
	Post_uuid      string    `json:"post_uuid"`
	User_uuid      string    `json:"user_uuid"`
	Username       string    `json:"username"`
	ProfilePicture string    `json:"profile_picture"`
	Content        string    `json:"content"`
	Category       string    `json:"categories"`
	Likes          int       `json:"likes"`
	Dislikes       int       `json:"dislikes"`
	Created_at     time.Time `json:"created_at"`
}

// CreatePost remplit le rôle de constructeur en initialisant et
// en retournant un pointeur vers une structure Post pas besoin de
// constructor en go.

func CreatePost(db *sql.DB, params map[string]interface{}) (*Post, error) {

	post_UUID, err := posts.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la génération du uuid: %v", err)
	}

	user_UUID, user_UUIDOK := params["user_uuid"].(string)
	content, contentOK := params["content"].(string)
	category, categoryOK := params["categories"].(string)

	if !user_UUIDOK || !contentOK || !categoryOK {
		return nil, errors.New("informations manquantes")
	}

	createPostQuery := `INSERT INTO posts (post_uuid, user_uuid, content, categories, likes, dislikes, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	creationDate := time.Now()

	_, err = server.RunQuery(createPostQuery, post_UUID, user_UUID, content, category, 0, 0, creationDate)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création du post: %v", err)
	}

	newPost := &Post{
		Post_uuid:  post_UUID,
		User_uuid:  user_UUID,
		Content:    content,
		Category:   category,
		Likes:      0,
		Dislikes:   0,
		Created_at: creationDate,
	}

	return newPost, nil
}

func FetchPost(db *sql.DB, params map[string]interface{}) ([]Post, error) {
	var fetchPostquery string
	var param string

	if post_UUID, ok := params["post_uuid"].(string); ok {
		fetchPostquery = `
			SELECT p.*, u.username, u.profile_picture 
			FROM posts p
			JOIN users u ON p.user_uuid = u.user_uuid
			WHERE p.post_uuid = ?`
		param = post_UUID
	} else if user_UUID, ok := params["user_uuid"].(string); ok {
		fetchPostquery = `
			SELECT p.*, u.username, u.profile_picture 
			FROM posts p
			JOIN users u ON p.user_uuid = u.user_uuid
			WHERE p.user_uuid = ?`
		param = user_UUID
	} else {
		return nil, errors.New("informations manquantes")
	}

	rows, err := server.RunQuery(fetchPostquery, param)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération du formulaire: %v", err)
	}

	var posts []Post

	for _, row := range rows {
		post := Post{
			Post_uuid:      row["post_uuid"].(string),
			User_uuid:      row["user_uuid"].(string),
			Username:       row["username"].(string),
			ProfilePicture: row["profile_picture"].(string),
			Content:        row["content"].(string),
			Category:       row["categories"].(string),
			Likes:          int(row["likes"].(int64)),
			Dislikes:       int(row["dislikes"].(int64)),
			Created_at:     row["created_at"].(time.Time),
		}

		if val, ok := row["post_uuid"].(string); ok {
			post.Post_uuid = val
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func FetchAllPosts(db *sql.DB) ([]Post, error) {
	fetchAllPostsQuery := `
		SELECT p.*, u.username, u.profile_picture 
		FROM posts p
		JOIN users u ON p.user_uuid = u.user_uuid
		ORDER BY p.created_at DESC`

	results, err := server.RunQuery(fetchAllPostsQuery)
	if err != nil {
		return nil, err
	}

	var posts []Post
	for _, row := range results {
		post := Post{}

		// Utiliser des assertions de type avec vérification de valeur nulle
		if v, ok := row["post_uuid"]; ok && v != nil {
			post.Post_uuid = v.(string)
		}
		if v, ok := row["user_uuid"]; ok && v != nil {
			post.User_uuid = v.(string)
		}
		if v, ok := row["username"]; ok && v != nil {
			post.Username = v.(string)
		}
		if v, ok := row["profile_picture"]; ok && v != nil {
			post.ProfilePicture = v.(string)
		}
		if v, ok := row["content"]; ok && v != nil {
			post.Content = v.(string)
		}
		if v, ok := row["categories"]; ok && v != nil {
			post.Category = v.(string)
		}
		if v, ok := row["likes"]; ok && v != nil {
			post.Likes = int(v.(int64))
		}
		if v, ok := row["dislikes"]; ok && v != nil {
			post.Dislikes = int(v.(int64))
		}
		if v, ok := row["created_at"]; ok && v != nil {
			post.Created_at = v.(time.Time)
		}

		posts = append(posts, post)
	}

	return posts, nil
}

/*--------Important /!\ Cela fonctionne et est la func finale quand user fait décommenter----------*/

/*
func DeletePost(db *sql.DB, params map[string]interface{}) (*Post, error) {
	post_UUID, post_UUIDOK := params["post_uuid"].(string)
	user_UUID, user_UUIDOK := params["user_uuid"].(string)

	if !post_UUIDOK || !user_UUIDOK {
		return nil, errors.New("informations manquantes")
	}

	// Vérification de l'existence du post et de l'auteur
	checkPostQuery := `SELECT user_uuid FROM posts WHERE post_uuid = ?`

	var postUserUUID string
	_, err := server.RunQuery(checkPostQuery, post_UUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("post non trouvé")
		}
		return nil, fmt.Errorf("erreur lors de la vérification du post: %v", err)
	}

	// Vérifier que l'utilisateur a le droit de supprimer le post
	if postUserUUID != user_UUID {
		return nil, errors.New("vous n'avez pas les droits pour supprimer ce post")
	}

	// Exécution de la requête de suppression
	deletePostQuery := `DELETE FROM posts WHERE post_uuid = ?`
	_, err = server.RunQuery(deletePostQuery, post_UUID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la suppression du post: %v", err)
	}

	return &Post{Post_uuid: post_UUID}, nil
}
-------------------------------------------------------------------------------------------------*/

func DeletePost(db *sql.DB, params map[string]interface{}) error {
	post_UUID, post_UUIDOK := params["post_uuid"].(string)

	if !post_UUIDOK {
		return errors.New("informations manquantes")
	}

	// Exécution de la requête de suppression
	deletePostQuery := `DELETE FROM posts WHERE post_uuid = ?`
	_, err := server.RunQuery(deletePostQuery, post_UUID)
	if err != nil {
		return fmt.Errorf("erreur lors de la suppression du post: %v", err)
	}

	return nil
}
