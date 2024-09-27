package post

import (
	"encoding/json"
	"forum/server"
	"forum/server/posts"
	"net/http"
)

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer le post_uuid depuis les paramètres de la requête
	var params map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := posts.DeletePost(server.Db, params); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retourner une réponse succès
	w.WriteHeader(http.StatusNoContent) // No Content, car aucune donnée à renvoyer
}
