package authentification

import (
	"encoding/json"
	"fmt"
	dbUser "forum/server/users"
	"net/http"
	"os"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if !r.Form.Has("username") || r.Method != "POST" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		fmt.Fprintln(os.Stderr, r.Form)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var usrFound dbUser.User
	var err error

	// Does user exist?
	usrFound, err = dbUser.FetchUserByEmail(r.FormValue("email"))
	if err != nil {
		http.Error(w, "Fatal error fetching", http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	// No => return error
	if usrFound == (dbUser.User{}) {
		/*
			renderTemplate(w, "./static/authentification/authentification.html", map[string]interface{}{
				"Error": "No user was found with this email address. Please try another.",
			})
		*/
		json.NewEncoder(w).Encode(map[string]interface{}{
			"Error": "No user was found with this email address. Please try another.",
		})
		return
	}

	// Yes but password invalid => return error
	if err = dbUser.CheckPassword(usrFound.EncryptedPassword, r.FormValue("password")); err != nil {
		/*
			renderTemplate(w, "./static/authentification/authentification.html", map[string]interface{}{
				"Error": "Password did not match. Please try another.",
			})
		*/
		json.NewEncoder(w).Encode(map[string]interface{}{
			"Error": "Password did not match. Please try another.",
		})
		return
	}

	fmt.Printf("User logged in: %s -> %s (%s)\n", usrFound.UUID, usrFound.Username, usrFound.Email)

	/*
		http.SetCookie(w, &http.Cookie{
			Name:   "UserLogged",
			Value:  usrFound.ToCookieValue(),
			MaxAge: 300, // 5 minutes
		})
	*/

	//renderTemplate(w, "./static/homePage/index.html", nil)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"CookieValue": usrFound.ToCookieValue(),
	})
}
