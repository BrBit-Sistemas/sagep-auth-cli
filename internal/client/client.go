package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/BrBit-Sistemas/sagep-auth-cli/internal/manifest"
)

// AuthClient é o cliente HTTP para comunicação com o sagep-auth
type AuthClient struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// NewAuthClient cria uma nova instância do AuthClient
func NewAuthClient(baseURL, token string) *AuthClient {
	return &AuthClient{
		BaseURL: baseURL,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SyncResultDTO representa o resultado de uma operação de sync
type SyncResultDTO struct {
	Code   string `json:"code"`
	Action string `json:"action"` // "created" ou "updated"
	ID     string `json:"id,omitempty"`
}

// SyncRoleResultDTO representa o resultado de sync de uma role
type SyncRoleResultDTO struct {
	Code        string         `json:"code"`
	Action      string         `json:"action"` // "created" ou "updated"
	ID          string         `json:"id,omitempty"`
	Permissions []SyncResultDTO `json:"permissions"`
}

// SyncResponse representa a resposta do endpoint /v1/applications/sync
type SyncResponse struct {
	Application SyncResultDTO       `json:"application"`
	Permissions []SyncResultDTO     `json:"permissions"`
	Roles       []SyncRoleResultDTO `json:"roles"`
}

// SyncApplication envia o manifest para o endpoint de sync do sagep-auth
func (c *AuthClient) SyncApplication(ctx context.Context, m *manifest.AuthManifest) (*SyncResponse, error) {
	// Converter manifest para JSON
	payload, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar manifest: %w", err)
	}

	// Construir requisição
	url := c.BaseURL + "/v1/applications/sync"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	// Headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	// Executar requisição
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar requisição: %w", err)
	}
	defer resp.Body.Close()

	// Ler resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	// Verificar status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("erro na API (status %d): %s", resp.StatusCode, string(body))
	}

	// Fazer unmarshal da resposta
	var syncResp SyncResponse
	if err := json.Unmarshal(body, &syncResp); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse da resposta: %w", err)
	}

	return &syncResp, nil
}

