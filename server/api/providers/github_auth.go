package providers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	redirectGithubURL = "http://localhost:8080/api/github_callback"
	//oauthGithubURL    = "https://accounts.google.com/o/oauth2/auth"
	//tokenGithubURL    = "https://github.com/login/oauth/access_token"
	//userInfoGithubURL = "https://api.github.com/user"
)

// Gestion du clic sur le bouton de connexion "Login with Github"
func HandleGithubLogin(w http.ResponseWriter, r *http.Request) {
	// Construire l'URL d'authentification Github manuellement
	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s",
		oauthGithubURL,
		os.Getenv("GITHUB_ID"),
		redirectGithubURL,
		"user:email",
		OAuthState,
	)

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// Gestion du callback après l'authentification Github
func HandleGithubCallback(w http.ResponseWriter, r *http.Request) {

}

func handleGithubLogin(w http.ResponseWriter, r *http.Request) {
	url := githubOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGithubCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != oauthStateString {
		log.Println("Invalid OAuth state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := githubOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Println("Code exchange failed: ", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Utilisez le token pour récupérer les informations de l'utilisateur
	client := githubOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		log.Println("Failed to get user info: ", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	defer resp.Body.Close()
	fmt.Fprintf(w, "GitHub login successful!")
}
