package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/BrBit-Sistemas/sagep-auth-cli/internal/client"
	"github.com/BrBit-Sistemas/sagep-auth-cli/internal/config"
	"github.com/BrBit-Sistemas/sagep-auth-cli/internal/manifest"
)

// RunSync executa o comando de sincronização
func RunSync(manifestPath string, cfg *config.Config) error {
	// Carregar manifest
	m, err := manifest.LoadManifest(manifestPath)
	if err != nil {
		return fmt.Errorf("erro ao carregar manifest: %w", err)
	}

	// Criar cliente
	authClient := client.NewAuthClient(cfg.AuthURL, cfg.AuthToken, cfg.AuthSecret)

	// Exibir informações iniciais
	fmt.Printf("Sincronizando aplicação: %s\n", m.Application.Code)
	fmt.Printf("URL do auth: %s\n\n", cfg.AuthURL)

	// Executar sync
	ctx := context.Background()
	resp, err := authClient.SyncApplication(ctx, m)
	if err != nil {
		return fmt.Errorf("erro ao sincronizar: %w", err)
	}

	// Calcular estatísticas
	permsCreated := 0
	permsUpdated := 0
	for _, perm := range resp.Permissions {
		if perm.Action == "created" {
			permsCreated++
		} else if perm.Action == "updated" {
			permsUpdated++
		}
	}

	rolesCreated := 0
	rolesUpdated := 0
	for _, role := range resp.Roles {
		if role.Action == "created" {
			rolesCreated++
		} else if role.Action == "updated" {
			rolesUpdated++
		}
	}

	usersCreated := 0
	usersUpdated := 0
	for _, user := range resp.Users {
		if user.Action == "created" {
			usersCreated++
		} else if user.Action == "updated" {
			usersUpdated++
		}
	}

	// Exibir resumo
	fmt.Printf("Application: %s (%s)\n", resp.Application.Code, resp.Application.Action)
	fmt.Printf("Permissions: %d (%d criadas, %d atualizadas)\n", len(resp.Permissions), permsCreated, permsUpdated)
	fmt.Printf("Roles:       %d (%d criadas, %d atualizadas)\n", len(resp.Roles), rolesCreated, rolesUpdated)
	if len(resp.Users) > 0 {
		fmt.Printf("Users:       %d (%d criados, %d atualizados)\n", len(resp.Users), usersCreated, usersUpdated)
	}
	fmt.Println("\nSync concluído com sucesso.")

	return nil
}

// RunSyncWithExit executa RunSync e faz os.Exit apropriado em caso de erro
func RunSyncWithExit(manifestPath string, cfg *config.Config) {
	if err := RunSync(manifestPath, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
		os.Exit(1)
	}
}

