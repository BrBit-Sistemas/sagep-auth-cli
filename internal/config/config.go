package config

import (
	"fmt"
	"os"
	"strings"
)

// Config contém as configurações do CLI
type Config struct {
	AuthURL   string
	AuthToken string
}

// LoadConfig carrega a configuração a partir de flags e variáveis de ambiente
// Ordem de precedência: flags > env vars > defaults
func LoadConfig(authURLFlag, authTokenFlag string) (*Config, error) {
	cfg := &Config{}

	// URL do auth: flag > env > default
	if authURLFlag != "" {
		cfg.AuthURL = strings.TrimSuffix(authURLFlag, "/")
	} else if envURL := os.Getenv("SAGEP_AUTH_URL"); envURL != "" {
		cfg.AuthURL = strings.TrimSuffix(envURL, "/")
	} else {
		return nil, fmt.Errorf("URL do auth não configurada. Use --url ou SAGEP_AUTH_URL")
	}

	// Token: flag > env
	if authTokenFlag != "" {
		cfg.AuthToken = authTokenFlag
	} else if envToken := os.Getenv("SAGEP_AUTH_TOKEN"); envToken != "" {
		cfg.AuthToken = envToken
	} else {
		return nil, fmt.Errorf("token de autenticação não configurado. Use --token ou SAGEP_AUTH_TOKEN")
	}

	return cfg, nil
}

