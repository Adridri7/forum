package authentification

import (
	"encoding/json"
	"fmt"
	dbUser "forum/server/users"
	generator "forum/server/utils"
	"net/http"
	"os"
	"time"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if !r.Form.Has("new-username") || r.Method != "POST" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		fmt.Fprintln(os.Stderr, r.Form)
		return
	}

	w.Header().Set("Content-Type", "application/json")

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
			/*
				renderTemplate(w, "./static/authentification/authentification.html", map[string]interface{}{
					"Error": "User already exists with this email address. Please try with another.",
				})
			*/
			json.NewEncoder(w).Encode(map[string]interface{}{
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

	/*
		http.SetCookie(w, &http.Cookie{
			Name:   "UserLogged",
			Value:  newUser.ToCookieValue(),
			MaxAge: 300, // 5 minutes
		})
	*/

	//renderTemplate(w, "./static/homePage/index.html", nil)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"CookieValue": newUser.ToCookieValue(),
	})
}
