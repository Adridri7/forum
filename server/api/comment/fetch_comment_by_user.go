package comments

import (
	"encoding/json"
	"forum/server"
	"forum/server/comments"
	"net/http"
)

func FetchUserCommentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var params map[string]interface{}

	// Récupère la demande du front et décode le JSON pour obtenir user_uuid
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Vérification de la présence de user_uuid
	user_UUID, ok := params["user_uuid"].(string)
	if !ok || user_UUID == "" {
		http.Error(w, "Missing or invalid user_uuid", http.StatusBadRequest)
		return
	}

	// Récupérer les posts de l'utilisateur
	commentData, err := comments.FetchUserComments(server.Db, user_UUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Encoder les données des posts en JSON et les écrire dans la réponse
	if err := json.NewEncoder(w).Encode(commentData); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
