package requests

import (
	"encoding/json"
	"fmt"
	"forum/server"
	request "forum/server/admin_requests"

	"net/http"
)

type ActionRequest struct {
	RequestUUID string `json:"request_uuid"`
	Action      string `json:"action"`
}

func HandleActionRequestAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var actionReq ActionRequest
	err := json.NewDecoder(r.Body).Decode(&actionReq)
	if err != nil {
		http.Error(w, "Données de requête invalides", http.StatusBadRequest)
		return
	}

	// Valider l'action pour être sûr que c'est soit "valid" soit "reject"
	if actionReq.Action != "approuve" && actionReq.Action != "reject" {
		http.Error(w, "Action invalide", http.StatusBadRequest)
		return
	}

	// Appeler la fonction pour gérer l'action en base de données
	err = request.HandleActionRequest(server.Db, actionReq.RequestUUID, actionReq.Action)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors du traitement de la requête : %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Action mise à jour avec succès"))
}
