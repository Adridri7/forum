package providers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// Gestion du clic sur le bouton de connexion "Login with Google"
func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	// Construire l'URL d'authentification Google manuellement
	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s",
		oauthGoogleURL,
		os.Getenv("GOOGLE_ID"),
		redirectGoogleURL,
		"openid email profile", // Chaine de caractères contenant les scopes
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

	// Afficher les informations utilisateur
	fmt.Fprintf(w, "User Info: %s", userInfo)
}

// Récupère le token OAuth2 en échangeant le code d'autorisation
func getGoogleOauthToken(code string) (*OAuthToken, error) {
	// Préparer la requête POST pour obtenir le token
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", os.Getenv("GOOGLE_ID"))
	data.Set("client_secret", os.Getenv("GOOGLE_SECRET"))
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
