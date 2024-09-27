package server

import (
	"log"
)

// runQuery exécute une requête SQL avec des paramètres et renvoie les résultats

/*--------------------------------------------------------------------------------------------------------*
|     																									  |
| - variable query = représente les requetes sql														  |
| - params ...interface{} : le second paramètre est un argument de type variadique (...interface{}),      |
| ce qui signifie que l'on peut passer un nombre quelconque d'arguments de différents types (interface{}  |
| peut être n'importe quel type en Go). Cela permet de fournir dynamiquement les valeurs pour les ?       |
| de la requête SQL (SELECT * FROM users WHERE id = ?).													  |
|																										  |
| []map[string]interface{} : la fonction retourne une slice ([]) de maps (map[string]interface{}),		  |
| où chaque map représente une row de la table avec des columns et leurs valeurs.					      |
|																										  |
*---------------------------------------------------------------------------------------------------------*/

func RunQuery(query string, params ...interface{}) ([]map[string]interface{}, error) {

	//----------------------------------------------------------------------//
	// Prépare la requête
	// Voir ./database.go "var db *sql.DB"
	// params ex: "SELECT * FROM users"
	//----------------------------------------------------------------------//

	rows, err := db.Query(query, params...)
	if err != nil {
		log.Printf("Erreur lors de l'exécution de la requête : %v", err)
		return nil, err
	}
	defer rows.Close()

	// Récupère les colonnes de la requête
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Création d'une slice pour stocker les valeurs

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// Stockage des résultats dans une liste de maps
	// Tous les éléments trouvés sont stockés et renvoyés

	var results []map[string]interface{}
	for rows.Next() {
		// Remplit les valeurs pour la ligne actuelle
		if err := rows.Scan(valuePtrs...); err != nil {
			log.Printf("Erreur lors du scan des résultats : %v", err)
			return nil, err
		}

		// Crée une map pour la ligne
		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			row[col] = val
		}
		results = append(results, row)
	}

	return results, nil
}
