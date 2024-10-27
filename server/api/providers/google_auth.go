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
	GOOGLE_ID     = "54777063980-5g9u1tgobagtb6m7s60i50e64qn0v04t.apps.googleusercontent.com"
	GOOGLE_SECRET = "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"
)

type GoogleUser struct {
	Name    string `json:"name"`
	Picture string `json:"picture"`
	Email   string `json:"email"`
	//EmailVerified string `json:"verified_email"`
}

// Gestion du clic sur le bouton de connexion "Login with Google"
func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	// Construire l'URL d'authentification Google manuellement
	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s",
		oauthGoogleURL,
		GOOGLE_ID,
		redirectGoogleURL,
		"openid+email+profile", // Chaine de caractères contenant les scopes
		OAuthState,
	)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// Gestion du callback après l'authentification Google
func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
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
	token, err := getGoogleOauthToken(code)
	if err != nil {
		http.Error(w, "Failed to get Google Oauth token", http.StatusInternalServerError)
		return
	}

	// Utiliser le token pour récupérer les informations utilisateur
	userInfo, err := getGoogleUserInfo(token.AccessToken)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	// Décoder les infos utilisateur renvoyées...
	var googUsr GoogleUser
	if err = json.Unmarshal([]byte(userInfo), &googUsr); err != nil {
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	// ... puis on vérifie que l'utilisateur n'existe pas déjà
	var usr dbUser.User
	if usr, err = dbUser.FetchUserByEmail(googUsr.Email); err != nil {
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

		if strings.Contains(googUsr.Picture[8:], "lh3.googleusercontent.com/a/") {
			usr.ProfilePicture = dbUser.RandomProfilPicture()
		} else {
			usr.ProfilePicture = googUsr.Picture
		}

		if strings.IndexByte(googUsr.Name, '(') == -1 {
			usr.Username = googUsr.Name
		} else {
			usr.Username = googUsr.Name[:strings.IndexByte(googUsr.Name, '(')-1]
		}

		usr.Email = googUsr.Email
		usr.EncryptedPassword = ""

		if err = dbUser.RegisterUser(usr.ToMap()); err != nil {
			http.Error(w, "{\"Error\": \"Fatal error add\"}", http.StatusInternalServerError)
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
	} else {
		if strings.IndexByte(googUsr.Name, '(') == -1 {
			usr.Username = googUsr.Name
		} else {
			usr.Username = googUsr.Name[:strings.IndexByte(googUsr.Name, '(')-1]
		}

		if strings.Contains(googUsr.Picture[8:], "lh3.googleusercontent.com/a/") {
			usr.ProfilePicture = dbUser.RandomProfilPicture()
		} else {
			usr.ProfilePicture = googUsr.Picture
		}

		if err = usr.UpdateUser(map[string]interface{}{
			"email":           usr.Email,
			"password":        "",
			"profile_picture": usr.ProfilePicture,
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
func getGoogleOauthToken(code string) (*OAuthToken, error) {
	// Préparer la requête POST pour obtenir le token
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", GOOGLE_ID)
	data.Set("client_secret", GOOGLE_SECRET)
	data.Set("redirect_uri", redirectGoogleURL)
	data.Set("grant_type", "authorization_code")

	// Faire la requête POST pour obtenir le token
	response, err := http.PostForm(tokenGoogleURL, data)
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
func getGoogleUserInfo(accessToken string) (string, error) {
	// Faire une requête GET avec le token d'accès
	req, err := http.NewRequest("GET", userInfoGoogleURL, nil)
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
