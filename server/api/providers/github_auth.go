package providers

import (
	//dbUser "forum/server/api/user"
	//utils "forum/server/utils"

	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	//"os"
)

const (
	GITHUB_ID     = "Ov23libnZpblLJfYUNV0"
	GITHUB_SECRET = "ddb158a0782ac5c6d9b0141784df6f0c9b2d33de"
)

// Gestion du clic sur le bouton de connexion "Login with Github"
func HandleGithubLogin(w http.ResponseWriter, r *http.Request) {
	// Construire l'URL d'authentification Github manuellement
	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s",
		oauthGithubURL,
		GITHUB_ID,
		redirectGithubURL,
		"user:email",
		OAuthState,
	)

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// Gestion du callback après l'authentification Github
func HandleGithubCallback(w http.ResponseWriter, r *http.Request) {
	// Vérifier que l'état correspond
	if r.URL.Query().Get("state") != OAuthState {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	// Récupérer le code d'autorisation
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No code in URL", http.StatusBadRequest)
		return
	}

	// Échange du code contre un token d'accès
	token, err := getGithubOauthToken(code)
	if err != nil {
		http.Error(w, "Failed to get Github Oauth token", http.StatusInternalServerError)
		return
	}

	// Utiliser le token pour récupérer les informations utilisateur
	userInfo, err := getGithubUserInfo(token.AccessToken)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(userInfo))

	/*
		// Décoder les infos utilisateur renvoyées...
		var githUsr GithubUser
		if err = json.Unmarshal([]byte(userInfo), &githUsr); err != nil {
			http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
			return
		}

		// ... puis on vérifie que l'utilisateur n'existe pas déjà
		var usr dbUser.User
		if usr, err = dbUser.FetchUserByEmail(githUsr.Email); err != nil {
			http.Error(w, "{\"Error\": \"Fatal error fetching\"}", http.StatusInternalServerError)
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		if usr == (dbUser.User{}) {
			if usr.UUID, err = utils.GenerateUUID(); err != nil {
				http.Error(w, "{\"Error\": \"Fatal error gen\"}", http.StatusInternalServerError)
				fmt.Fprintln(os.Stderr, err.Error())
				return
			}

			if strings.Contains(githUsr.Picture[8:], "lh3.googleusercontent.com/a/") {
				usr.ProfilePicture = dbUser.RandomProfilPicture()
			} else {
				usr.ProfilePicture = githUsr.Picture
			}

			if strings.IndexByte(githUsr.Name, '(') == -1 {
				usr.Username = githUsr.Name
			} else {
				usr.Username = githUsr.Name[:strings.IndexByte(githUsr.Name, '(')-1]
			}

			usr.Email = githUsr.Email
			usr.EncryptedPassword = ""
			usr.Role = "user"

			if err = dbUser.RegisterUser(usr.ToMap()); err != nil {
				http.Error(w, "{\"Error\": \"Fatal error add\"}", http.StatusInternalServerError)
				fmt.Fprintln(os.Stderr, err.Error())
				return
			}
		} else {
			if strings.IndexByte(githUsr.Name, '(') == -1 {
				usr.Username = githUsr.Name
			} else {
				usr.Username = githUsr.Name[:strings.IndexByte(githUsr.Name, '(')-1]
			}

			if strings.Contains(githUsr.Picture[8:], "lh3.googleusercontent.com/a/") {
				usr.ProfilePicture = dbUser.RandomProfilPicture()
			} else {
				usr.ProfilePicture = githUsr.Picture
			}

			if err = usr.UpdateUser(map[string]interface{}{
				"email":           usr.Email,
				"password":        nil,
				"profile_picture": usr.ProfilePicture,
				"role":            "user",
				"username":        usr.Username,
				"user_uuid":       usr.UUID,
			}); err != nil {
				http.Error(w, "{\"Error\": \"Fatal error update\"}", http.StatusInternalServerError)
				fmt.Fprintln(os.Stderr, err.Error())
				return
			}
		}

		http.SetCookie(w, &http.Cookie{
			Name:   "UserLogged",
			Value:  usr.ToCookieValue(),
			Path:   "/",
			MaxAge: 300, // 5 minutes
		})

		fmt.Printf("User logged in: %s -> %s (%s)\n", usr.UUID, usr.Username, usr.Email)

		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	*/
}

// Récupère le token OAuth2 en échangeant le code d'autorisation
func getGithubOauthToken(code string) (*OAuthToken, error) {
	// Préparer la requête POST pour obtenir le token
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", GITHUB_ID)
	data.Set("client_secret", GITHUB_SECRET)
	data.Set("redirect_uri", redirectGithubURL)
	data.Set("grant_type", "authorization_code")

	fmt.Println(data)

	// Faire la requête POST pour obtenir le token
	response, err := http.PostForm(tokenGithubURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to request token: %v", err)
	}
	defer response.Body.Close()

	// Lire la réponse et décoder le token
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read token response: %v", err)
	}

	var token OAuthToken
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token response: %v", err)
	}

	return &token, nil
}

// Utilise le token d'accès pour récupérer les informations utilisateur
func getGithubUserInfo(accessToken string) (string, error) {
	fmt.Println(accessToken)

	// Faire une requête GET avec le token d'accès
	req, err := http.NewRequest("GET", userInfoGithubURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get user info: %v", err)
	}
	defer response.Body.Close()

	// Lire la réponse et extraire les informations utilisateur
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read user info: %v", err)
	}

	return string(body), nil
}
