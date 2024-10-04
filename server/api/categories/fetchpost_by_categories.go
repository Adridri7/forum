package categories

import (
	"encoding/json"
	"forum/server"
	"forum/server/posts"
	"net/http"
)

func FetchPostByCategoriesHandler(w http.ResponseWriter, r *http.Request) {

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
	category, ok := params["categories"].(string)
	if !ok || category == "" {
		http.Error(w, "Missing or invalid categories", http.StatusBadRequest)
		return
	}

	// Récupération des commentaires basés sur le post_uuid
	categoriesData, err := posts.FetchPostsByCategory(server.Db, category)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Si trouvé renvoie la réponse en format JSON pour le fronted etc....
	w.Header().Set("Content", "application/json")
	if err := json.NewEncoder(w).Encode(categoriesData); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
