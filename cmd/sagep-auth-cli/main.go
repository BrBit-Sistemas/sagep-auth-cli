package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/BrBit-Sistemas/sagep-auth-cli/internal/commands"
	"github.com/BrBit-Sistemas/sagep-auth-cli/internal/config"
)

const (
	defaultManifestPath = "./auth-manifest.yaml"
)

func main() {
	// Definir flags
	var (
		manifestPath = flag.String("manifest", defaultManifestPath, "Caminho do arquivo manifest YAML")
		manifestPathShort = flag.String("m", defaultManifestPath, "Caminho do arquivo manifest YAML (short)")
		authURL = flag.String("url", "", "URL base do serviço sagep-auth (override)")
		authToken = flag.String("token", "", "Token de autenticação (override)")
		help = flag.Bool("help", false, "Exibir ajuda")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Uso: %s sync [opções]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Comandos:\n")
		fmt.Fprintf(os.Stderr, "  sync    Sincroniza o manifest com o serviço sagep-auth\n\n")
		fmt.Fprintf(os.Stderr, "Opções:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nVariáveis de ambiente:\n")
		fmt.Fprintf(os.Stderr, "  SAGEP_AUTH_URL     URL base do serviço sagep-auth\n")
		fmt.Fprintf(os.Stderr, "  SAGEP_AUTH_TOKEN   Token de autenticação\n\n")
		fmt.Fprintf(os.Stderr, "Exemplo:\n")
		fmt.Fprintf(os.Stderr, "  %s sync --manifest ./auth-manifest.yaml\n", os.Args[0])
	}

	flag.Parse()

	// Verificar se foi solicitada ajuda
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// Verificar comando
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Erro: comando não especificado\n\n")
		flag.Usage()
		os.Exit(1)
	}

	command := args[0]

	// Usar -m se fornecido, senão usar --manifest
	manifest := *manifestPath
	if *manifestPathShort != defaultManifestPath {
		manifest = *manifestPathShort
	}

	switch command {
	case "sync":
		// Carregar configuração
		cfg, err := config.LoadConfig(*authURL, *authToken)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro de configuração: %v\n", err)
			os.Exit(1)
		}

		// Executar sync
		commands.RunSyncWithExit(manifest, cfg)

	default:
		fmt.Fprintf(os.Stderr, "Erro: comando desconhecido '%s'\n\n", command)
		fmt.Fprintf(os.Stderr, "Comandos disponíveis: sync\n")
		os.Exit(1)
	}
}

