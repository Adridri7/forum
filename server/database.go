package server

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Variable globale pour la base de données
var Db *sql.DB

// init() est appelé automatiquement avant le main() afin de vérifier la connexion à la db
func init() {
	var err error

	// Ouvre la connexion à la base de données ici
	// Le fichier forumdatabase.db est à la racine, c'est ce fichier qui contient toute la database
	// Grâce à l'extension sqlite de vscode, nous pouvons visualiser cela plus facilement
	// Clic droit sur forumdatabase.db
	// Open database
	// Magie on peut voir les tables avec les columns er rows

	// Db, err = sql.Open("sqlite3", "./forumdatabase.db")
	// if err != nil {
	// 	log.Fatalf("Erreur lors de l'ouverture de la base de données : %v", err)
	// }

	// // Vérifie la connexion
	// if err = Db.Ping(); err != nil {
	// 	log.Fatalf("Erreur lors de la connexion à la base de données : %v", err)
	// }
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		Db, err = sql.Open("sqlite3", "./forumdatabase.db")
		if err == nil {
			err = Db.Ping()
			if err == nil {
				break
			}
		}
		log.Printf("Tentative de connexion à la base de données échouée (%d/%d). Nouvelle tentative dans 5 secondes...", i+1, maxRetries)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatalf("Impossible de se connecter à la base de données après %d tentatives : %v", maxRetries, err)
	}
	createTables(Db)

	log.Println("Connexion à la base de données réussie")
}

func createTables(db *sql.DB) {
	// Requête pour créer la table users
	createUsersTable := `
    CREATE TABLE IF NOT EXISTS users (
        user_uuid TEXT PRIMARY KEY,
        username TEXT,
        email TEXT,
        password TEXT,
        role TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        profil_picture TEXT
    );`

	// Requête pour créer la table posts
	createPostsTable := `
    CREATE TABLE IF NOT EXISTS posts (
        post_uuid TEXT PRIMARY KEY,
        user_uuid TEXT,
        content TEXT,
        categories TEXT,
        likes INTEGER DEFAULT 0,
        dislikes INTEGER DEFAULT 0,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        post_image TEXT,
        FOREIGN KEY (user_uuid) REFERENCES users(user_uuid)
    );`

	// Requête pour créer la table comments
	createCommentsTable := `
    CREATE TABLE IF NOT EXISTS comments (
        comment_id TEXT PRIMARY KEY,
        post_uuid TEXT,
        user_uuid TEXT,
        content TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        likes INTEGER DEFAULT 0,
        dislikes INTEGER DEFAULT 0,
        FOREIGN KEY (post_uuid) REFERENCES posts(post_uuid),
        FOREIGN KEY (user_uuid) REFERENCES users(user_uuid)
    );`

	// Requête pour créer la table post_reactions
	createPostsReactionsTable := `
    CREATE TABLE IF NOT EXISTS post_reactions (
        post_uuid TEXT,
        user_uuid TEXT,
        action TEXT CHECK(action IN ('like', 'dislike')),
        PRIMARY KEY (post_uuid, user_uuid),
        FOREIGN KEY (post_uuid) REFERENCES posts(post_uuid),
        FOREIGN KEY (user_uuid) REFERENCES users(user_uuid)
    );`

	// Requête pour créer la table comment_reactions
	createCommentReactionsTable := `
    CREATE TABLE IF NOT EXISTS comment_reactions (
        comment_id TEXT,
        user_uuid TEXT,
        action TEXT CHECK(action IN ('like', 'dislike')),
        PRIMARY KEY (comment_id, user_uuid),
        FOREIGN KEY (comment_id) REFERENCES comments(comment_id),
        FOREIGN KEY (user_uuid) REFERENCES users(user_uuid)
    );`

	// Exécution des requêtes pour créer les tables
	statements := []struct {
		name      string
		statement string
	}{
		{"users", createUsersTable},
		{"posts", createPostsTable},
		{"comments", createCommentsTable},
		{"post_reactions", createPostsReactionsTable},
		{"comment_reactions", createCommentReactionsTable},
	}

	var createdTables []string

	for _, stmt := range statements {
		_, err := db.Exec(stmt.statement)
		if err != nil {
			log.Fatalf("Erreur lors de la création de la table %s: %v", stmt.name, err)
		}
		// Ajoute le nom de la table créée
		createdTables = append(createdTables, stmt.name)
	}

	if len(createdTables) > 0 {
		fmt.Printf("Tables créées avec succès : %s\n", createdTables)
	} else {
		fmt.Println("Aucune table n'a été créée.")
	}
}
