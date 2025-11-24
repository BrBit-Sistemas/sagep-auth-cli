package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/BrBit-Sistemas/sagep-auth-cli/internal/manifest"
	"gopkg.in/yaml.v3"
)

// inferSubjectAndAction é um wrapper para a função do pacote manifest
func inferSubjectAndAction(code string) (string, string, bool) {
	return manifest.InferSubjectAndAction(code)
}

// getActionOptions retorna as opções de ações formatadas com descrições
// Baseado em melhores práticas: AWS IAM, Google Cloud, CASL.js
func getActionOptions() []string {
	return []string{
		"read - Consultar registros (listar ou visualizar individual)",
		"create - Criar novos registros",
		"update - Atualizar registros existentes",
		"delete - Remover registros",
		"manage - Controle total (todas as operações: read, create, update, delete)",
		"view - Visualizar interface/telas (usado principalmente para menus)",
	}
}

// extractActionValue extrai apenas o valor da ação (sem descrição)
// Ex: "read - Consultar..." → "read"
func extractActionValue(selectedOption string) string {
	parts := strings.SplitN(selectedOption, " - ", 2)
	if len(parts) > 0 {
		return strings.TrimSpace(parts[0])
	}
	return strings.TrimSpace(selectedOption)
}

// findActionOption encontra a opção formatada correspondente a uma ação
// Ex: "read" → "read - Consultar registros..."
func findActionOption(action string) string {
	options := getActionOptions()
	for _, opt := range options {
		if extractActionValue(opt) == action {
			return opt
		}
	}
	// Fallback: retorna a primeira opção se não encontrar
	if len(options) > 0 {
		return options[0]
	}
	return action
}

type InitAnswers struct {
	AppName        string
	AppCode        string
	AppDescription string
	
	CreateUsers    bool
	Users          []UserAnswer
	
	CreatePermissions bool
	Permissions      []PermissionAnswer
	
	CreateRoles    bool
	Roles          []RoleAnswer
}

type UserAnswer struct {
	Email    string
	Password string
	Name     string
	TenantID string // Opcional: unidade do usuário (deixe vazio para usuário global)
	IsMaster bool
	Roles    []string
}

type PermissionAnswer struct {
	Code        string
	Subject     string
	Action      string
	Description string
	Conditions  string
}

type RoleAnswer struct {
	Code        string
	Name        string
	System      bool
	Description string
	Permissions []string
}

func RunInit(manifestPath string) error {
	// Verificar se manifest já existe
	var existingManifest *manifest.AuthManifest
	manifestExists := false
	if _, err := os.Stat(manifestPath); err == nil {
		manifestExists = true
		loaded, err := manifest.LoadManifest(manifestPath)
		if err != nil {
			fmt.Printf("\n⚠️  Manifest existe mas não pôde ser carregado: %v\n", err)
			var proceed bool
			if err := survey.AskOne(&survey.Confirm{
				Message: "Deseja sobrescrever o manifest existente?",
				Default: false,
			}, &proceed); err != nil || !proceed {
				return fmt.Errorf("operação cancelada")
			}
		} else {
			existingManifest = loaded
		}
	}

	// Se manifest existe, perguntar o que fazer
	if manifestExists && existingManifest != nil {
		fmt.Printf("\n📄 Manifest encontrado: %s\n", manifestPath)
		fmt.Printf("   - Aplicação: %s (%s)\n", existingManifest.Application.Name, existingManifest.Application.Code)
		fmt.Printf("   - Permissões: %d\n", len(existingManifest.Permissions))
		fmt.Printf("   - Roles: %d\n", len(existingManifest.Roles))
		fmt.Printf("   - Usuários: %d\n\n", len(existingManifest.Users))
		
		var action string
		if err := survey.AskOne(&survey.Select{
			Message: "O que deseja fazer?",
			Options: []string{
				"adicionar - Adicionar novos recursos ao manifest existente",
				"sobrescrever - Criar um novo manifest (perde dados existentes)",
				"cancelar - Não fazer nada",
			},
			Default: "adicionar - Adicionar novos recursos ao manifest existente",
		}, &action); err != nil {
			return err
		}

		if strings.Contains(action, "cancelar") {
			return fmt.Errorf("operação cancelada")
		}

		if strings.Contains(action, "sobrescrever") {
			var confirm bool
			if err := survey.AskOne(&survey.Confirm{
				Message: "⚠️  ATENÇÃO: Isso vai apagar todos os recursos existentes. Continuar?",
				Default: false,
			}, &confirm); err != nil || !confirm {
				return fmt.Errorf("operação cancelada")
			}
			existingManifest = nil // Resetar para criar do zero
		}
	}

	var answers InitAnswers
	
	// Se tem manifest existente e vamos adicionar, carregar dados atuais
	if existingManifest != nil {
		answers.AppCode = existingManifest.Application.Code
		answers.AppName = existingManifest.Application.Name
		answers.AppDescription = existingManifest.Application.Description
		
		// Converter permissions existentes
		answers.Permissions = make([]PermissionAnswer, len(existingManifest.Permissions))
		for i, p := range existingManifest.Permissions {
			answers.Permissions[i] = PermissionAnswer{
				Code:        p.Code,
				Subject:     p.Subject,
				Action:      p.Action,
				Description: p.Description,
				Conditions:  p.Conditions,
			}
		}
		
		// Converter roles existentes
		answers.Roles = make([]RoleAnswer, len(existingManifest.Roles))
		for i, r := range existingManifest.Roles {
			answers.Roles[i] = RoleAnswer{
				Code:        r.Code,
				Name:        r.Name,
				System:      r.System,
				Description: r.Description,
				Permissions: r.Permissions,
			}
		}
		
		// Converter users existentes
		answers.Users = make([]UserAnswer, len(existingManifest.Users))
		for i, u := range existingManifest.Users {
			tenantID := ""
			if u.TenantID != nil {
				tenantID = *u.TenantID
			}
			answers.Users[i] = UserAnswer{
				Email:    u.Email,
				Password: u.Password,
				Name:     u.Name,
				TenantID: tenantID,
				Roles:    u.Roles,
			}
		}
		
		fmt.Println("\n📝 Modo: Adicionar recursos ao manifest existente")
		fmt.Println("═══════════════════════════════════════════════════════\n")
	} else {
		fmt.Println("\n🚀 Criando novo manifest para integração com sagep-auth")
		fmt.Println("═══════════════════════════════════════════════════════\n")
	}

	// 1. Informações da Aplicação (apenas se criando novo manifest)
	if existingManifest == nil {
		var appInput struct {
			AppName        string
			AppDescription string
		}
		if err := survey.Ask([]*survey.Question{
		{
			Name: "appName",
			Prompt: &survey.Input{
				Message: "Nome da aplicação (ex: Biopass, CRV, Core):",
				Help:    "Informe apenas o nome básico. O CLI gerará o código automaticamente.",
			},
			Validate: survey.Required,
		},
		{
			Name: "appDescription",
			Prompt: &survey.Input{
				Message: "Descrição (opcional):",
			},
		},
	}, &appInput); err != nil {
		return err
	}
	
	// Inferir code e name a partir do input
	appNameInput := strings.TrimSpace(appInput.AppName)
	answers.AppName = manifest.InferApplicationName(appNameInput)
	answers.AppCode = manifest.InferApplicationCode(appNameInput)
	answers.AppDescription = strings.TrimSpace(appInput.AppDescription)
	
	// Mostrar o que foi inferido
	fmt.Printf("\n   ✅ Informações inferidas:\n")
	fmt.Printf("      Código: %s\n", answers.AppCode)
	fmt.Printf("      Nome:   %s\n", answers.AppName)
	
	// Permitir editar se necessário
	var confirmApp bool
	if err := survey.AskOne(&survey.Confirm{
		Message: "Confirmar informações da aplicação?",
		Default: true,
	}, &confirmApp); err != nil {
		return err
	}
	
	if !confirmApp {
		// Permitir editar manualmente
		if err := survey.Ask([]*survey.Question{
			{
				Name: "appCode",
				Prompt: &survey.Input{
					Message: "Código da aplicação:",
					Default: answers.AppCode,
				},
				Validate: survey.Required,
			},
			{
				Name: "appName",
				Prompt: &survey.Input{
					Message: "Nome da aplicação:",
					Default: answers.AppName,
				},
				Validate: survey.Required,
			},
		}, &answers); err != nil {
			return err
		}
		answers.AppCode = strings.ToLower(strings.TrimSpace(answers.AppCode))
	} else {
		fmt.Printf("✅ Usando aplicação existente: %s (%s)\n\n", answers.AppName, answers.AppCode)
	}
		answers.AppName = strings.TrimSpace(answers.AppName)
	}

	// 2. Usuários
	if err := survey.AskOne(&survey.Confirm{
		Message: "Deseja criar usuários iniciais?",
		Default: true,
	}, &answers.CreateUsers); err != nil {
		return err
	}

	if answers.CreateUsers {
		for {
			var user UserAnswer
			
			isMaster := false
			if err := survey.AskOne(&survey.Confirm{
				Message: "Este é um usuário Master?",
				Default: len(answers.Users) == 0,
				Help:    "Usuário Master tem acesso total (bypass de permissões)",
			}, &isMaster); err != nil {
				break
			}
			user.IsMaster = isMaster

			userQuestions := []*survey.Question{
				{
					Name: "email",
					Prompt: &survey.Input{
						Message: "Email do usuário:",
					},
					Validate: survey.Required,
				},
				{
					Name: "password",
					Prompt: &survey.Password{
						Message: "Senha:",
					},
					Validate: survey.Required,
				},
				{
					Name: "name",
					Prompt: &survey.Input{
						Message: "Nome completo:",
					},
					Validate: survey.Required,
				},
				{
					Name: "tenantID",
					Prompt: &survey.Input{
						Message: "Unidade (tenant_id) do usuário:",
						Help:    "Opcional - deixe vazio para usuário global. Especialmente útil para primeiro usuário/bootstrap (ex: unidade-005)",
					},
				},
			}

			if err := survey.Ask(userQuestions, &user); err != nil {
				break
			}

			// Limpar tenant_id se vazio (para não incluir no YAML)
			user.TenantID = strings.TrimSpace(user.TenantID)

			if isMaster {
				user.Roles = []string{"master"}
			} else if len(answers.Roles) > 0 {
				roleCodes := make([]string, len(answers.Roles))
				for i, r := range answers.Roles {
					roleCodes[i] = r.Code
				}
				
				var selectedRoles []string
				if err := survey.AskOne(&survey.MultiSelect{
					Message: "Selecione as roles para este usuário:",
					Options: roleCodes,
				}, &selectedRoles); err == nil {
					user.Roles = selectedRoles
				}
			}

			answers.Users = append(answers.Users, user)

			var addMore bool
			if err := survey.AskOne(&survey.Confirm{
				Message: "Adicionar outro usuário?",
				Default: false,
			}, &addMore); err != nil || !addMore {
				break
			}
		}
	}

	// 3. Permissões
	if err := survey.AskOne(&survey.Confirm{
		Message: "Deseja criar permissões agora?",
		Default: true,
		Help:    "Você pode adicionar mais depois executando 'init' novamente",
	}, &answers.CreatePermissions); err != nil {
		return err
	}

	if answers.CreatePermissions {
		fmt.Println("\n💡 Você pode criar permissões de Menu ou de Recurso (entidade).\n")

		for {
			var perm PermissionAnswer
			
			// 1. Perguntar tipo de permissão
			var permType string
			if err := survey.AskOne(&survey.Select{
				Message: "Tipo de permissão:",
				Options: []string{"Menu", "Recurso (entidade)"},
				Help:    "Menu: controle de visibilidade | Recurso: operações em entidades",
			}, &permType); err != nil {
				break
			}
			
			if permType == "Menu" {
				// 2a. Permissão de Menu - UX simplificada
				var menuInput struct {
					MenuName string
				}
				if err := survey.Ask([]*survey.Question{
					{
						Name: "menuName",
						Prompt: &survey.Input{
							Message: "Nome do menu (ex: Dashboard, Participantes):",
							Help:    "Informe apenas o nome. O CLI criará automaticamente Menu:{Nome}",
						},
						Validate: survey.Required,
					},
				}, &menuInput); err != nil {
					break
				}
				
				// Inferir automaticamente
				code, subject, action := manifest.InferMenuPermission(menuInput.MenuName)
				perm.Code = code
				perm.Subject = subject
				perm.Action = action
				
				fmt.Printf("\n   ✅ Permissão de menu criada:\n")
				fmt.Printf("      Code:    %s\n", perm.Code)
				fmt.Printf("      Subject: %s\n", perm.Subject)
				fmt.Printf("      Action:  %s\n", perm.Action)
				
			} else {
				// 2b. Permissão de Recurso - UX simplificada
				var resourceInput struct {
					Entidade string
					Action   string
				}
				if err := survey.Ask([]*survey.Question{
					{
						Name: "entidade",
						Prompt: &survey.Input{
							Message: "Nome da entidade (ex: participantes, devices, users):",
							Help:    "Use minúsculo, plural (como o frontend verifica no CASL.js)",
						},
						Validate: survey.Required,
					},
					{
						Name: "action",
						Prompt: &survey.Select{
							Message: "Operação permitida:",
							Options: getActionOptions(),
							Help:    "Ação que será permitida nesta entidade. Baseado em melhores práticas de autorização (CASL.js, AWS IAM, Google Cloud)",
						},
					},
				}, &resourceInput); err != nil {
					break
				}
				
				// Extrair apenas o valor da ação (sem descrição)
				actionValue := extractActionValue(resourceInput.Action)
				
				// Inferir automaticamente usando appCode
				code, subject, actionOut := manifest.InferResourcePermission(resourceInput.Entidade, actionValue, answers.AppCode)
				perm.Code = code
				perm.Subject = subject
				perm.Action = actionOut
				
				fmt.Printf("\n   ✅ Permissão de recurso criada:\n")
				fmt.Printf("      Code:    %s\n", perm.Code)
				fmt.Printf("      Subject: %s\n", perm.Subject)
				fmt.Printf("      Action:  %s\n", perm.Action)
			}

			// 3. Garantir que subject e action estão preenchidos (fallback de segurança)
			if perm.Subject == "" || perm.Action == "" {
				subject, action, inferred := inferSubjectAndAction(perm.Code)
				if inferred {
					perm.Subject = subject
					perm.Action = action
				} else {
					fmt.Println("\n   ⚠️  Erro: Não foi possível inferir subject e action.")
					fmt.Println("   Por favor, tente novamente.")
					continue
				}
			}
			
			// 4. Permitir editar se necessário (opcional)
			var confirm bool
			if err := survey.AskOne(&survey.Confirm{
				Message: "Confirmar permissão criada?",
				Default: true,
			}, &confirm); err != nil {
				break
			}
			
			if !confirm {
				// Solicitar edição manual
				if err := survey.Ask([]*survey.Question{
					{
						Name: "code",
						Prompt: &survey.Input{
							Message: "Code:",
							Default: perm.Code,
						},
						Validate: survey.Required,
					},
					{
						Name: "subject",
						Prompt: &survey.Input{
							Message: "Subject:",
							Default: perm.Subject,
						},
						Validate: survey.Required,
					},
					{
						Name: "action",
						Prompt: &survey.Select{
							Message: "Action:",
							Options: getActionOptions(),
							Default: findActionOption(perm.Action),
							Help:    "Ação que será permitida nesta entidade. Baseado em melhores práticas de autorização (CASL.js, AWS IAM, Google Cloud)",
						},
					},
				}, &perm); err != nil {
					break
				}
			}

			// 5. Garantir que subject e action estão preenchidos
			perm.Subject = strings.TrimSpace(perm.Subject)
			// Extrair apenas o valor da ação se for uma opção formatada
			perm.Action = extractActionValue(strings.TrimSpace(perm.Action))
			
			if perm.Subject == "" || perm.Action == "" {
				fmt.Println("   ❌ Erro: Subject e Action são obrigatórios para compatibilidade com CASL.js")
				fmt.Println("   Por favor, tente novamente ou edite o manifest manualmente.")
				continue
			}

			// 6. Solicitar description e conditions
			if err := survey.Ask([]*survey.Question{
				{
					Name: "description",
					Prompt: &survey.Input{
						Message: "Descrição:",
					},
				},
				{
					Name: "conditions",
					Prompt: &survey.Input{
						Message: "Conditions (JSON opcional, ex: {\"userId\": \"${user.id}\"}):",
						Help:    "Deixe vazio se não precisar de condições",
					},
				},
			}, &perm); err != nil {
				break
			}

			perm.Description = strings.TrimSpace(perm.Description)
			perm.Conditions = strings.TrimSpace(perm.Conditions)
			answers.Permissions = append(answers.Permissions, perm)

			var addMore bool
			if err := survey.AskOne(&survey.Confirm{
				Message: "Adicionar outra permissão?",
				Default: true,
			}, &addMore); err != nil || !addMore {
				break
			}
		}
	}

	// 4. Roles
	if err := survey.AskOne(&survey.Confirm{
		Message: "Deseja criar roles base agora?",
		Default: true,
		Help:    "Roles base (system: true) são protegidas e não podem ser editadas via API",
	}, &answers.CreateRoles); err != nil {
		return err
	}

	if answers.CreateRoles {
		for {
			var role RoleAnswer

			if err := survey.Ask([]*survey.Question{
				{
					Name: "code",
					Prompt: &survey.Input{
						Message: "Código da role:",
						Help:    "Ex: biopass.admin, master",
					},
					Validate: survey.Required,
				},
				{
					Name: "name",
					Prompt: &survey.Input{
						Message: "Nome da role:",
						Help:    "Ex: Administrador BioPass",
					},
					Validate: survey.Required,
				},
				{
					Name: "description",
					Prompt: &survey.Input{
						Message: "Descrição:",
					},
				},
			}, &role); err != nil {
				break
			}

			if err := survey.AskOne(&survey.Confirm{
				Message: "Esta é uma role base do sistema (system: true)?",
				Default: true,
				Help:    "Roles base são protegidas e não podem ser editadas via API",
			}, &role.System); err != nil {
				break
			}

			// Trim do código antes de verificar se é master
			role.Code = strings.TrimSpace(role.Code)
			isMasterRole := strings.ToLower(role.Code) == "master"

		// Master não precisa de permissões - o sistema retorna {action: "manage", subject: "all"} automaticamente
		// IMPORTANTE: Master sempre deve ter permissions: [] para garantir que o sistema retorne o acesso total
		if isMasterRole {
			role.Permissions = []string{}
			fmt.Println("   ℹ️  Role 'master' não precisa de permissões")
			fmt.Println("      O sistema retorna automaticamente: {action: \"manage\", subject: \"all\"}")
			} else {
			// Selecionar permissões para roles não-master
			// IMPORTANTE: Wildcards funcionam (ex: biopass.*), mas cada permission no banco
			// precisa ter subject e action corretos para compatibilidade com CASL.js
			if len(answers.Permissions) > 0 {
				permOptions := make([]string, len(answers.Permissions))
				for i, p := range answers.Permissions {
					permOptions[i] = p.Code
				}

				var selectedPerms []string
				if err := survey.AskOne(&survey.MultiSelect{
					Message: "Selecione as permissões para esta role:",
					Options: permOptions,
					Help:    "Você pode usar wildcards no YAML manualmente depois (ex: biopass.*)\n" +
						"Nota: Wildcards funcionam, mas cada permission precisa ter subject/action corretos no banco",
				}, &selectedPerms); err == nil {
					role.Permissions = selectedPerms
				}
			} else {
				if err := survey.AskOne(&survey.Input{
					Message: "Permissões (separadas por vírgula ou wildcard como biopass.*):",
					Help:    "Ex: biopass.* ou biopass.devices.read,biopass.devices.create\n" +
						"Nota: Wildcards funcionam, mas cada permission precisa ter subject/action corretos no banco",
				}, &role.Permissions); err != nil {
					role.Permissions = []string{}
				} else {
					permsStr := strings.TrimSpace(role.Permissions[0])
					if permsStr != "" {
						role.Permissions = strings.Split(permsStr, ",")
						for i := range role.Permissions {
							role.Permissions[i] = strings.TrimSpace(role.Permissions[i])
						}
					}
				}
			}
			}

			// Code já foi trimado acima
			role.Name = strings.TrimSpace(role.Name)
			role.Description = strings.TrimSpace(role.Description)
			answers.Roles = append(answers.Roles, role)

			var addMore bool
			if err := survey.AskOne(&survey.Confirm{
				Message: "Adicionar outra role?",
				Default: true,
			}, &addMore); err != nil || !addMore {
				break
			}
		}
	}

	// Verificar se algum usuário tem role "master" e garantir que a role master existe
	hasMasterUser := false
	for _, user := range answers.Users {
		for _, role := range user.Roles {
			if strings.ToLower(role) == "master" {
				hasMasterUser = true
				break
			}
		}
		if hasMasterUser {
			break
		}
	}

	// Se há usuário master, garantir que a role master existe e tem permissions vazias
	if hasMasterUser {
		hasMasterRole := false
		for i := range answers.Roles {
			if strings.ToLower(answers.Roles[i].Code) == "master" {
				hasMasterRole = true
				// Garantir que permissions está vazio
				answers.Roles[i].Permissions = []string{}
				break
			}
		}
		
		// Se não existe role master, criar automaticamente
		if !hasMasterRole {
			fmt.Println("\n   ⚠️  Usuário Master detectado, mas role 'master' não foi criada.")
			fmt.Println("   ✅ Criando role 'master' automaticamente com permissions vazias...")
			answers.Roles = append(answers.Roles, RoleAnswer{
				Code:        "master",
				Name:        "Master",
				System:      true,
				Description: "Role Master - acesso total ao sistema",
				Permissions: []string{},
			})
		}
	}

	// Criar manifest a partir das respostas
	m := buildManifestFromAnswers(answers)

	// Salvar arquivo
	return saveManifest(m, manifestPath)
}

func buildManifestFromAnswers(answers InitAnswers) *manifest.AuthManifest {
	m := &manifest.AuthManifest{
		Application: manifest.Application{
			Code:        answers.AppCode,
			Name:        answers.AppName,
			Description: answers.AppDescription,
		},
		Permissions: make([]manifest.Permission, len(answers.Permissions)),
		Roles:       make([]manifest.Role, len(answers.Roles)),
		Users:       make([]manifest.User, len(answers.Users)),
	}

	for i, p := range answers.Permissions {
		m.Permissions[i] = manifest.Permission{
			Code:        p.Code,
			Subject:     p.Subject,
			Action:      p.Action,
			Description: p.Description,
			Conditions:  p.Conditions,
		}
	}

	for i, r := range answers.Roles {
		// Master não precisa de permissões - garantir array vazio sempre
		// IMPORTANTE: O sistema detecta role master e retorna {action: "manage", subject: "all"} automaticamente
		permissions := r.Permissions
		if strings.ToLower(r.Code) == "master" {
			permissions = []string{}
		}
		
		m.Roles[i] = manifest.Role{
			Code:        r.Code,
			Name:        r.Name,
			System:      r.System,
			Description: r.Description,
			Permissions: permissions,
		}
	}

	for i, u := range answers.Users {
		var tenantID *string
		if strings.TrimSpace(u.TenantID) != "" {
			tenantIDValue := strings.TrimSpace(u.TenantID)
			tenantID = &tenantIDValue
		}
		
		m.Users[i] = manifest.User{
			Email:    u.Email,
			Password: u.Password,
			Name:     u.Name,
			TenantID: tenantID,
			Active:   true,
			Roles:    u.Roles,
		}
	}

	return m
}

func saveManifest(m *manifest.AuthManifest, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo: %w", err)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	
	if err := encoder.Encode(m); err != nil {
		return fmt.Errorf("erro ao escrever YAML: %w", err)
	}

	fmt.Printf("\n✅ Manifest criado com sucesso: %s\n", path)
	fmt.Printf("\n📋 Resumo:\n")
	fmt.Printf("   - Aplicação: %s (%s)\n", m.Application.Name, m.Application.Code)
	fmt.Printf("   - Permissões: %d\n", len(m.Permissions))
	fmt.Printf("   - Roles: %d\n", len(m.Roles))
	fmt.Printf("   - Usuários: %d\n", len(m.Users))
	fmt.Printf("\n🚀 Próximo passo: Execute 'sync' para sincronizar com o servidor\n")

	return nil
}

