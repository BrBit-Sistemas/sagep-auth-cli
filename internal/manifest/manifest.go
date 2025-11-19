package manifest

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Application representa a aplicação no manifest
type Application struct {
	Code        string `yaml:"code" json:"code"`
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
}

// Permission representa uma permissão no manifest
type Permission struct {
	Code        string `yaml:"code" json:"code"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
}

// Role representa uma role no manifest
type Role struct {
	Code        string   `yaml:"code" json:"code"`
	Name        string   `yaml:"name" json:"name"`
	System      bool     `yaml:"system" json:"system"`
	Description string   `yaml:"description,omitempty" json:"description,omitempty"`
	Permissions []string `yaml:"permissions" json:"permissions"`
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
	}

	// Validar roles
	for i, role := range m.Roles {
		if role.Code == "" {
			return fmt.Errorf("roles[%d].code não pode ser vazio", i)
		}
		if role.Name == "" {
			return fmt.Errorf("roles[%d].name não pode ser vazio", i)
		}
		if len(role.Permissions) == 0 {
			return fmt.Errorf("roles[%d].permissions não pode estar vazio", i)
		}
	}

	return nil
}

