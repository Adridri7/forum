package authentification

import (
	"fmt"
	dbUser "forum/server/users"
	"net/http"
	"os"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if !r.Form.Has("username") || r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
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
		http.Error(w, "No user was found with this email address. Please try another.", http.StatusNotFound)
		return
	}

	// Yes but password invalid => return error
	if err = dbUser.CheckPassword(usrFound.EncryptedPassword, r.FormValue("password")); err != nil {
		http.Error(w, "Password did not match. Please try another.", http.StatusUnauthorized)
		return
	}

	fmt.Printf("User logged in: %s -> %s (%s)\n", usrFound.UUID, usrFound.Username, usrFound.Email)

	http.SetCookie(w, &http.Cookie{
		Name:   "UserLogged",
		Value:  usrFound.ToCookieValue(),
		MaxAge: 300, // 5 minutes
	})

	w.WriteHeader(http.StatusOK)
}
