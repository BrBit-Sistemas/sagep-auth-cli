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
# Primeiro, configure o ambiente Go (garante que repositório público não precise autenticação)
make setup-go-env

# Depois, compile o CLI
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

**Opção 4: Instalar diretamente do GitHub**

```bash
go install github.com/BrBit-Sistemas/sagep-auth-cli/cmd/sagep-auth-cli@latest
```

**Configuração do ambiente Go (GOPRIVATE):**

⚠️ **Importante:** Apenas `github.com/BrBit-Sistemas/sagep-auth-cli` é **público**. Todos os outros repositórios em `github.com/BrBit-Sistemas` são **privados**.

Para garantir que o Go possa acessar o `sagep-auth-cli` publicamente sem expor os outros repositórios privados, execute:

```bash
# Opção 1: Usar o script automático (recomendado)
make setup-go-env
```

O script `scripts/setup-go-env.sh` implementa a estratégia correta:

1. **Garante que `github.com/BrBit-Sistemas` está no `GOPRIVATE`**
   - Protege todos os outros repositórios privados da organização
   - Se não estiver, o script adiciona automaticamente

2. **Adiciona exceção via `GONOPROXY` e `GONOSUMDB`**
   - Adiciona `github.com/BrBit-Sistemas/sagep-auth-cli` às exceções
   - Permite acesso público apenas a este repositório específico
   - Mantém os outros repositórios protegidos pelo `GOPRIVATE`

**Opção 2: Configurar manualmente**

```bash
# 1. Garantir que GOPRIVATE contém github.com/BrBit-Sistemas (protege outros repositórios)
go env -w GOPRIVATE="github.com/BrBit-Sistemas"

# 2. Adicionar exceção para o repositório público
go env -w GONOPROXY="github.com/BrBit-Sistemas/sagep-auth-cli"
go env -w GONOSUMDB="github.com/BrBit-Sistemas/sagep-auth-cli"
```

**Importante:** Se você já tiver outros módulos privados no `GOPRIVATE`, adicione `github.com/BrBit-Sistemas` à lista existente:

```bash
# Exemplo: se você já tem outros módulos privados
go env -w GOPRIVATE="github.com/empresa1,github.com/BrBit-Sistemas,github.com/empresa2"
```

**Importante:** `GOPRIVATE`, `GONOPROXY` e `GONOSUMDB` são configurações do ambiente Go, não do CLI. Elas afetam apenas a instalação via `go install`, não o funcionamento do CLI em si.

## Configuração

### Arquivo .env (Recomendado)

O CLI suporta um arquivo `.env` na raiz do projeto para configurar as variáveis de ambiente. Isso é útil para desenvolvimento local e evita expor credenciais na linha de comando.

**Primeiro passo:** Copie o arquivo de exemplo:

```bash
cp .env.example .env
```

**Importante:** Edite o arquivo `.env` e configure:

- `SAGEP_AUTH_URL`: URL do serviço sagep-auth (obrigatório)
- `SAGEP_AUTH_SECRET`: Secret compartilhado para HMAC (bootstrap) **OU**
- `SAGEP_AUTH_TOKEN`: Token JWT (uso normal após bootstrap)

Você precisa configurar **pelo menos um** (`SAGEP_AUTH_SECRET` ou `SAGEP_AUTH_TOKEN`).

**Segundo passo:** Edite o arquivo `.env` e preencha com seus valores:

**Opção 1: Bootstrap inicial (sem aplicação/usuário ainda)**

```env
SAGEP_AUTH_URL=http://localhost:8080
SAGEP_AUTH_SECRET=seu-secret-compartilhado-aqui
```

**Opção 2: Uso normal (após bootstrap, com autenticação JWT)**

```env
SAGEP_AUTH_URL=https://auth.sagep.com.br
SAGEP_AUTH_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Importante:** O arquivo `.env` está no `.gitignore` e não será commitado no repositório.

### Variáveis de Ambiente

#### Obrigatória

- `SAGEP_AUTH_URL`: URL base do serviço sagep-auth (ex: `https://auth.sagep.com.br` ou `http://localhost:8080`)

#### Autenticação (escolha uma opção)

Você precisa configurar **pelo menos uma** das seguintes variáveis:

- **`SAGEP_AUTH_SECRET`** (Recomendado para bootstrap): Secret compartilhado com o servidor `sagep-auth` para autenticação HMAC. Deve ser o mesmo valor de `BOOTSTRAP_SECRET` no servidor. Permite criar aplicações iniciais sem autenticação JWT.
  - **Gere um secret seguro:** `openssl rand -base64 32`
  
- **`SAGEP_AUTH_TOKEN`** (Uso normal após bootstrap): Token JWT obtido via autenticação. Use quando a aplicação já foi criada e você tem um usuário/token.

**Como funciona:**
- Se `SAGEP_AUTH_SECRET` estiver configurado → CLI usa HMAC (bootstrap)
- Se `SAGEP_AUTH_TOKEN` estiver configurado → CLI usa JWT (uso normal)
- O servidor `sagep-auth` aceita ambos no endpoint `/applications/sync`

**Nota sobre `GOPRIVATE`:** Esta variável de ambiente do Go **não é necessária** para o funcionamento do CLI. Ela é uma configuração do ambiente Go que afeta como o Go baixa módulos. Se você tiver problemas ao instalar o CLI via `go install`, verifique se `GOPRIVATE` não está configurado para `github.com/BrBit-Sistemas` (veja seção de instalação acima).

### Ordem de Precedência

1. **Flags de linha de comando** (maior precedência) - `--url`, `--token` ou `--secret`
2. **Arquivo `.env`** - na raiz do projeto
3. **Variáveis de ambiente do sistema** - `SAGEP_AUTH_URL`, `SAGEP_AUTH_SECRET` ou `SAGEP_AUTH_TOKEN`

**Nota:** Se nenhuma das opções acima fornecer `SAGEP_AUTH_URL` e pelo menos uma das opções de autenticação (`SAGEP_AUTH_SECRET` ou `SAGEP_AUTH_TOKEN`), o CLI retornará um erro.

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

**Exemplos:**

⚠️ **Importante:** Os flags devem vir **antes** do comando `sync`:

```bash
# Opção 1: Usando arquivo .env (recomendado)
# Certifique-se de ter o arquivo .env configurado na raiz do projeto
sagep-auth-cli --manifest ./auth-manifest.yaml sync

# Opção 2: Usando variáveis de ambiente do sistema (bootstrap com HMAC)
export SAGEP_AUTH_URL=https://auth.sagep.com.br
export SAGEP_AUTH_SECRET=seu-secret-compartilhado
sagep-auth-cli --manifest ./auth-manifest.yaml sync

# Opção 2b: Usando variáveis de ambiente do sistema (uso normal com JWT)
export SAGEP_AUTH_URL=https://auth.sagep.com.br
export SAGEP_AUTH_TOKEN=seu-token-jwt-aqui
sagep-auth-cli --manifest ./auth-manifest.yaml sync

# Opção 3: Usando flags (override completo)
sagep-auth-cli \
  --manifest ./auth-manifest.yaml \
  --url https://auth.sagep.com.br \
  --token seu-token-jwt-aqui \
  sync

# Opção 4: Usando flag curta para manifest
sagep-auth-cli -m ./auth-manifest.yaml sync

# Opção 5: Usando manifest padrão (./auth-manifest.yaml)
# Se o arquivo se chamar exatamente auth-manifest.yaml na raiz
sagep-auth-cli sync
```

**Nota:** A ordem `sagep-auth-cli sync --manifest` também funciona, mas a ordem recomendada é `sagep-auth-cli --manifest sync` (flags antes do comando).

**Ordem de precedência:** Flags > `.env` > Variáveis de ambiente do sistema

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
# Exemplo para GitHub Actions (bootstrap com HMAC)
- name: Sync Auth Manifest
  env:
    SAGEP_AUTH_URL: ${{ secrets.SAGEP_AUTH_URL }}
    SAGEP_AUTH_SECRET: ${{ secrets.SAGEP_AUTH_SECRET }}
  run: |
    sagep-auth-cli sync --manifest ./auth-manifest.yaml

# Exemplo para GitHub Actions (uso normal com JWT)
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

