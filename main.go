package main

import (
	"fmt"
	"forum/server/api"
	"log"
	"net/http"
)

// Pour l'instant nous testons uniquement des requetes simples

func main() {
	/*
		// Exemple de requête pour créer une table `users`
		createTableQuery := `
		CREATE TABLE IF NOT EXISTS users (
			user_uuid TEXT PRIMARY KEY,
			username TEXT,
			email TEXT,
			password TEXT,
			created_at TEXT
		);`

		// Exécute la requête pour créer la table
		_, err := server.RunQuery(createTableQuery)
		if err != nil {
			log.Fatalf("Erreur lors de la création de la table : %v", err)
		}

		// Insère un utilisateur dans la table
		insertQuery := `
		INSERT INTO users (user_uuid, username, email, password, created_at)
		VALUES ('123e4567-e89b-12d3-a456-426614174000', 'JohnDoe', 'john@example.com', 'password123', '2024-09-25');`

		_, err = server.RunQuery(insertQuery)
		if err != nil {
			log.Fatalf("Erreur lors de l'insertion de l'utilisateur : %v", err)
		}

		// Sélectionne tous les utilisateurs
		selectQuery := "SELECT * FROM users"
		users, err := server.RunQuery(selectQuery)
		if err != nil {
			log.Fatalf("Erreur lors de la récupération des utilisateurs : %v", err)
		}

		// Affiche les utilisateurs
		for _, user := range users {
			fmt.Printf("Utilisateur : %v\n", user)
		}*/
	http.HandleFunc("/api/createPost", api.CreatePostHandler)

	http.Handle("/", http.FileServer(http.Dir("./static")))

	fmt.Println("Serveur démarré sur le port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
