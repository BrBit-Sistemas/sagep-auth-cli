package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

// Config contém as configurações do CLI
type Config struct {
	AuthURL    string
	AuthToken  string // JWT token (uso normal)
	AuthSecret string // Secret para HMAC (bootstrap)
}

// LoadConfig carrega a configuração a partir de flags, arquivo .env e variáveis de ambiente
// Ordem de precedência: flags > .env > env vars do sistema
// As variáveis SAGEP_AUTH_URL e (SAGEP_AUTH_SECRET ou SAGEP_AUTH_TOKEN) são obrigatórias
func LoadConfig(authURLFlag, authTokenFlag, authSecretFlag string) (*Config, error) {
	// Tentar carregar arquivo .env (se existir)
	// Primeiro tenta no diretório atual, depois procura a raiz do projeto
	envPath := ".env"
	if _, err := os.Stat(envPath); err != nil {
		// Se não encontrou no diretório atual, procura na raiz do projeto
		projectRoot, err := FindProjectRoot()
		if err == nil {
			rootEnvPath := filepath.Join(projectRoot, ".env")
			if _, err := os.Stat(rootEnvPath); err == nil {
				envPath = rootEnvPath
			}
		}
	}

	// Carregar .env se existir
	if _, err := os.Stat(envPath); err == nil {
		if err := godotenv.Load(envPath); err != nil {
			return nil, fmt.Errorf("erro ao carregar arquivo .env: %w", err)
		}
	}

	cfg := &Config{}

	// URL do auth: flag > .env > env vars do sistema
	if authURLFlag != "" {
		cfg.AuthURL = strings.TrimSuffix(authURLFlag, "/")
	} else {
		envURL := os.Getenv("SAGEP_AUTH_URL")
		if envURL == "" {
			return nil, fmt.Errorf("SAGEP_AUTH_URL é obrigatória. Configure via --url, arquivo .env ou variável de ambiente")
		}
		cfg.AuthURL = strings.TrimSuffix(envURL, "/")
	}

	// Token: flag > .env > env vars do sistema (opcional se tiver secret)
	if authTokenFlag != "" {
		cfg.AuthToken = authTokenFlag
	} else {
		cfg.AuthToken = os.Getenv("SAGEP_AUTH_TOKEN")
	}

	// Secret: flag > .env > env vars do sistema (opcional se tiver token)
	if authSecretFlag != "" {
		cfg.AuthSecret = authSecretFlag
	} else {
		cfg.AuthSecret = os.Getenv("SAGEP_AUTH_SECRET")
	}

	// Validar que tem pelo menos um (token OU secret)
	if cfg.AuthToken == "" && cfg.AuthSecret == "" {
		return nil, fmt.Errorf("SAGEP_AUTH_TOKEN ou SAGEP_AUTH_SECRET é obrigatório. Configure via --token/--secret, arquivo .env ou variável de ambiente")
	}

	return cfg, nil
}

// FindProjectRoot procura a raiz do projeto procurando por .env ou go.mod
func FindProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		// Verifica se existe .env ou go.mod neste diretório
		if _, err := os.Stat(filepath.Join(dir, ".env")); err == nil {
			return dir, nil
		}
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Chegou na raiz do sistema de arquivos
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("raiz do projeto não encontrada")
}

