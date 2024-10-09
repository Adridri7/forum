package server

import (
	"database/sql"
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

	log.Println("Connexion à la base de données réussie")
}
