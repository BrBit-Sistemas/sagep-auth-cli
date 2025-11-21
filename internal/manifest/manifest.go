package manifest

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Application representa a aplicação no manifest
type Application struct {
	Code        string `yaml:"code" json:"code"`
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
}

// Permission representa uma permissão no manifest
// IMPORTANTE: Subject e Action são obrigatórios para compatibilidade com CASL.js
// Wildcards funcionam nas roles (ex: biopass.*), mas cada permission no banco
// precisa ter subject e action corretos para o CASL.js funcionar corretamente
type Permission struct {
	Code        string `yaml:"code" json:"code"`                                 // Identificador único (ex: "biopass.devices.read")
	Subject     string `yaml:"subject,omitempty" json:"subject,omitempty"`       // OBRIGATÓRIO: Recurso para CASL.js (ex: "Device", "User", "Menu:Dashboard")
	Action      string `yaml:"action,omitempty" json:"action,omitempty"`        // OBRIGATÓRIO: Ação para CASL.js (ex: "read", "create", "update", "delete", "manage", "view")
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	Conditions  string `yaml:"conditions,omitempty" json:"conditions,omitempty"` // JSON opcional com condições (ex: {"userId": "${user.id}"})
}

// Role representa uma role no manifest
type Role struct {
	Code        string   `yaml:"code" json:"code"`
	Name        string   `yaml:"name" json:"name"`
	System      bool     `yaml:"system" json:"system"`
	Description string   `yaml:"description,omitempty" json:"description,omitempty"`
	Permissions []string `yaml:"permissions" json:"permissions"` // Lista de codes de permissions ou wildcards (ex: ["biopass.*"])
	// IMPORTANTE: Role "master" deve ter permissions: [] (vazio)
	// O sistema detecta role master e retorna automaticamente {action: "manage", subject: "all"} para CASL.js
}

// User representa um usuário no manifest
type User struct {
	Email    string   `yaml:"email" json:"email"`
	Password string   `yaml:"password" json:"password"`
	Name     string   `yaml:"name" json:"name"`
	Active   bool     `yaml:"active,omitempty" json:"active,omitempty"`
	Roles    []string `yaml:"roles" json:"roles"`
}

// AuthManifest representa o manifest completo
type AuthManifest struct {
	Application Application  `yaml:"application" json:"application"`
	Permissions []Permission  `yaml:"permissions" json:"permissions"`
	Roles       []Role        `yaml:"roles" json:"roles"`
	Users       []User        `yaml:"users,omitempty" json:"users,omitempty"`
}

// LoadManifest lê e valida um arquivo de manifest YAML
func LoadManifest(path string) (*AuthManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo manifest: %w", err)
	}

	var manifest AuthManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse do YAML: %w", err)
	}

	// Validações
	if err := validateManifest(&manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

// validateManifest valida o conteúdo do manifest
func validateManifest(m *AuthManifest) error {
	// Validar application
	if m.Application.Code == "" {
		return fmt.Errorf("application.code não pode ser vazio")
	}
	if m.Application.Name == "" {
		return fmt.Errorf("application.name não pode ser vazio")
	}

	// Validar permissions
	for i, perm := range m.Permissions {
		if perm.Code == "" {
			return fmt.Errorf("permissions[%d].code não pode ser vazio", i)
		}
		// Subject e Action são obrigatórios para compatibilidade com CASL.js
		// (exceto se o sistema conseguir inferir do code no backend)
		if perm.Subject == "" {
			return fmt.Errorf("permissions[%d].subject não pode ser vazio (necessário para CASL.js)", i)
		}
		if perm.Action == "" {
			return fmt.Errorf("permissions[%d].action não pode ser vazio (necessário para CASL.js)", i)
		}
		// Validar que action é uma ação válida do CASL.js
		validActions := []string{"read", "create", "update", "delete", "manage", "view"}
		actionValid := false
		for _, validAction := range validActions {
			if perm.Action == validAction {
				actionValid = true
				break
			}
		}
		if !actionValid {
			return fmt.Errorf("permissions[%d].action deve ser uma das ações válidas do CASL.js: read, create, update, delete, manage, view (atual: %s)", i, perm.Action)
		}
	}

	// Validar roles
	for i, role := range m.Roles {
		if role.Code == "" {
			return fmt.Errorf("roles[%d].code não pode ser vazio", i)
		}
		if role.Name == "" {
			return fmt.Errorf("roles[%d].name não pode ser vazio", i)
		}
		
		// Master sempre deve ter permissions vazio
		if strings.ToLower(role.Code) == "master" {
			if len(role.Permissions) > 0 {
				return fmt.Errorf("roles[%d] (master) deve ter permissions vazio - o sistema concede acesso total automaticamente", i)
			}
		} else {
			// Outras roles devem ter pelo menos uma permission
			if len(role.Permissions) == 0 {
				return fmt.Errorf("roles[%d].permissions não pode estar vazio", i)
			}
		}
	}

	return nil
}

