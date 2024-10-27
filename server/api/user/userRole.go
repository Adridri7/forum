package user

import (
	"encoding/json"
	"net/http"
)

func UpdateUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var params map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user_UUID, ok := params["user_uuid"].(string)
	if !ok || user_UUID == "" {
		http.Error(w, "Missing or invalid user_uuid", http.StatusBadRequest)
		return
	}

	action, ok := params["action"].(string)
	if !ok || action == "" {
		http.Error(w, "Missing or invalid action", http.StatusBadRequest)
		return
	}

	if err := UpdateUserRole(user_UUID, action); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

/*
	- L'idée serait dans 'Request' de faire une sorte d'interface où l'on pourrait écrire un message qui sera directement envoyer
	à l'admin et ce dernier pourra traiter la requête dans l'onglet 'Modération' !
*/
