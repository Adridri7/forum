package providers

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	OAuthState = "pseudo-random" // À sécuriser avec un état aléatoire en production

	redirectDiscordURL    = "https://localhost:8080/api/discord_callback"
	oauthDiscordURL       = "https://discord.com/oauth2/authorize"
	tokenDiscordURL       = "https://discord.com/api/v10/oauth2/token"
	tokenRevokeDiscordURL = "https://discord.com/api/v10/oauth2/token/revoke"
	userInfoDiscordURL    = "https://discord.com/api/v10/users/@me"
	getPPDiscordURL       = "https://cdn.discordapp.com/avatars/"

	redirectGithubURL = "https://localhost:8080/api/github_callback"
	oauthGithubURL    = "https://github.com/login/oauth/authorize"
	tokenGithubURL    = "https://github.com/login/oauth/access_token"
	userInfoGithubURL = "https://api.github.com/user"

	redirectGoogleURL = "https://localhost:8080/api/google_callback"
	oauthGoogleURL    = "https://accounts.google.com/o/oauth2/auth"
	tokenGoogleURL    = "https://oauth2.googleapis.com/token"
	userInfoGoogleURL = "https://www.googleapis.com/oauth2/v2/userinfo"
)

// Structure pour stocker le token OAuth2
type OAuthToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	IdToken     string `json:"id_token"`
	Scope       string `json:"scope"`
}

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

		line = []byte(strings.TrimSpace(string(line)))

		if err = os.Setenv(string(line)[:strings.Index(string(line), "=")], string(line)[strings.Index(string(line), "=")+1:]); err != nil {
			return fmt.Errorf("setting Env variables: %v", err)
		}

		// Debug
		//fmt.Printf("New env variable: %s = %s\n", string(line)[:strings.Index(string(line), "=")], os.Getenv(string(line)[:strings.Index(string(line), "=")]))
	}

	return nil
}
