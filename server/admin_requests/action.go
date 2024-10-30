package request

import (
	"database/sql"
	"fmt"
	"forum/server"
)

func HandleActionRequest(db *sql.DB, request_uuid string, action string) error {
	// Vérifier si l'utilisateur a déjà interagi
	checkQuery := `SELECT action FROM admin_requests WHERE request_uuid = ?`
	_, err := server.RunQuery(checkQuery, request_uuid)
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification de l'action: %v", err)
	}

	// Mettre à jour la colonne `action` avec "approuve" ou "reject"
	updateQuery := `UPDATE admin_requests SET action = ? WHERE request_uuid = ?`
	_, err = server.RunQuery(updateQuery, action, request_uuid)
	if err != nil {
		return fmt.Errorf("erreur lors de la mise à jour de l'action: %v", err)
	}

	return nil
}
