# 🚀 Guia Completo: Bootstrap do Sistema de Autenticação

Este guia mostra o passo a passo completo desde a configuração inicial do `sagep-auth` até a sincronização do manifest de uma aplicação (ex: `sagep-biopass-admin`).

---

## 📋 Pré-requisitos

- Go 1.21+ instalado
- PostgreSQL rodando
- Acesso aos repositórios:
  - `sagep-auth` (servidor de autenticação)
  - `sagep-auth-cli` (CLI para sincronização)
  - `sagep-biopass-admin` (aplicação exemplo)

---

## 🔧 PASSO 1: Configurar o Servidor `sagep-auth`

### 1.1. Clonar e entrar no diretório

```bash
cd ~/source/BrBit/sagep-auth
```

### 1.2. Criar arquivo `.env` a partir do exemplo

```bash
cp env.example .env
```

### 1.3. Gerar o Secret HMAC (BOOTSTRAP_SECRET)

```bash
# Gere um secret seguro (32 bytes em base64)
openssl rand -base64 32
```

**Exemplo de saída:**
```
Kx9mP2vQ8nR5tY7wZ3aB6cD9eF1gH4iJ7kL0mN2oP5qR8sT1uV4wX7yZ0aB3cD6eF9gH2iJ5kL8mN1oP4qR7sT0uV3wX6yZ9aB2cD5eF8gH1iJ4kL7mN0oP3qR6sT9uV2wX5yZ8aB1cD4eF7gH0iJ3kL6mN9oP2qR5sT8uV1wX4yZ7aB0cD3eF6gH9iJ2kL5mN8oP1qR4sT7uV0wX3yZ6aB9cD2eF5gH8iJ1kL4mN7oP0qR3sT6uV9wX2yZ5aB8cD1eF4gH7iJ0kL3mN6oP9qR2sT5uV8wX1yZ4aB7cD0eF3gH6iJ9kL2mN5oP8qR1sT4uV7wX0yZ3aB6cD9eF2gH5iJ8kL1mN4oP7qR0sT3uV6wX9yZ2aB5cD8eF1gH4iJ7kL0mN3oP6qR9sT2uV5wX8yZ1aB4cD7eF0gH3iJ6kL9mN2oP5qR8sT1uV4wX7yZ0aB3cD6eF9gH2iJ5kL8mN1oP4qR7sT0uV3wX6yZ9aB2cD5eF8gH1iJ4kL7mN0oP3qR6sT9uV2wX5yZ8aB1cD4eF7gH0iJ3kL6mN9oP2qR5sT8uV1wX4yZ7aB0cD3eF6gH9iJ2kL5mN8oP1qR4sT7uV0wX3yZ6aB9cD2eF5gH8iJ1kL4mN7oP0qR3sT6uV9wX2yZ5aB8cD1eF4gH7iJ0kL3mN6oP9qR2sT5uV8wX1yZ4aB7cD0eF3gH6iJ9kL2mN5oP8qR1sT4uV7wX0yZ3aB6cD9eF2gH5iJ8kL1mN4oP7qR0sT3uV6wX9yZ2aB5cD8eF1gH4iJ7kL0mN3oP6qR9sT2uV5wX8yZ1aB4cD7eF0gH3iJ6kL9mN2oP5qR8sT1uV4wX7yZ0aB3cD6eF9gH2iJ5kL8mN1oP4qR7sT0uV3wX6yZ9aB2cD5eF8gH1iJ4kL7mN0oP3qR6sT9uV2wX5yZ8aB1cD4eF7gH0iJ3kL6mN9oP2qR5sT8uV1wX4yZ7aB0cD3eF6gH9iJ2kL5mN8oP1qR4sT7uV0wX3yZ6aB9cD2eF5gH8iJ1kL4mN7oP0qR3sT6uV9wX2yZ5aB8cD1eF4gH7iJ0kL3mN6oP9qR2sT5uV8wX1yZ4aB7cD0eF3gH6iJ9kL2mN5oP8qR1sT4uV7wX0yZ3aB6cD9eF2gH5iJ8kL1mN4oP7qR0sT3uV6wX9yZ2aB5cD8eF1gH4iJ7kL0mN3oP6qR9sT2uV5wX8yZ1aB4cD7eF0gH3iJ6kL9mN2oP5qR8sT1uV4wX7yZ0aB3cD6eF9gH2iJ5kL8mN1oP4qR7sT0uV3wX6yZ9aB2cD5eF8gH1iJ4kL7mN0oP3qR6sT9uV2wX5yZ8aB1cD4eF7gH0iJ3
```

**⚠️ IMPORTANTE:** Guarde este secret! Você vai precisar dele no CLI.

### 1.4. Editar o arquivo `.env`

Abra o arquivo `.env` e configure as variáveis obrigatórias:

```env
# ==============================================
# SERVIDOR
# ==============================================
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# ==============================================
# BANCO DE DADOS
# ==============================================
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=sua-senha-aqui
DB_NAME=sagep_auth
DB_SSLMODE=disable

# ==============================================
# JWT
# ==============================================
JWT_SECRET=sua-chave-jwt-secreta-aqui
JWT_ACCESS_EXPIRATION=24h
JWT_REFRESH_EXPIRATION=168h

# ==============================================
# CORS
# ==============================================
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173

# ==============================================
# DOCS
# ==============================================
DOCS_ENABLED=true
DOCS_PATH=/docs

# ==============================================
# BOOTSTRAP SECRET (CRIAÇÃO INICIAL DE APLICAÇÕES)
# ==============================================
# ⚠️ Cole aqui o secret gerado no passo 1.3
BOOTSTRAP_SECRET=Kx9mP2vQ8nR5tY7wZ3aB6cD9eF1gH4iJ7kL0mN2oP5qR8sT1uV4wX7yZ0aB3cD6eF9gH2iJ5kL8mN1oP4qR7sT0uV3wX6yZ9aB2cD5eF8gH1iJ4kL7mN0oP3qR6sT9uV2wX5yZ8aB1cD4eF7gH0iJ3kL6mN9oP2qR5sT8uV1wX4yZ7aB0cD3eF6gH9iJ2kL5mN8oP1qR4sT7uV0wX3yZ6aB9cD2eF5gH8iJ1kL4mN7oP0qR3sT6uV9wX2yZ5aB8cD1eF4gH7iJ0kL3mN6oP9qR2sT5uV8wX1yZ4aB7cD0eF3gH6iJ9kL2mN5oP8qR1sT4uV7wX0yZ3aB6cD9eF2gH5iJ8kL1mN4oP7qR0sT3uV6wX9yZ2aB5cD8eF1gH4iJ7kL0mN3oP6qR9sT2uV5wX8yZ1aB4cD7eF0gH3iJ6kL9mN2oP5qR8sT1uV4wX7yZ0aB3cD6eF9gH2iJ5kL8mN1oP4qR7sT0uV3wX6yZ9aB2cD5eF8gH1iJ4kL7mN0oP3qR6sT9uV2wX5yZ8aB1cD4eF7gH0iJ3kL6mN9oP2qR5sT8uV1wX4yZ7aB0cD3eF6gH9iJ2kL5mN8oP1qR4sT7uV0wX3yZ6aB9cD2eF5gH8iJ1kL4mN7oP0qR3sT6uV9wX2yZ5aB8cD1eF4gH7iJ0kL3mN6oP9qR2sT5uV8wX1yZ4aB7cD0eF3gH6iJ9kL2mN5oP8qR1sT4uV7wX0yZ3aB6cD9eF2gH5iJ8kL1mN4oP7qR0sT3uV6wX9yZ2aB5cD8eF1gH4iJ7kL0mN3oP6qR9sT2uV5wX8yZ1aB4cD7eF0gH3iJ6kL9mN2oP5qR8sT1uV4wX7yZ0aB3cD6eF9gH2iJ5kL8mN1oP4qR7sT0uV3wX6yZ9aB2cD5eF8gH1iJ4kL7mN0oP3qR6sT9uV2wX5yZ8aB1cD4eF7gH0iJ3
```

### 1.5. Aplicar migrations no banco

```bash
# Certifique-se de que o banco está rodando e as migrations serão aplicadas automaticamente ao iniciar o servidor
# Ou execute manualmente se necessário
```

### 1.6. Iniciar o servidor

```bash
# Desenvolvimento
go run cmd/api/main.go

# Ou via Makefile (se disponível)
make run
```

**✅ Verificação:** Acesse `http://localhost:8080/health` e confirme que retorna `{"status":"ok"}`

---

## 🛠️ PASSO 2: Configurar o CLI `sagep-auth-cli`

### 2.1. Clonar e entrar no diretório

```bash
cd ~/source/BrBit/sagep-auth-cli
```

### 2.2. Criar arquivo `.env` a partir do exemplo

```bash
cp .env.example .env
```

### 2.3. Editar o arquivo `.env`

Abra o arquivo `.env` e configure:

```env
# ==============================================
# URL DO SERVIÇO SAGEP-AUTH
# ==============================================
SAGEP_AUTH_URL=http://localhost:8080

# ==============================================
# AUTENTICAÇÃO (BOOTSTRAP)
# ==============================================
# ⚠️ Cole aqui o MESMO secret gerado no passo 1.3
# Deve ser idêntico ao BOOTSTRAP_SECRET do servidor
SAGEP_AUTH_SECRET=Kx9mP2vQ8nR5tY7wZ3aB6cD9eF1gH4iJ7kL0mN2oP5qR8sT1uV4wX7yZ0aB3cD6eF9gH2iJ5kL8mN1oP4qR7sT0uV3wX6yZ9aB2cD5eF8gH1iJ4kL7mN0oP3qR6sT9uV2wX5yZ8aB1cD4eF7gH0iJ3kL6mN9oP2qR5sT8uV1wX4yZ7aB0cD3eF6gH9iJ2kL5mN8oP1qR4sT7uV0wX3yZ6aB9cD2eF5gH8iJ1kL4mN7oP0qR3sT6uV9wX2yZ5aB8cD1eF4gH7iJ0kL3mN6oP9qR2sT5uV8wX1yZ4aB7cD0eF3gH6iJ9kL2mN5oP8qR1sT4uV7wX0yZ3aB6cD9eF2gH5iJ8kL1mN4oP7qR0sT3uV6wX9yZ2aB5cD8eF1gH4iJ7kL0mN3oP6qR9sT2uV5wX8yZ1aB4cD7eF0gH3iJ6kL9mN2oP5qR8sT1uV4wX7yZ0aB3cD6eF9gH2iJ5kL8mN1oP4qR7sT0uV3wX6yZ9aB2cD5eF8gH1iJ4kL7mN0oP3qR6sT9uV2wX5yZ8aB1cD4eF7gH0iJ3kL6mN9oP2qR5sT8uV1wX4yZ7aB0cD3eF6gH9iJ2kL5mN8oP1qR4sT7uV0wX3yZ6aB9cD2eF5gH8iJ1kL4mN7oP0qR3sT6uV9wX2yZ5aB8cD1eF4gH7iJ0kL3mN6oP9qR2sT5uV8wX1yZ4aB7cD0eF3gH6iJ9kL2mN5oP8qR1sT4uV7wX0yZ3aB6cD9eF2gH5iJ8kL1mN4oP7qR0sT3uV6wX9yZ2aB5cD8eF1gH4iJ7kL0mN3oP6qR9sT2uV5wX8yZ1aB4cD7eF0gH3iJ6kL9mN2oP5qR8sT1uV4wX7yZ0aB3cD6eF9gH2iJ5kL8mN1oP4qR7sT0uV3wX6yZ9aB2cD5eF8gH1iJ4kL7mN0oP3qR6sT9uV2wX5yZ8aB1cD4eF7gH0iJ3
```

**⚠️ IMPORTANTE:** O `SAGEP_AUTH_SECRET` deve ser **exatamente igual** ao `BOOTSTRAP_SECRET` do servidor!

### 2.4. Compilar o CLI (opcional, para uso global)

```bash
# Opção 1: Build local
go build -o sagep-auth-cli ./cmd/sagep-auth-cli

# Opção 2: Instalar globalmente
go install ./cmd/sagep-auth-cli

# Opção 3: Usar diretamente com go run
# (não precisa compilar)
```

**✅ Verificação:** Teste o CLI

```bash
# Se instalou globalmente
sagep-auth-cli --help

# Se não instalou, use go run
go run ./cmd/sagep-auth-cli/main.go --help
```

---

## 📦 PASSO 3: Criar o Manifest na Aplicação `sagep-biopass-admin`

### 3.1. Entrar no diretório da aplicação

```bash
cd ~/source/BrBit/sagep-biopass-admin
```

### 3.2. Criar o arquivo `auth-manifest.yaml`

Crie o arquivo `auth-manifest.yaml` na raiz do projeto:

```yaml
# ============================================================================
# Manifest de Autenticação e Autorização - SAGEP BioPass Admin
# ============================================================================
# Este arquivo descreve a aplicação, suas permissões e roles base para o
# sistema de autenticação centralizado (sagep-auth).
# ============================================================================

application:
  code: sagep-biopass
  name: SAGEP Biopass
  description: Sistema de controle de ponto biométrico

# ============================================================================
# Permissões
# ============================================================================
# Cada permissão representa uma ação ou acesso a um recurso.
# Formato: biopass.{recurso}.{ação}
# ============================================================================

permissions:
  # Dashboard
  - code: biopass.dashboard.view
    description: Visualizar o dashboard principal

  # Dispositivos (CRUD completo)
  - code: biopass.devices.read
    description: Listar e visualizar dispositivos
  - code: biopass.devices.create
    description: Criar novos dispositivos
  - code: biopass.devices.update
    description: Editar dispositivos existentes
  - code: biopass.devices.delete
    description: Remover ou inativar dispositivos

  # Participantes (CRUD completo)
  - code: biopass.participants.read
    description: Listar e visualizar participantes
  - code: biopass.participants.create
    description: Criar novos participantes
  - code: biopass.participants.update
    description: Editar participantes existentes
  - code: biopass.participants.delete
    description: Remover ou inativar participantes

  # Registros de Ponto
  - code: biopass.attendance.read
    description: Visualizar registros de ponto/atendimento
  - code: biopass.attendance.create
    description: Criar registros de ponto manualmente
  - code: biopass.attendance.update
    description: Editar registros de ponto
  - code: biopass.attendance.delete
    description: Remover registros de ponto

  # Relatórios
  - code: biopass.reports.read
    description: Visualizar relatórios e análises

  # Usuários (administração dentro do contexto biopass)
  - code: biopass.users.read
    description: Listar e visualizar usuários do sistema
  - code: biopass.users.manage
    description: Criar, editar e gerenciar usuários e suas roles dentro da aplicação

  # Configurações
  - code: biopass.settings.view
    description: Acessar configurações da aplicação
  - code: biopass.settings.manage
    description: Gerenciar configurações da aplicação

  # ============================================================================
  # Permissões de Menu (CASL)
  # ============================================================================
  # Permissões específicas para controle de visibilidade de menus
  # Formato: Menu:{NomeDoMenu}
  # ============================================================================

  # Menus - Visibilidade
  - code: Menu:Dashboard
    description: Exibir menu Dashboard
  - code: Menu:Devices
    description: Exibir menu Dispositivos
  - code: Menu:Participants
    description: Exibir menu Participantes
  - code: Menu:Attendance
    description: Exibir menu Registros de Ponto
  - code: Menu:Reports
    description: Exibir menu Relatórios
  - code: Menu:Users
    description: Exibir menu Usuários
  - code: Menu:Settings
    description: Exibir menu Configurações

# ============================================================================
# Roles Base do Sistema
# ============================================================================
# Roles com system: true são criadas automaticamente pelo sagep-auth
# e não podem ser deletadas ou modificadas diretamente.
# ============================================================================

roles:
  # Role de Administrador Completo
  - code: biopass.admin
    name: Administrador BioPass
    system: true
    description: Acesso completo a todas as funcionalidades do BioPass
    permissions:
      # Dashboard
      - biopass.dashboard.view
      # Dispositivos
      - biopass.devices.read
      - biopass.devices.create
      - biopass.devices.update
      - biopass.devices.delete
      # Participantes
      - biopass.participants.read
      - biopass.participants.create
      - biopass.participants.update
      - biopass.participants.delete
      # Registros de Ponto
      - biopass.attendance.read
      - biopass.attendance.create
      - biopass.attendance.update
      - biopass.attendance.delete
      # Relatórios
      - biopass.reports.read
      # Usuários
      - biopass.users.read
      - biopass.users.manage
      # Configurações
      - biopass.settings.view
      - biopass.settings.manage
      # Menus
      - Menu:Dashboard
      - Menu:Devices
      - Menu:Participants
      - Menu:Attendance
      - Menu:Reports
      - Menu:Users
      - Menu:Settings

  # Role de Usuário Operacional
  - code: biopass.operator
    name: Operador BioPass
    system: true
    description: Acesso para operações do dia a dia (visualizar e criar registros)
    permissions:
      # Dashboard
      - biopass.dashboard.view
      # Dispositivos (somente leitura)
      - biopass.devices.read
      # Participantes (somente leitura)
      - biopass.participants.read
      # Registros de Ponto (criar e visualizar)
      - biopass.attendance.read
      - biopass.attendance.create
      # Relatórios (somente leitura)
      - biopass.reports.read
      # Menus
      - Menu:Dashboard
      - Menu:Devices
      - Menu:Participants
      - Menu:Attendance
      - Menu:Reports

  # Role de Visualizador/Relatórios
  - code: biopass.viewer
    name: Visualizador BioPass
    system: true
    description: Acesso somente leitura para visualização de dados e relatórios
    permissions:
      # Dashboard
      - biopass.dashboard.view
      # Dispositivos (somente leitura)
      - biopass.devices.read
      # Participantes (somente leitura)
      - biopass.participants.read
      # Registros de Ponto (somente leitura)
      - biopass.attendance.read
      # Relatórios
      - biopass.reports.read
      # Menus
      - Menu:Dashboard
      - Menu:Devices
      - Menu:Participants
      - Menu:Attendance
      - Menu:Reports
```

### 3.3. Verificar se o arquivo está correto

```bash
# Verificar sintaxe YAML (se tiver yamllint instalado)
yamllint auth-manifest.yaml

# Ou simplesmente abrir e verificar manualmente
cat auth-manifest.yaml
```

---

## 🔄 PASSO 4: Sincronizar o Manifest com o Servidor

### 4.1. Certificar-se de que o servidor está rodando

```bash
# No terminal do sagep-auth, verifique se está rodando
curl http://localhost:8080/health
# Deve retornar: {"status":"ok","service":"sagep-auth","version":"1.0.0"}
```

### 4.2. Executar o sync (bootstrap inicial)

**Opção A: Usando o CLI compilado/instalado**

```bash
# Se instalou globalmente
sagep-auth-cli sync --manifest ./auth-manifest.yaml

# Ou se compilou localmente
./sagep-auth-cli sync --manifest ./auth-manifest.yaml
```

**Opção B: Usando go run (sem compilar)**

```bash
# Do diretório sagep-auth-cli
go run ./cmd/sagep-auth-cli/main.go sync --manifest ~/source/BrBit/sagep-biopass-admin/auth-manifest.yaml
```

**Opção C: Usando variáveis de ambiente (sem .env)**

```bash
export SAGEP_AUTH_URL=http://localhost:8080
export SAGEP_AUTH_SECRET=Kx9mP2vQ8nR5tY7wZ3aB6cD9eF1gH4iJ7kL0mN2oP5qR8sT1uV4wX7yZ0aB3cD6eF9gH2iJ5kL8mN1oP4qR7sT0uV3wX6yZ9aB2cD5eF8gH1iJ4kL7mN0oP3qR6sT9uV2wX5yZ8aB1cD4eF7gH0iJ3kL6mN9oP2qR5sT8uV1wX4yZ7aB0cD3eF6gH9iJ2kL5mN8oP1qR4sT7uV0wX3yZ6aB9cD2eF5gH8iJ1kL4mN7oP0qR3sT6uV9wX2yZ5aB8cD1eF4gH7iJ0kL3mN6oP9qR2sT5uV8wX1yZ4aB7cD0eF3gH6iJ9kL2mN5oP8qR1sT4uV7wX0yZ3aB6cD9eF2gH5iJ8kL1mN4oP7qR0sT3uV6wX9yZ2aB5cD8eF1gH4iJ7kL0mN3oP6qR9sT2uV5wX8yZ1aB4cD7eF0gH3iJ6kL9mN2oP5qR8sT1uV4wX7yZ0aB3cD6eF9gH2iJ5kL8mN1oP4qR7sT0uV3wX6yZ9aB2cD5eF8gH1iJ4kL7mN0oP3qR6sT9uV2wX5yZ8aB1cD4eF7gH0iJ3

go run ./cmd/sagep-auth-cli/main.go sync --manifest ~/source/BrBit/sagep-biopass-admin/auth-manifest.yaml
```

### 4.3. Verificar o resultado

**Saída esperada (sucesso):**

```
✅ Sincronização concluída com sucesso!

📊 Resumo:
   - Aplicação: sagep-biopass (criada)
   - Permissions: 20 criadas
   - Roles: 3 criadas
   - Role Permissions: 45 vinculações criadas
```

**Se houver erro:**

- **401 Unauthorized**: Verifique se `BOOTSTRAP_SECRET` e `SAGEP_AUTH_SECRET` são idênticos
- **Connection refused**: Verifique se o servidor `sagep-auth` está rodando
- **Invalid manifest**: Verifique a sintaxe YAML do `auth-manifest.yaml`

---

## ✅ PASSO 5: Verificar no Servidor

### 5.1. Acessar a documentação (Redoc)

Abra no navegador:
```
http://localhost:8080/docs
```

### 5.2. Verificar via API (opcional)

```bash
# Listar aplicações (requer autenticação JWT)
# Primeiro, você precisa criar um usuário master e autenticar
curl -X GET http://localhost:8080/v1/applications \
  -H "Authorization: Bearer {seu-token-jwt}"
```

---

## 🔐 PASSO 6: Próximos Passos (Após Bootstrap)

### 6.1. Criar usuário Master (se ainda não existir)

Você precisará criar um usuário master manualmente no banco ou via script de seed.

### 6.2. Autenticar e obter token JWT

```bash
curl -X POST http://localhost:8080/v1/authenticate \
  -H "Content-Type: application/json" \
  -d '{
    "email": "master@sagep.com.br",
    "password": "sua-senha",
    "application_code": "sagep-biopass"
  }'
```

### 6.3. Usar JWT para próximos syncs (opcional)

Após ter um token JWT, você pode atualizar o `.env` do CLI:

```env
# Remover ou comentar SAGEP_AUTH_SECRET
# SAGEP_AUTH_SECRET=...

# Adicionar token JWT
SAGEP_AUTH_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

Agora os próximos syncs usarão JWT ao invés de HMAC.

---

## 📝 Resumo das Variáveis de Ambiente

### `sagep-auth/.env`
```env
BOOTSTRAP_SECRET=<secret-gerado-com-openssl>
```

### `sagep-auth-cli/.env`
```env
SAGEP_AUTH_URL=http://localhost:8080
SAGEP_AUTH_SECRET=<mesmo-secret-do-servidor>  # Para bootstrap
# OU
SAGEP_AUTH_TOKEN=<token-jwt>  # Para uso normal (após bootstrap)
```

---

## 🎯 Checklist Final

- [ ] Servidor `sagep-auth` configurado e rodando
- [ ] `BOOTSTRAP_SECRET` gerado e configurado no servidor
- [ ] CLI `sagep-auth-cli` configurado
- [ ] `SAGEP_AUTH_SECRET` configurado no CLI (mesmo valor do servidor)
- [ ] `auth-manifest.yaml` criado na aplicação
- [ ] Sync executado com sucesso
- [ ] Aplicação, permissions e roles criadas no banco

---

## 🆘 Troubleshooting

### Erro: "Bootstrap não configurado"
- Verifique se `BOOTSTRAP_SECRET` está definido no `.env` do servidor

### Erro: "Assinatura HMAC inválida"
- Verifique se `BOOTSTRAP_SECRET` e `SAGEP_AUTH_SECRET` são **exatamente iguais**
- Verifique se não há espaços extras ou quebras de linha

### Erro: "Timestamp muito antigo"
- O timestamp tem validade de 5 minutos
- Execute o sync novamente

### Erro: "Connection refused"
- Verifique se o servidor `sagep-auth` está rodando
- Verifique se a URL está correta (`SAGEP_AUTH_URL`)

---

**🎉 Pronto!** Seu sistema de autenticação está configurado e sincronizado!

