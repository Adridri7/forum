package providers

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// Structure pour stocker le token OAuth2
type OAuthToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	IdToken     string `json:"id_token"`
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

		if err = os.Setenv(string(line)[:strings.Index(string(line), "=")], string(line)[strings.Index(string(line), "=")+1:]); err != nil {
			return fmt.Errorf("setting Env variables: %v", err)
		}

		// Debug
		//fmt.Printf("New env variable: %s = %s", string(line)[:strings.Index(string(line), "=")], os.Getenv(string(line)[:strings.Index(string(line), "=")]))
	}

	return nil
}
