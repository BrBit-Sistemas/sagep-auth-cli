package client

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
	Secret     string // Secret para HMAC (bootstrap)
	HTTPClient *http.Client
}

// NewAuthClient cria uma nova instância do AuthClient
func NewAuthClient(baseURL, token, secret string) *AuthClient {
	return &AuthClient{
		BaseURL: baseURL,
		Token:   token,
		Secret:  secret,
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

// calculateHMAC calcula a assinatura HMAC do body + timestamp
func calculateHMAC(body []byte, timestamp int64, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	mac.Write([]byte(fmt.Sprintf("%d", timestamp)))
	return hex.EncodeToString(mac.Sum(nil))
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
	
	// Autenticação: HMAC (bootstrap) OU JWT (uso normal)
	if c.Secret != "" {
		// Usar HMAC para bootstrap
		timestamp := time.Now().Unix()
		signature := calculateHMAC(payload, timestamp, c.Secret)
		req.Header.Set("X-Signature", signature)
		req.Header.Set("X-Timestamp", fmt.Sprintf("%d", timestamp))
	} else if c.Token != "" {
		// Usar JWT para uso normal
	req.Header.Set("Authorization", "Bearer "+c.Token)
	} else {
		return nil, fmt.Errorf("SAGEP_AUTH_SECRET ou SAGEP_AUTH_TOKEN é obrigatório")
	}

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

