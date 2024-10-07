package providers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

const (
	redirectDiscordURL    = "http%3A%2F%2Flocalhost%3A8080%2Fapi%2Fdiscord_callback"
	oauthDiscordURL       = "https://discord.com/oauth2/authorize"
	tokenDiscordURL       = "https://discord.com/api/oauth2/token"
	tokenRevokeDiscordURL = "https://discord.com/api/oauth2/token/revoke"
)

// Gestion du clic sur le bouton de connexion "Login with Discord"
func HandleDiscordLogin(w http.ResponseWriter, r *http.Request) {
	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s",
		oauthDiscordURL,
		os.Getenv("DISCORD_ID"),
		redirectDiscordURL,
		"identify+email",
		OAuthState,
	)

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// Gestion du callback après l'authentification Discord
func HandleDiscordCallback(w http.ResponseWriter, r *http.Request) {

}

// Récupère le token OAuth2 en échangeant le code d'autorisation
func getDiscordOauthToken(code string) (*OAuthToken, error) {
	// Préparer la requête POST pour obtenir le token
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", os.Getenv("DISCORD_ID"))
	data.Set("client_secret", os.Getenv("DISCORD_SECRET"))
	data.Set("redirect_uri", redirectDiscordURL)
	data.Set("grant_type", "authorization_code")

	// Faire la requête POST pour obtenir le token
	response, err := http.PostForm(tokenDiscordURL, data)
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
