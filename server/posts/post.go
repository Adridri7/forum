package posts

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/server"
	posts "forum/server/utils"
	"net/http"
	"strings"
	"time"
)

type Post struct {
	Post_uuid      string    `json:"post_uuid"`
	User_uuid      string    `json:"user_uuid"`
	Username       string    `json:"username"`
	ProfilePicture string    `json:"profile_picture"`
	Content        string    `json:"content"`
	Category       []string  `json:"categories"`
	Likes          int       `json:"likes"`
	Dislikes       int       `json:"dislikes"`
	Created_at     time.Time `json:"created_at"`
}

func CreatePost(db *sql.DB, r *http.Request, params map[string]interface{}) (*Post, error) {

	post_UUID, _ := posts.GenerateUUID()

	// extraire le uuid du cookie
	user_UUID, err := posts.GetUserFromCookie(r)

	if err != nil {
		return nil, fmt.Errorf("{erreur lors de la génération du uuid: %v}", err)
	}

	//user_UUID := params["user_uuid"].(string)
	content, contentOK := params["content"].(string)

	if !contentOK {
		return nil, errors.New("informations manquantes")
	}

	// Extraire les hashtags du content
	categories := posts.ExtractHashtags(content)

	createPostQuery := `INSERT INTO posts (post_uuid, content, user_uuid, categories, likes, dislikes, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	creationDate := time.Now()

	_, err = server.RunQuery(createPostQuery, post_UUID, content, user_UUID, strings.Join(categories, ","), 0, 0, creationDate)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création du post: %v", err)
	}

	newPost := &Post{
		Post_uuid:  post_UUID,
		User_uuid:  user_UUID,
		Content:    content,
		Category:   categories,
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
			Category:       strings.Split(row["categories"].(string), ","),
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

	rows, err := server.RunQuery(fetchAllPostsQuery)
	if err != nil {
		return nil, fmt.Errorf("database query failed: %v", err)
	}

	var posts []Post
	for _, row := range rows {
		post := Post{
			Post_uuid:      row["post_uuid"].(string),
			User_uuid:      row["user_uuid"].(string),
			Username:       row["username"].(string),
			ProfilePicture: row["profile_picture"].(string),
			Content:        row["content"].(string),
			Category:       strings.Split(row["categories"].(string), ","),
			Likes:          int(row["likes"].(int64)),
			Dislikes:       int(row["dislikes"].(int64)),
			Created_at:     row["created_at"].(time.Time),
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func DeletePost(db *sql.DB, params map[string]interface{}) error {
	post_UUID, post_UUIDOK := params["post_uuid"].(string)

	if !post_UUIDOK {
		return errors.New("informations manquantes")
	}

	deletePostQuery := `DELETE FROM posts WHERE post_uuid = ?`
	_, err := server.RunQuery(deletePostQuery, post_UUID)
	if err != nil {
		return fmt.Errorf("erreur lors de la suppression du post: %v", err)
	}

	return nil
}

func FetchAllCategories(db *sql.DB) ([]string, error) {
	fetchAllCategoriesQuery := `
        SELECT DISTINCT categories
        FROM posts
        WHERE categories IS NOT NULL AND categories <> ''`

	rows, err := server.RunQuery(fetchAllCategoriesQuery)
	if err != nil {
		return nil, fmt.Errorf("database query failed: %v", err)
	}

	var categories []string
	for _, row := range rows {
		// On récupère la colonne 'categories' et on la divise par les virgules
		if catStr, ok := row["categories"].(string); ok {
			for _, category := range strings.Split(catStr, ",") {
				// Éliminer les espaces superflus et vérifier si la catégorie n'est pas déjà présente
				trimmedCategory := strings.TrimSpace(category)
				if trimmedCategory != "" {
					categories = append(categories, "#"+trimmedCategory)
				}
			}
		}
	}

	// Utiliser un map pour éliminer les doublons
	categoryMap := make(map[string]struct{})
	for _, category := range categories {
		categoryMap[category] = struct{}{}
	}

	// Convertir le map en slice
	uniqueCategories := make([]string, 0, len(categoryMap))
	for category := range categoryMap {
		uniqueCategories = append(uniqueCategories, category)
	}

	return uniqueCategories, nil
}

// Mini algo pour trouver les tendances

/*----------------------------------------------------------------------------------------------------

Explication de la Fonction
Requête SQL :

La requête utilise une sous-requête pour séparer les catégories (hashtags) en lignes distinctes.
SUBSTRING_INDEX et numbers génèrent les catégories individuelles à partir d'une chaîne qui contient plusieurs catégories séparées par des virgules.
Le nombre d'occurrences de chaque catégorie est ensuite compté et celles ayant plus de 10 occurrences sont sélectionnées.
Traitement des Résultats :

La fonction retourne un map[string]int où la clé est le nom de la catégorie (hashtag) et la valeur est le nombre de posts qui utilisent cette catégorie.

----------------------------------------------------------------------------------------------------*/

func FetchCategoryRanking(db *sql.DB) (map[string]int, error) {
	// Requête pour compter le nombre d'occurrences de chaque catégorie
	fetchCategoryRankingQuery := `
        SELECT category, COUNT(*) AS count
        FROM (
            SELECT TRIM(SUBSTRING_INDEX(SUBSTRING_INDEX(categories, ',', numbers.n), ',', -1)) AS category
            FROM posts
            INNER JOIN (
                SELECT 1 AS n UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL 
                SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9 UNION ALL SELECT 10
            ) AS numbers ON CHAR_LENGTH(categories) - CHAR_LENGTH(REPLACE(categories, ',', '')) >= numbers.n - 1
        ) AS subquery
        WHERE category <> ''
        GROUP BY category
        HAVING count > 10
        ORDER BY count DESC`

	// Exécute la requête
	rows, err := server.RunQuery(fetchCategoryRankingQuery)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération du classement des catégories: %v", err)
	}

	// Préparez le mappage des catégories et leur compte
	categoryRanking := make(map[string]int)
	for _, row := range rows {
		if category, ok := row["category"].(string); ok {
			if count, ok := row["count"].(int64); ok {
				categoryRanking[category] = int(count)
			}
		}
	}

	return categoryRanking, nil
}

func FetchPostsByCategory(db *sql.DB, category string) ([]Post, error) {
	fetchPostsByCategoryQuery := `
		SELECT p.*, u.username, u.profile_picture 
		FROM posts p
		JOIN users u ON p.user_uuid = u.user_uuid
		WHERE p.categories LIKE ?
		ORDER BY p.created_at DESC`

	category = strings.TrimPrefix(category, "#")
	param := "%" + category + "%"

	// Log pour voir la requête et les paramètres
	fmt.Println("Executing Query:", fetchPostsByCategoryQuery)
	fmt.Println("With parameter:", param)

	rows, err := server.RunQuery(fetchPostsByCategoryQuery, param)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des posts par catégorie: %v", err)
	}

	var posts []Post
	for _, row := range rows {
		post := Post{
			Post_uuid:      row["post_uuid"].(string),
			User_uuid:      row["user_uuid"].(string),
			Username:       row["username"].(string),
			ProfilePicture: row["profile_picture"].(string),
			Content:        row["content"].(string),
			Category:       strings.Split(row["categories"].(string), ","),
			Likes:          int(row["likes"].(int64)),
			Dislikes:       int(row["dislikes"].(int64)),
			Created_at:     row["created_at"].(time.Time),
		}
		posts = append(posts, post)
	}

	return posts, nil
}
