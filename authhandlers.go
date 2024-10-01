package main

import (
	"fmt"
	dbUser "forum/server/users"
	generator "forum/server/utils"
	"net/http"
	"os"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if !r.Form.Has("username") || r.Method != "POST" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		fmt.Fprintln(os.Stderr, r.Form)
		return
	}

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
		renderTemplate(w, "./static/authentification/authentification.html", map[string]interface{}{
			"Error": "No user was found with this email address. Please try another.",
		})
		return
	}

	// Yes but password invalid => return error
	if err = dbUser.CheckPassword(usrFound.EncryptedPassword, r.FormValue("password")); err != nil {
		renderTemplate(w, "./static/authentification/authentification.html", map[string]interface{}{
			"Error": "Password did not match. Please try another.",
		})
		return
	}

	fmt.Printf("User logged in: %s -> %s (%s)\n", usrFound.UUID, usrFound.Username, usrFound.Email)

	http.SetCookie(w, &http.Cookie{
		Name:   "UserLogged",
		Value:  usrFound.ToCookieValue(),
		MaxAge: 300, // 5 minutes
	})

	renderTemplate(w, "./static/homePage/index.html", nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if !r.Form.Has("new-username") || r.Method != "POST" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		fmt.Fprintln(os.Stderr, r.Form)
		return
	}

	var err error

	// Verify that user doesn't exist already
	{
		var usrFound dbUser.User

		usrFound, err = dbUser.FetchUserByEmail(r.FormValue("email"))
		if err != nil {
			http.Error(w, "Fatal error fetching", http.StatusInternalServerError)
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		if usrFound != (dbUser.User{}) {
			renderTemplate(w, "./static/authentification/authentification.html", map[string]interface{}{
				"Error": "User already exists with this email address. Please try with another.",
			})
			return
		}
	}

	// TODO : Profile picture

	var uuid string
	if uuid, err = generator.GenerateUUID(); err != nil {
		http.Error(w, "Fatal error gen", http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	newUser := dbUser.NewUser(uuid, r.FormValue("new-username"), r.FormValue("email"), r.FormValue("new-password"), time.Now(), "user", "")

	if err = dbUser.RegisterUser(newUser.ToMap()); err != nil {
		http.Error(w, "Fatal error add", http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	fmt.Printf("New user registered: %s -> %s (%s)\n", newUser.UUID, newUser.Username, newUser.Email)

	http.SetCookie(w, &http.Cookie{
		Name:   "UserLogged",
		Value:  newUser.ToCookieValue(),
		MaxAge: 300, // 5 minutes
	})

	renderTemplate(w, "./static/homePage/index.html", nil)
}
