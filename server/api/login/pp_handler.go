package authentification

import (
	"encoding/json"
	"fmt"
	dbUser "forum/server/api/user"
	"net/http"
	"os"
)

func PP_Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var usr dbUser.User

	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		http.Error(w, "Fatal error decode id", http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if usr.ProfilePicture, err = dbUser.FetchPPByID(usr.UUID); err != nil {
		http.Error(w, "Fatal error query pp", http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if usr.ProfilePicture == "" {
		fmt.Printf("User not found for ID \"%s\"\n", usr.UUID)
	}
	json.NewEncoder(w).Encode(usr.ProfilePicture)
}
