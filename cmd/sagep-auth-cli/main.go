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
		authToken = flag.String("token", "", "Token JWT de autenticação (override, uso normal)")
		authSecret = flag.String("secret", "", "Secret compartilhado para HMAC (override, bootstrap)")
		help = flag.Bool("help", false, "Exibir ajuda")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Uso: %s [opções] <comando>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Comandos:\n")
		fmt.Fprintf(os.Stderr, "  init    Cria um novo manifest interativamente\n")
		fmt.Fprintf(os.Stderr, "  sync    Sincroniza o manifest com o serviço sagep-auth\n\n")
		fmt.Fprintf(os.Stderr, "Opções:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nVariáveis de ambiente:\n")
		fmt.Fprintf(os.Stderr, "  SAGEP_AUTH_URL     URL base do serviço sagep-auth (obrigatório)\n")
		fmt.Fprintf(os.Stderr, "  SAGEP_AUTH_SECRET  Secret compartilhado para HMAC (bootstrap, opcional se tiver TOKEN)\n")
		fmt.Fprintf(os.Stderr, "  SAGEP_AUTH_TOKEN   Token JWT (uso normal, opcional se tiver SECRET)\n")
		fmt.Fprintf(os.Stderr, "\n  Você precisa configurar pelo menos um: SAGEP_AUTH_SECRET ou SAGEP_AUTH_TOKEN\n\n")
		fmt.Fprintf(os.Stderr, "Exemplos:\n")
		fmt.Fprintf(os.Stderr, "  %s init  # Cria manifest interativamente\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --manifest ./auth-manifest.yaml sync\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -m ./auth-manifest.yaml sync\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s sync  # usa ./auth-manifest.yaml (padrão)\n", os.Args[0])
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

	// Detectar se flags foram passados após o comando (ordem incorreta)
	if len(args) > 1 {
		nextArg := args[1]
		if nextArg == "--manifest" || nextArg == "-m" || nextArg == "--url" || nextArg == "--token" || nextArg == "--secret" {
			fmt.Fprintf(os.Stderr, "❌ Erro: Os flags devem vir ANTES do comando!\n\n")
			fmt.Fprintf(os.Stderr, "❌ Forma incorreta: %s %s %s ...\n", os.Args[0], args[0], nextArg)
			fmt.Fprintf(os.Stderr, "✅ Forma correta:   %s %s %s ...\n\n", os.Args[0], nextArg, args[0])
			fmt.Fprintf(os.Stderr, "Exemplos:\n")
			fmt.Fprintf(os.Stderr, "  %s --manifest ./auth-manifest.yaml sync\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "  %s -m ./auth-manifest.yaml sync\n", os.Args[0])
			os.Exit(1)
		}
	}

	command := args[0]

	// Determinar qual manifest usar: -m tem precedência sobre --manifest
	// Se nenhum for fornecido, usa o default
	manifest := defaultManifestPath
	
	// Verificar se --manifest foi usado (diferente do default)
	if *manifestPath != defaultManifestPath {
		manifest = *manifestPath
	}
	
	// Verificar se -m foi usado (tem precedência sobre --manifest)
	if *manifestPathShort != defaultManifestPath {
		manifest = *manifestPathShort
	}

	switch command {
	case "init":
		// Determinar caminho do manifest para criar
		initManifestPath := defaultManifestPath
		if *manifestPath != defaultManifestPath {
			initManifestPath = *manifestPath
		}
		if *manifestPathShort != defaultManifestPath {
			initManifestPath = *manifestPathShort
		}

		if err := commands.RunInit(initManifestPath); err != nil {
			fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
			os.Exit(1)
		}

	case "sync":
		// Carregar configuração
		cfg, err := config.LoadConfig(*authURL, *authToken, *authSecret)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro de configuração: %v\n", err)
			os.Exit(1)
		}

		// Executar sync
		commands.RunSyncWithExit(manifest, cfg)

	default:
		fmt.Fprintf(os.Stderr, "Erro: comando desconhecido '%s'\n\n", command)
		fmt.Fprintf(os.Stderr, "Comandos disponíveis: init, sync\n")
		os.Exit(1)
	}
}

