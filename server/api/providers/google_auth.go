package providers

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Remplacez par vos identifiants OAuth2 Google
const (
	redirectURL    = "http://localhost:8080/api/google_callback"
	oauthGoogleURL = "https://accounts.google.com/o/oauth2/auth"
	tokenURL       = "https://oauth2.googleapis.com/token"
	userInfoURL    = "https://www.googleapis.com/oauth2/v2/userinfo"
	oauthState     = "pseudo-random" // À sécuriser avec un état aléatoire en production
)

func LoadEnvVariables() error {
	envFile, err := os.Open("./.env")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening ENV file: %v\n", err)
	}

	reader := bufio.NewReader(envFile)

	for err == nil {
		var line []byte

		if line, err = reader.ReadBytes('\n'); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("reading Env variables: %v", err)
		}

		if err = os.Setenv(string(line)[:strings.Index(string(line), "=")], string(line)[strings.Index(string(line), "=")+1:]); err != nil {
			return fmt.Errorf("setting Env variables: %v", err)
		}

		// Debug
		//fmt.Printf("New env variable: %s = %s", string(line)[:strings.Index(string(line), "=")], os.Getenv(string(line)[:strings.Index(string(line), "=")]))
	}

	return nil
}

// Gestion du clic sur le bouton de connexion "Login with Google"
func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	// Construire l'URL d'authentification Google manuellement
	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s",
		oauthGoogleURL,
		os.Getenv("GOOGLE_ID"),
		redirectURL,
		"url:https://www.googleapis.com/auth/userinfo.email",
		oauthState,
	)

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// Gestion du callback après l'authentification Google
func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Vérifier que l'état correspond
	if r.URL.Query().Get("state") != oauthState {
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
	data.Set("redirect_uri", redirectURL)
	data.Set("grant_type", "authorization_code")

	// Faire la requête POST pour obtenir le token
	response, err := http.PostForm(tokenURL, data)
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
	req, err := http.NewRequest("GET", userInfoURL, nil)
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

// Structure pour stocker le token OAuth2
type OAuthToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	IdToken     string `json:"id_token"`
}
