package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/BrBit-Sistemas/sagep-auth-cli/internal/manifest"
	"gopkg.in/yaml.v3"
)

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
	Description string
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
		fmt.Println("\n💡 Dica: Formato sugerido para códigos de permissão: {app}.{recurso}.{ação}")
		fmt.Println("   Exemplo: biopass.devices.read, biopass.devices.create")
		fmt.Println("   Ou para menus: Menu:Dashboard, Menu:Devices\n")

		for {
			var perm PermissionAnswer

			if err := survey.Ask([]*survey.Question{
				{
					Name: "code",
					Prompt: &survey.Input{
						Message: "Código da permissão:",
						Help:    "Ex: biopass.devices.read, Menu:Dashboard",
					},
					Validate: survey.Required,
				},
				{
					Name: "description",
					Prompt: &survey.Input{
						Message: "Descrição:",
					},
				},
			}, &perm); err != nil {
				break
			}

			perm.Code = strings.TrimSpace(perm.Code)
			perm.Description = strings.TrimSpace(perm.Description)
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

			// Selecionar permissões
			if len(answers.Permissions) > 0 {
				permOptions := make([]string, len(answers.Permissions))
				for i, p := range answers.Permissions {
					permOptions[i] = p.Code
				}

				var selectedPerms []string
				if err := survey.AskOne(&survey.MultiSelect{
					Message: "Selecione as permissões para esta role:",
					Options: permOptions,
					Help:    "Você pode usar wildcards no YAML manualmente depois (ex: biopass.*)",
				}, &selectedPerms); err == nil {
					role.Permissions = selectedPerms
				}
			} else {
				if err := survey.AskOne(&survey.Input{
					Message: "Permissões (separadas por vírgula ou wildcard como biopass.*):",
					Help:    "Ex: biopass.* ou biopass.devices.read,biopass.devices.create",
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

			role.Code = strings.TrimSpace(role.Code)
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
			Description: p.Description,
		}
	}

	for i, r := range answers.Roles {
		m.Roles[i] = manifest.Role{
			Code:        r.Code,
			Name:        r.Name,
			System:      r.System,
			Description: r.Description,
			Permissions: r.Permissions,
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

