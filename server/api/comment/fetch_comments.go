package comments

import (
	"encoding/json"
	"forum/server"
	"forum/server/comments"
	"net/http"
)

func FetchCommentHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var params map[string]interface{}

	// Récupere la demande du front et decode le JSON post_uuid
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Vérification de la présence de post_uuid
	postUUID, ok := params["post_uuid"].(string)
	if !ok || postUUID == "" {
		http.Error(w, "Missing or invalid post_uuid", http.StatusBadRequest)
		return
	}

	// Récupération des commentaires basés sur le post_uuid
	commentData, err := comments.FetchComment(server.Db, map[string]interface{}{
		"post_uuid": postUUID,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Si trouvé renvoie la réponse en format JSON pour le fronted etc....
	w.Header().Set("Content", "application/json")
	if err := json.NewEncoder(w).Encode(commentData); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
