# 📜 Forum

## ✍🏼 Authors
- [@Adrien](https://github.com/Adridri7/)
- [@Pierre](https://github.com/pcaboor/)
- [@Esteban](https://github.com/MrLepoischiche)
- [@Gabriel](https://github.com/Rookate)

## 📦 Documentation

# Forum Backend en Go avec SQLite

Ce projet implémente un backend basique de forum en utilisant Go (Golang) et une base de données SQLite. Il est structuré pour inclure une gestion des utilisateurs, des publications, et des commentaires. Ce document détaille la configuration de la base de données, l'initialisation du projet et l'exécution des requêtes SQL. Le tout est organisé en plusieurs fichiers.

## Table des Matières
1. [Installation](#installation)
2. [Structure du Projet](#structure-du-projet)
3. [Initialisation de la Base de Données (`initDB`)](#initialisation-de-la-base-de-données-initdb)
4. [Exécution des Requêtes SQL (`runQuery`)](#exécution-des-requêtes-sql-runquery)
5. [Exemples de Code](#exemples-de-code)
6. [Lancement du Serveur](#lancement-du-serveur)
7. [Contributions](#contributions)

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