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
	Post_uuid  string    `json:"post_uuid"`
	User_uuid  string    `json:"user_uuid"`
	Content    string    `json:"content"`
	Category   string    `json:"categories"`
	Likes      int       `json:"likes"`
	Dislikes   int       `json:"dislikes"`
	Created_at time.Time `json:"created_at"`
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
