# 📜 Forum

## ✍🏼 Authors
- [@Adrien](https://github.com/Adridri7/)
- [@Pierre](https://github.com/pcaboor/)
- [@Esteban](https://github.com/MrLepoischiche)
- [@Gabriel](https://github.com/Rookate)

## 📦 Documentation
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- [UUID](https://github.com/gofrs/uuid?tab=readme-ov-file)
- [Sqlite3](https://github.com/mattn/go-sqlite3?tab=readme-ov-file)


# Forum Backend en Go avec SQLite

Ce projet implémente un backend basique de forum en utilisant Go (Golang) et une base de données SQLite. Il est structuré pour inclure une gestion des utilisateurs, des publications, et des commentaires. Ce document détaille la configuration de la base de données, l'initialisation du projet et l'exécution des requêtes SQL. Le tout est organisé en plusieurs fichiers.

## Table des Matières
1. [Installation](#installation)
2. [Structure du Projet](#structure-du-projet)
3. [Initialisation de la Base de Données (`initDB`)](#initialisation-de-la-base-de-données-initdb)
4. [Exécution des Requêtes SQL (`runQuery`)](#exécution-des-requêtes-sql-runquery)
5. [Exemples de Code](#exemples-de-code)
6. [Lancement du Serveur](#lancement-du-serveur)

## Installation

### Pré-requis
- [Go](https://golang.org/dl/) (version 1.16 ou supérieure)
- SQLite installé (optionnel si vous utilisez une base de données existante)
- Git installé pour cloner le projet

### Instructions d'installation

1. Clonez ce dépôt GitHub :

    ```bash
    git clone https://github.com/Adridri7/forum.git
    cd forum
    ```

2. Initialisez le projet Go et les modules :

    ```bash
    go mod init forum
    go mod tidy
    ```

3. Ajoutez le support pour SQLite :

    ```bash
    go get -u github.com/mattn/go-sqlite3
    ```

## Structure du Projet




## Initialisation de la Base de Données (`initDB`)

La fonction `initDB` est utilisée pour créer et initialiser la connexion à la base de données SQLite. Elle configure la base de données, vérifie la connexion et initialise les tables nécessaires.

### Exemple de code

```go
package server

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Variable globale pour stocker la connexion à la base de données
var db *sql.DB

// initDB initialise la connexion à la base de données SQLite
func initDB(dataSourceName string) error {
	var err error
	db, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return err
	}

	// Vérification de la connexion
	if err = db.Ping(); err != nil {
		return err
	}

	log.Println("Connexion à la base de données SQLite réussie !")
	return nil
}
```

### Explication

- sql.Open("sqlite3", dataSourceName) : Ouvre une connexion avec le  pilote sqlite3.
- db.Ping() : Vérifie que la connexion est active.
Si la connexion est établie, un message est affiché dans la console.

# Exécution des requêtes sql runquery

```go
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
```

# Exemples de code

Exemple d'une structure du serveur 

```golang
func main() {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/api/login", authentification.LoginHandler)
	mux.HandleFunc("/api/get-pp", authentification.PP_handler)


	server := http.Server{
		Addr:              ":8080",
		Handler:           mux,
		MaxHeaderBytes:    1 << 26, // 4 MB
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      45 * time.Second,
		IdleTimeout:       3 * time.Minute,
	}

	log.Println("Server started on http://localhost:8080")

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start server: %s", err)
	}
}
```

Voici un exemple d'un handler

```golang
package authentification

import (
	"encoding/json"
	"fmt"
	dbUser "forum/server/api/user"
	"net/http"
	"os"
)

func PP_Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var usr dbUser.User

	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		http.Error(w, "Fatal error decode id", http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if usr.ProfilePicture, err = dbUser.FetchPPByID(usr.UUID); err != nil {
		http.Error(w, "Fatal error query pp", http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if usr.ProfilePicture == "" {
		fmt.Printf("User not found for ID \"%s\"\n", usr.UUID)
	}
	json.NewEncoder(w).Encode(usr.ProfilePicture)
}

```

# Lancement du serveur

```sh
cd ~/forum
go run .
```