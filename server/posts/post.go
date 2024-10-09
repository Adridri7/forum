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
	Post_image     string    `json:"post_image"`
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

	image, imageOK := params["post_image"].(string)

	if !imageOK {
		image = ""
	}

	if !contentOK {
		return nil, errors.New("informations manquantes")
	}

	// Extraire les hashtags du content
	categories := posts.ExtractHashtags(content)

	createPostQuery := `INSERT INTO posts (post_uuid, content, user_uuid, categories, likes, dislikes, created_at, post_image) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	creationDate := time.Now()

	_, err = server.RunQuery(createPostQuery, post_UUID, content, user_UUID, strings.Join(categories, ","), 0, 0, creationDate, image)
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
		Post_image: image,
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
		WHERE p.post_uuid = ?`
		//	WHERE p.user_uuid = ?`
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

		if _, ok := row["post_image"]; ok {
			post.Post_image = row["post_image"].(string)
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

		if data, ok := row["post_image"]; ok && data != nil {
			post.Post_image = row["post_image"].(string)
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

	/*---------------------------------------------------------------------------------------------------

			  Cette requête a pour objectif de compter et de classer les catégories d'un ensemble de posts.
			  Les catégories sont initialement stockées sous forme de chaînes de caractères, avec plusieurs
			  catégories séparées par des virgules dans chaque post (par exemple : "sports,news,technology").
			  La requête utilise une CTE récursive (Common Table Expression) pour décomposer ces chaînes de caractères
			  en catégories individuelles, puis compte leur fréquence d'apparition.

			  WITH RECURSIVE : Démarre une expression de table commune (CTE) récursive, appelée split.

			 - split(category, rest) : La CTE split est définie avec deux colonnes :
			 - category : Pour stocker chaque catégorie extraite individuellement.
			 - rest : Pour stocker la portion restante de la chaîne de catégories qui doit encore être traitée.

			 - SELECT '' : Cette première partie de la requête initialise la CTE avec une valeur vide pour la colonne category (car on n'a pas encore extrait de catégories).
			 - categories || ',' : Concatène la chaîne de catégories avec une virgule finale (categories + ','), pour s'assurer que chaque catégorie est séparée proprement par une virgule, même si la chaîne se termine par la dernière catégorie.
	         - FROM posts : Prend les données de la table posts. Cela génère une ligne pour chaque post avec la liste complète de catégories, prête à être traitée.

			 UNION ALL : Combine les résultats de la requête initiale avec les résultats générés par l'itération suivante de la CTE.

			 ---------------------------------------------------------------------------------------
			 /!\ UNION ALL est utilisé car il inclut tous les résultats, y compris les doublons. /!\
			 ---------------------------------------------------------------------------------------

			 SELECT : Cette partie de la requête est la portion récursive de la CTE, qui va extraire chaque catégorie une par une de la chaîne de caractères rest.

			 substr(rest, 0, instr(rest, ',')) :

			 instr(rest, ',') : Trouve la position de la première virgule dans la chaîne rest. Cette position correspond à la fin de la prochaine catégorie.
			 substr(rest, 0, instr(rest, ',')) : Utilise substr pour extraire la sous-chaîne (catégorie) de rest, qui commence à l'index 0 jusqu'à la position de la virgule (exclus).

			 Par exemple, si rest est "sports,news,technology,", alors substr extrait "sports" lors de la première itération.
			 substr(rest, instr(rest, ',') + 1) :

			Utilise substr pour extraire la portion restante de rest après la première virgule trouvée. Cela devient la nouvelle chaîne rest pour la prochaine itération.
			Par exemple, si rest est "sports,news,technology,", alors rest devient "news,technology," après la première itération.

			SELECT category, COUNT(*) AS count : Sélectionne chaque category unique et calcule combien de fois elle apparaît.
			FROM split : Utilise les résultats de la CTE split.
			WHERE category != '' : Exclut toutes les lignes où category est une chaîne vide (celles-ci proviennent de la première ligne initiale SELECT '').
			GROUP BY category : Regroupe les résultats par category pour obtenir le nombre total de chaque catégorie.
			HAVING count > 10 : Filtre pour inclure seulement les catégories qui apparaissent plus de 10 fois.
			ORDER BY count DESC : Trie les résultats par le nombre de catégories en ordre décroissant pour afficher d'abord les catégories les plus fréquentes.

	/*---------------------------------------------------------------------------------------------------*/

	fetchCategoryRankingQuery := `
        WITH RECURSIVE split(category, rest) AS (
            SELECT '', categories || ',' FROM posts
            UNION ALL
            SELECT
                substr(rest, 0, instr(rest, ',')),
                substr(rest, instr(rest, ',') + 1)
            FROM split
            WHERE rest != ''
        )
        SELECT category, COUNT(*) AS count
        FROM split
        WHERE category != ''
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

		if data, ok := row["post_image"]; ok && data != nil {
			post.Post_image = row["post_image"].(string)
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func FetchPostMostLiked(db *sql.DB) (map[string]int, error) {
	// Requête pour récupérer les posts les plus likés
	fetchPostMostLikedQuery := `
        SELECT post_uuid, SUM(likes) AS likes_count
        FROM posts
        GROUP BY post_uuid
        ORDER BY likes_count DESC;`

	// Exécute la requête en utilisant server.RunQuery
	rows, err := server.RunQuery(fetchPostMostLikedQuery)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des posts les plus likés: %v", err)
	}

	// Prépare le mappage des posts et leur nombre de likes
	postMostLikedRanking := make(map[string]int)
	for _, row := range rows {
		if postID, ok := row["post_uuid"].(string); ok {
			if count, ok := row["likes_count"].(int64); ok {
				postMostLikedRanking[postID] = int(count)
			}
		}
	}

	return postMostLikedRanking, nil
}

//  HAVING likes_count > 10

func FetchUserPosts(db *sql.DB, user_uuid string) ([]Post, error) {

	fetchUserPostsQuery := `
        SELECT p.*, u.username, u.profile_picture 
        FROM posts p
        JOIN users u ON p.user_uuid = u.user_uuid
        WHERE p.user_uuid = ?  -- Filtrer par user_uuid
        ORDER BY p.created_at DESC`

	rows, err := server.RunQuery(fetchUserPostsQuery, user_uuid)
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
