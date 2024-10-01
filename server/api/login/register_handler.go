package authentification

import (
	"fmt"
	dbUser "forum/server/users"
	generator "forum/server/utils"
	"net/http"
	"os"
	"time"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if !r.Form.Has("new-username") || r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
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
			http.Error(w, "User already exists with this email address. Please try with another.", http.StatusUnauthorized)
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

	w.WriteHeader(http.StatusOK)
}
