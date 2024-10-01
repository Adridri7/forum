package authentification

import (
	"encoding/json"
	"fmt"
	dbUser "forum/server/users"
	"io"
	"net/http"
	"os"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		fmt.Fprintln(os.Stderr, r.Form)
		return
	}

	var reqBody []byte
	var usr dbUser.User
	var err error

	if reqBody, err = io.ReadAll(r.Body); err != nil {
		http.Error(w, "Fatal error body", http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	if len(reqBody) == 0 {
		http.Error(w, "Body empty", http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(reqBody, &usr); err != nil {
		http.Error(w, "Fatal error marshal", http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	{
		var usrFound dbUser.User

		// Does user exist?
		usrFound, err = dbUser.FetchUserByEmail(usr.Email)
		if err != nil {
			http.Error(w, "Fatal error fetching", http.StatusInternalServerError)
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		// No => return error
		if usrFound == (dbUser.User{}) {
			http.Error(w, "No user was found with this email address. Please try another.", http.StatusNotFound)
			return
		}

		// Yes but password invalid => return error
		if err = dbUser.CheckPassword(usrFound.EncryptedPassword, usr.EncryptedPassword); err != nil {
			http.Error(w, "Password did not match. Please try another.", http.StatusUnauthorized)
			return
		}

		usr = usrFound
	}

	fmt.Printf("User logged in: %s -> %s (%s)\n", usr.UUID, usr.Username, usr.Email)

	http.SetCookie(w, &http.Cookie{
		Name:   "UserLogged",
		Value:  usr.ToCookieValue(),
		MaxAge: 300, // 5 minutes
	})

	w.WriteHeader(http.StatusOK)
}
