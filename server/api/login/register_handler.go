package authentification

import (
	"encoding/json"
	"fmt"
	dbUser "forum/server/api/user"
	generator "forum/server/utils"
	"io"
	"net/http"
	"os"
	"time"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, "{\"Error\": \"Method not allowed\"}", http.StatusBadRequest)
		fmt.Fprintln(os.Stderr, r.Form)
		return
	}

	var reqBody []byte
	var newUser dbUser.User
	var err error

	if reqBody, err = io.ReadAll(r.Body); err != nil {
		http.Error(w, "{\"Error\": \"Fatal error body\"}", http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	if len(reqBody) == 0 {
		http.Error(w, "{\"Error\": \"Body empty\"}", http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(reqBody, &newUser); err != nil {
		http.Error(w, "{\"Error\": \"Fatal error marshal\"}", http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	// Verify that user doesn't exist already
	{
		var usrFound dbUser.User

		usrFound, err = dbUser.FetchUserByEmail(newUser.Email)
		if err != nil {
			http.Error(w, "{\"Error\": \"Fatal error fetching\"}", http.StatusInternalServerError)
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		if usrFound != (dbUser.User{}) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict) // Utilisation de 409 Conflict

			// Crée une structure d'erreur pour le JSON
			errorResponse := map[string]string{
				"error": "User already exists with this email address. Please try with another.",
			}

			// Encode la réponse en JSON
			if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
				http.Error(w, "{\"error\": \"Failed to encode error response.\"}", http.StatusInternalServerError)
				return
			}

			return
		}

		UsernameFound, _ := dbUser.IsUsernameTaken(newUser.Username)
		fmt.Println(UsernameFound)
		if UsernameFound {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict) // Utilisation de 409 Conflict

			// Crée une structure d'erreur pour le JSON
			errorResponse := map[string]string{
				"error": "User already exists with this Username. Please try with another.",
			}

			// Encode la réponse en JSON
			if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
				http.Error(w, "{\"error\": \"Failed to encode error response.\"}", http.StatusInternalServerError)
				return
			}

			return
		}
	}

	var uuid string

	if uuid, err = generator.GenerateUUID(); err != nil {
		http.Error(w, "{\"Error\": \"Fatal error gen\"}", http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	if newUser.EncryptedPassword, err = dbUser.HashPassword(newUser.EncryptedPassword); err != nil {
		http.Error(w, "{\"Error\": \"Fatal error hash\"}", http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	newUser.UUID = uuid
	newUser.CreatedAt = time.Now()
	newUser.Role = "user"

	if err = dbUser.RegisterUser(newUser.ToMap()); err != nil {
		http.Error(w, "{\"Error\": \"Fatal error add\"}", http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	fmt.Printf("New user registered: %s -> %s (%s)\n", newUser.UUID, newUser.Username, newUser.Email)

	sessionID, _ := generator.GenerateUUID() // Génère un UUID unique
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  sessionID,
		Path:   "/",
		MaxAge: 3600,
	})
	Sessions[sessionID] = newUser

	w.WriteHeader(http.StatusOK)
}
