# sagep-auth-cli

CLI em Go para sincronizar manifests de autenticação de aplicações SAGEP com o serviço central `sagep-auth`.

## Propósito

O `sagep-auth-cli` permite que qualquer aplicação SAGEP (ex: `sagep-biopass`, `sagep-crv`, etc.) sincronize suas permissões e roles base com o serviço de autenticação central através de um arquivo manifest em YAML.

## Instalação

### Requisitos

- Go 1.21 ou superior

### Build Local

**Opção 1: Usando Makefile (recomendado)**

```bash
make build
```

**Opção 2: Build direto**

```bash
go build -o sagep-auth-cli ./cmd/sagep-auth-cli
```

**Opção 3: Instalar globalmente (desde o diretório do projeto)**

```bash
# A partir da raiz do projeto
go install ./cmd/sagep-auth-cli
```

O binário será instalado em `$GOPATH/bin/sagep-auth-cli`. Certifique-se de que `$GOPATH/bin` está no seu `PATH` para usar o comando diretamente, ou use o caminho completo.

**Nota:** O comando `go install github.com/brbit/sagep-auth-cli/cmd/sagep-auth-cli@latest` só funcionará após o módulo ser publicado em um repositório Git público. Para uso local, use uma das opções acima.

## Configuração

### Variáveis de Ambiente

O CLI suporta as seguintes variáveis de ambiente:

- `SAGEP_AUTH_URL`: URL base do serviço sagep-auth (ex: `https://auth.sagep.com.br`)
- `SAGEP_AUTH_TOKEN`: Token ou API key para autenticação na API do auth

### Ordem de Precedência

1. Flags de linha de comando (maior precedência)
2. Variáveis de ambiente
3. Valores default (quando aplicável)

## Uso

### Comando Sync

O comando principal é `sync`, que sincroniza um manifest YAML com o serviço `sagep-auth`:

```bash
sagep-auth-cli sync --manifest ./auth-manifest.yaml
```

**Opções:**

- `--manifest` ou `-m`: Caminho do arquivo manifest YAML (default: `./auth-manifest.yaml`)
- `--url`: Override da URL do auth (opcional)
- `--token`: Override do token de autenticação (opcional)
- `--help`: Exibir ajuda

**Exemplo:**

```bash
# Usando variáveis de ambiente
export SAGEP_AUTH_URL=https://auth.sagep.com.br
export SAGEP_AUTH_TOKEN=seu-token-aqui
sagep-auth-cli sync

# Usando flags
sagep-auth-cli sync --manifest ./auth-manifest.yaml --url https://auth.sagep.com.br --token seu-token
```

## Formato do Manifest

O arquivo `auth-manifest.yaml` deve seguir este formato:

```yaml
application:
  code: sagep-biopass
  name: SAGEP Biopass
  description: Sistema de controle de ponto biométrico

permissions:
  - code: biopass.devices.read
    description: Listar e visualizar dispositivos
  - code: biopass.devices.create
    description: Criar dispositivos
  - code: biopass.devices.update
    description: Atualizar dispositivos
  - code: biopass.devices.delete
    description: Excluir dispositivos

roles:
  - code: BIOPASS_ADMIN
    name: Administrador Biopass
    system: true
    description: Acesso total ao sistema Biopass
    permissions:
      - biopass.*
  - code: BIOPASS_DEVICES_READONLY
    name: Consulta de Dispositivos Biopass
    system: true
    description: Apenas leitura de dispositivos
    permissions:
      - biopass.devices.read
```

### Campos Obrigatórios

- `application.code`: Código único da aplicação
- `application.name`: Nome da aplicação
- `permissions[].code`: Código da permissão
- `roles[].code`: Código da role
- `roles[].name`: Nome da role
- `roles[].system`: Boolean indicando se é role base (true) ou customizada (false)
- `roles[].permissions`: Lista de códigos de permissões ou wildcards (ex: `biopass.*`)

### Wildcards

O manifest suporta wildcards em permissões. Qualquer prefixo seguido de `.*` expande para todas as permissões que começam com esse prefixo:

- `biopass.*` → todas as permissões que começam com `biopass.`
- `crv.orders.*` → todas as permissões que começam com `crv.orders.`

## Integração com CI/CD

O CLI foi projetado para ser usado em pipelines de CI/CD. Após o deploy (ou durante o build), execute o sync para garantir que o serviço de autenticação esteja atualizado com as permissões e roles base da aplicação:

```yaml
# Exemplo para GitHub Actions
- name: Sync Auth Manifest
  env:
    SAGEP_AUTH_URL: ${{ secrets.SAGEP_AUTH_URL }}
    SAGEP_AUTH_TOKEN: ${{ secrets.SAGEP_AUTH_TOKEN }}
  run: |
    sagep-auth-cli sync --manifest ./auth-manifest.yaml
```

## Saída do Comando

O comando exibe um resumo da sincronização:

```
Sincronizando aplicação: sagep-biopass
URL do auth: https://auth.sagep.com.br

Application: sagep-biopass (updated)
Permissions: 4 (2 criadas, 2 atualizadas)
Roles:       2 (1 criada, 1 atualizada)

Sync concluído com sucesso.
```

## Tratamento de Erros

Em caso de erro, o CLI:

1. Exibe uma mensagem clara indicando o problema
2. Retorna código de saída diferente de zero
3. Fornece informações úteis para debug

**Exemplos de erros comuns:**

- Manifest não encontrado
- YAML inválido
- Campos obrigatórios ausentes
- Erro de conexão com a API
- Token inválido ou expirado
- Status HTTP não 2xx

## Estrutura do Projeto

```
sagep-auth-cli/
├── cmd/
│   └── sagep-auth-cli/
│       └── main.go          # Ponto de entrada do CLI
├── internal/
│   ├── config/               # Leitura de configurações
│   │   └── config.go
│   ├── manifest/             # Structs e validação do YAML
│   │   └── manifest.go
│   ├── client/               # HTTP client para sagep-auth
│   │   └── client.go
│   └── commands/              # Implementação dos comandos
│       └── sync.go
├── go.mod
└── README.md
```

## Desenvolvimento

### Executar Localmente

```bash
# Desenvolvimento
go run cmd/sagep-auth-cli/main.go sync --manifest ./auth-manifest.yaml

# Build
go build -o sagep-auth-cli cmd/sagep-auth-cli/main.go
```

### Adicionar Novos Comandos

O projeto está organizado para facilitar a adição de novos comandos. Basta:

1. Criar uma nova função em `internal/commands/`
2. Adicionar o case no `switch` do `main.go`
3. Seguir o mesmo padrão de tratamento de erros e saída

## Licença

[Adicionar licença conforme necessário]

