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
	fmt.Println("\n🚀 Criando novo manifest para integração com sagep-auth")
	fmt.Println("═══════════════════════════════════════════════════════\n")

	var answers InitAnswers

	// 1. Informações da Aplicação
	if err := survey.Ask([]*survey.Question{
		{
			Name: "appCode",
			Prompt: &survey.Input{
				Message: "Código da aplicação (slug, ex: sagep-biopass):",
				Help:    "Será usado como identificador único. Ex: sagep-biopass, sagep-crv",
			},
			Validate: survey.Required,
		},
		{
			Name: "appName",
			Prompt: &survey.Input{
				Message: "Nome da aplicação:",
				Help:    "Nome amigável exibido no sistema. Ex: SAGEP Biopass",
			},
			Validate: survey.Required,
		},
		{
			Name: "appDescription",
			Prompt: &survey.Input{
				Message: "Descrição (opcional):",
			},
		},
	}, &answers); err != nil {
		return err
	}
	answers.AppCode = strings.ToLower(strings.TrimSpace(answers.AppCode))
	answers.AppName = strings.TrimSpace(answers.AppName)
	answers.AppDescription = strings.TrimSpace(answers.AppDescription)

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
			}

			if err := survey.Ask(userQuestions, &user); err != nil {
				break
			}

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
		fmt.Println("\n💡 Dica: Formato sugerido para códigos de permissão:")
		fmt.Println("   - {app}.{recurso}.{ação}: biopass.devices.read")
		fmt.Println("   - {recurso}.{ação}: Device.read")
		fmt.Println("   - Menu:{Nome}: Menu:Dashboard")
		fmt.Println("   O sistema tentará inferir Subject e Action automaticamente.\n")

		for {
			var perm PermissionAnswer

			// 1. Solicitar code
			if err := survey.Ask([]*survey.Question{
				{
					Name: "code",
					Prompt: &survey.Input{
						Message: "Código da permissão:",
						Help:    "Ex: biopass.devices.read, Device.read, Menu:Dashboard",
					},
					Validate: survey.Required,
				},
			}, &perm); err != nil {
				break
			}

			perm.Code = strings.TrimSpace(perm.Code)

			// 2. Tentar inferir subject e action
			subject, action, inferred := inferSubjectAndAction(perm.Code)
			
			if inferred {
				perm.Subject = subject
				perm.Action = action
				fmt.Printf("\n   ✅ Inferência automática:\n")
				fmt.Printf("      Subject: %s\n", subject)
				fmt.Printf("      Action:  %s\n", action)
				
				// 3. Permitir editar se necessário
				var confirm bool
				if err := survey.AskOne(&survey.Confirm{
					Message: "Confirmar subject e action inferidos?",
					Default: true,
				}, &confirm); err != nil {
					break
				}

				if !confirm {
					// Solicitar edição
					if err := survey.Ask([]*survey.Question{
						{
							Name: "subject",
							Prompt: &survey.Input{
								Message: "Subject (recurso, ex: Device, User, Menu:Dashboard):",
								Default: subject,
							},
							Validate: survey.Required,
						},
						{
							Name: "action",
							Prompt: &survey.Select{
								Message: "Action (ação CASL.js):",
								Options: []string{"read", "create", "update", "delete", "manage", "view"},
								Default: action,
							},
						},
					}, &perm); err != nil {
						break
					}
				}
			} else {
				// 4. Se não conseguiu inferir, solicitar explicitamente
				fmt.Println("\n   ⚠️  Não foi possível inferir automaticamente.")
				if err := survey.Ask([]*survey.Question{
					{
						Name: "subject",
						Prompt: &survey.Input{
							Message: "Subject (recurso, ex: Device, User, Menu:Dashboard):",
							Help:    "Nome do recurso para CASL.js",
						},
						Validate: survey.Required,
					},
					{
						Name: "action",
						Prompt: &survey.Select{
							Message: "Action (ação CASL.js):",
							Options: []string{"read", "create", "update", "delete", "manage", "view"},
							Help:    "Ação que será permitida",
						},
					},
				}, &perm); err != nil {
					break
				}
			}

			// 5. Garantir que subject e action estão preenchidos
			perm.Subject = strings.TrimSpace(perm.Subject)
			perm.Action = strings.TrimSpace(perm.Action)
			
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
		m.Users[i] = manifest.User{
			Email:    u.Email,
			Password: u.Password,
			Name:     u.Name,
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

