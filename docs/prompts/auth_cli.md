Quero que você atue como um arquiteto de plataforma e implemente, do zero, um CLI chamado `sagep-auth-cli`, em Go, para integrar qualquer aplicação SAGEP com o serviço de autenticação central `sagep-auth`, usando manifests em YAML.

Objetivo geral:
- O CLI deve ler um arquivo de manifest (`auth-manifest.yaml`) de uma aplicação qualquer (ex: sagep-biopass, sagep-crv, etc.).
- Converter esse manifest para um payload JSON.
- Enviar esse payload para o endpoint `/v1/applications/sync` do serviço `sagep-auth`.
- Exibir um resumo do que foi sincronizado.
- Ser reutilizável para qualquer app, sem lógica específica por aplicação.

### 1. Linguagem, estrutura e padrão do CLI

Use Go e organize o projeto mais ou menos assim:

- `cmd/sagep-auth-cli/main.go` → ponto de entrada do CLI.
- `internal/config` → leitura de configs (URL do auth, token, etc.).
- `internal/manifest` → structs e leitura/validação do YAML.
- `internal/client` → HTTP client para chamar o `sagep-auth`.
- `internal/commands` → implementação do comando `sync`.

Não precisa usar cobra obrigatoriamente, mas pode usar se quiser. Se preferir algo mais simples, use `flag` ou outra lib leve, contanto que o código fique organizado e extensível.

### 2. Comando principal do CLI

Implemente o comando:

```bash
sagep-auth-cli sync --manifest ./auth-manifest.yaml

Requisitos:

--manifest (ou -m) é o caminho do arquivo YAML.

Se não for informado, default para ./auth-manifest.yaml.

Flags extras opcionais:

--url (override da URL do auth, se necessário).

Se flags não forem passadas, use variáveis de ambiente:

SAGEP_AUTH_URL → URL base do sagep-auth (ex: https://auth.sagep.com.br
).

SAGEP_AUTH_TOKEN → token ou API key para autenticação na API do auth.

Ordem de precedência:

Flags > env vars > valores default (se fizer sentido).

3. Formato esperado do manifest

O CLI deve ler um YAML com esse formato (modelo):

application:
  code: sagep-biopass
  name: SAGEP Biopass

permissions:
  - code: biopass.devices.read
    description: Listar e visualizar dispositivos
  - code: biopass.devices.create
    description: Criar dispositivos

roles:
  - code: BIOPASS_ADMIN
    name: Administrador Biopass
    system: true
    permissions:
      - biopass.*
  - code: BIOPASS_DEVICES_READONLY
    name: Consulta de Dispositivos Biopass
    system: true
    permissions:
      - biopass.devices.read

Crie structs em Go para representar isso:

Application (code, name, etc).

Permission (code, description).

Role (code, name, system, permissions []string).

AuthManifest (application, permissions, roles).

4. Leitura e validação do manifest

No pacote internal/manifest:

Função LoadManifest(path string) (*AuthManifest, error) que:

Lê o arquivo YAML.

Faz unmarshal para os structs.

Valida:

application.code não pode ser vazio.

permissions não pode ter code vazio.

roles não pode ter code vazio.

roles devem referenciar permissions via permissions: []string (podendo conter wildcards ex: biopass.*).

Em caso de erro, retornar uma mensagem amigável indicando o problema.

5. Client HTTP para o sagep-auth

No pacote internal/client:

Criar uma struct AuthClient com:

BaseURL string

Token string (para header de auth)

Um http.Client interno.

Criar um método:
func (c *AuthClient) SyncApplication(ctx context.Context, manifest *AuthManifest) (*SyncResponse, error)
Que:

Constrói o payload em JSON esperado pelo sagep-auth para /v1/applications/sync.

Basicamente, é o próprio conteúdo do manifest, convertido para JSON:

application

permissions

roles

Envia um POST para: c.BaseURL + "/v1/applications/sync".

Inclui header de auth, por exemplo:

Authorization: Bearer <token>

Ou outro header se o padrão da sua API for diferente (deixe configurável).

Lida com erros de rede, status code não 2xx, etc.

Faz unmarshal da resposta em uma struct SyncResponse (você pode sugerir um modelo de resposta, ex: contagem de criados/atualizados).

Estruture o SyncResponse de forma genérica, por exemplo:

{
  "application": { "code": "sagep-biopass", "id": "..." },
  "stats": {
    "permissions_created": 4,
    "permissions_updated": 0,
    "roles_created": 3,
    "roles_updated": 1
  }
}

E crie a struct correspondente em Go.

6. Implementação do comando sync

No pacote internal/commands:

Criar uma função RunSync(manifestPath string, cfg Config) error que:

Carrega o manifest com LoadManifest.

Cria um AuthClient com base na config (URL, token).

Chama SyncApplication.

Exibe um resumo legível no stdout, por exemplo:

Sincronizando aplicação: sagep-biopass
URL do auth: https://auth.sagep.com.br

Permissions: 4 (2 criadas, 2 atualizadas)
Roles:       3 (3 criadas, 0 atualizadas)

Sync concluído com sucesso.

Em caso de erro, exibe mensagem clara e retorna código de saída diferente de zero.

7. Arquivo main.go

No cmd/sagep-auth-cli/main.go:

Ler as flags (--manifest, --url, etc).

Ler as variáveis de ambiente (SAGEP_AUTH_URL, SAGEP_AUTH_TOKEN).

Construir uma struct Config com os valores finais.

Chamar commands.RunSync(...).

Tratar erros e os.Exit(1) em caso de falha.

8. README do sagep-auth-cli

Crie um README.md explicando:

Propósito do CLI.

Instalação (como rodar localmente).

Variáveis de ambiente suportadas:

SAGEP_AUTH_URL

SAGEP_AUTH_TOKEN

Uso básico:

# dentro do repositório de qualquer app (que tenha auth-manifest.yaml)
sagep-auth-cli sync --manifest ./auth-manifest.yaml
Exemplo de manifest.

Explicação rápida de como isso será usado na pipeline de CI/CD:

Depois do deploy (ou no build), rodar o sagep-auth-cli sync para garantir que o auth esteja atualizado com as permissions e roles base da aplicação.

9. Qualidade de código

Comentar as principais funções.

Tratar erros de forma clara (mensagens amigáveis).

Organizar o código para ser fácil de evoluir (por exemplo, no futuro adicionar outros comandos além de sync).

Por favor, implemente tudo isso passo a passo, criando os arquivos e pastas necessários, e no final me mostre um resumo da estrutura de diretórios e exemplos de uso do CLI.
