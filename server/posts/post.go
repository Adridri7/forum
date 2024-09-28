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
		return nil, fmt.Errorf("erreur lors de la génération de uuid: %v", err)
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

	post_UUID, post_UUIDOK := params["post_uuid"].(string)
	user_UUID, user_UUIDOK := params["user_uuid"].(string)

	var fetchPostquery string

	// Param car on souhaite récuperer de la data depuis la db
	var param string

	// Préparer les requetes SQL et fetch si user_uuid ou post_uuid
	if post_UUIDOK {
		fetchPostquery = `SELECT * FROM posts WHERE post_uuid = ?`
		param = post_UUID
	} else if user_UUIDOK {
		fetchPostquery = `SELECT * FROM posts WHERE user_uuid = ?`
		param = user_UUID
	} else {
		return nil, errors.New("informations manquantes")
	}

	rows, err := server.RunQuery(fetchPostquery, param)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération du formulaire: %v", err)
	}

	var posts []Post

	// Lecture des résultats à partir du slice de maps
	for _, row := range rows {

		post := Post{
			Post_uuid:  row["post_uuid"].(string),
			User_uuid:  row["user_uuid"].(string),
			Content:    row["content"].(string),
			Category:   row["categories"].(string),
			Likes:      int(row["likes"].(int64)),     // Conversion de int64 vers int
			Dislikes:   int(row["dislikes"].(int64)),  // Conversion de int64 vers int
			Created_at: row["created_at"].(time.Time), // Conversion en time.Time
		}

		/*
			Si une colonne retournée par la base de données est NULL,
			alors l'accès direct avec row["clé"].(type) va provoquer un panic.
			Pour éviter cela, on s'assure que toutes les colonnes sont non-null
			(NOT NULL) ou on ajoute des vérifications supplémentaires, par exemple :
		*/

		if val, ok := row["post_uuid"].(string); ok {
			post.Post_uuid = val
		}
		posts = append(posts, post)
	}

	return posts, nil

}

// FetchAllPosts récupère tous les posts de la base de données
func FetchAllPosts(db *sql.DB) ([]Post, error) {

	results, err := server.RunQuery("SELECT post_uuid, user_uuid, content, categories, likes, dislikes, created_at FROM posts")
	if err != nil {
		return nil, err
	}

	// Convertir les résultats en slice de Post
	var posts []Post
	for _, row := range results {

		post := Post{
			Post_uuid:  row["post_uuid"].(string),
			User_uuid:  row["user_uuid"].(string),
			Content:    row["content"].(string),
			Category:   row["categories"].(string),
			Likes:      int(row["likes"].(int64)),
			Dislikes:   int(row["dislikes"].(int64)),
			Created_at: row["created_at"].(time.Time),
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
