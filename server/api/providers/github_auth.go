package providers

import (
	authentification "forum/server/api/login"
	dbUser "forum/server/api/user"
	utils "forum/server/utils"

	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	GITHUB_ID     = "Ov23libnZpblLJfYUNV0"
	GITHUB_SECRET = "ddb158a0782ac5c6d9b0141784df6f0c9b2d33de"
)

type GithubUser struct {
	Username  string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
}

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
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// Utiliser le token pour récupérer les informations utilisateur
	userInfo, err := getGithubUserInfo(token.AccessToken)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// Décoder les infos utilisateur renvoyées...
	var githUsr GithubUser
	if err = json.Unmarshal([]byte(userInfo), &githUsr); err != nil {
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	// ... puis on vérifie que l'utilisateur n'existe pas déjà
	var usr dbUser.User
	if githUsr.Email != "" {
		if usr, err = dbUser.FetchUserByEmail(githUsr.Email); err != nil {
			http.Error(w, "{\"Error\": \"Fatal error fetching\"}", http.StatusInternalServerError)
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
	}

	if usr == (dbUser.User{}) {
		if usr.UUID, err = utils.GenerateUUID(); err != nil {
			http.Error(w, "{\"Error\": \"Fatal error gen\"}", http.StatusInternalServerError)
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		usr.Username = githUsr.Username
		usr.Email = githUsr.Email
		usr.EncryptedPassword = ""
		usr.Role = "user"
		usr.ProfilePicture = githUsr.AvatarURL

		if err = dbUser.RegisterUser(usr.ToMap()); err != nil {
			http.Error(w, "{\"Error\": \"Fatal error add\"}", http.StatusInternalServerError)
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
	} else {
		usr.Username = githUsr.Username
		usr.ProfilePicture = githUsr.AvatarURL

		if err = usr.UpdateUser(map[string]interface{}{
			"email":           usr.Email,
			"password":        "",
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

	sessionID, _ := utils.GenerateUUID() // Génère un UUID unique
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  sessionID,
		Path:   "/",
		MaxAge: 3600,
	})
	authentification.Sessions[sessionID] = usr

	http.Redirect(w, r, "/", http.StatusPermanentRedirect)

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
	if err := json.Unmarshal(jsonified(string(body)), &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token response: %v", err)
	}

	return &token, nil
}

// Utilise le token d'accès pour récupérer les informations utilisateur
func getGithubUserInfo(accessToken string) (string, error) {
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

func jsonified(s string) []byte {
	res := make([]byte, 0)
	res = append(res, '{')

	for _, field := range strings.Split(s, "&") {
		key, value := strings.Split(field, "=")[0], strings.Split(field, "=")[1]
		res = append(res, []byte("\""+key+"\": \""+value+"\", ")...)
	}

	res = res[:len(res)-2]
	res = append(res, '}')

	return res
}
