# 🚀 Bootstrap Completo - SAGEP Auth

Guia passo a passo para configurar e sincronizar o sistema de autenticação.

---

## 📋 Pré-requisitos

- Go 1.21+
- PostgreSQL rodando
- Repositórios: `sagep-auth`, `sagep-auth-cli`, `sagep-biopass-admin`

---

## 🔧 PASSO 1: Configurar `sagep-auth` (Servidor)

### 1.1. Clonar e configurar

```bash
cd ~/source/BrBit/sagep-auth
cp env.example .env
```

### 1.2. Gerar Secret HMAC

```bash
openssl rand -base64 32
```

**⚠️ IMPORTANTE:** Guarde este secret! Você precisará dele no CLI.

### 1.3. Editar `.env`

```env
# Servidor
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Banco
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=sua-senha
DB_NAME=sagep_auth
DB_SSLMODE=disable

# JWT
JWT_SECRET=sua-chave-jwt-secreta
JWT_ACCESS_EXPIRATION=24h
JWT_REFRESH_EXPIRATION=168h

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173

# Docs
DOCS_ENABLED=true
DOCS_PATH=/docs

# Bootstrap (cole o secret gerado no passo 1.2)
BOOTSTRAP_SECRET=seu-secret-aqui
```

### 1.4. Iniciar servidor

```bash
go run cmd/api/main.go
```

**✅ Verificar:** `http://localhost:8080/health` deve retornar `{"status":"ok"}`

---

## 🛠️ PASSO 2: Configurar `sagep-auth-cli`

### 2.1. Clonar e configurar

```bash
cd ~/source/BrBit/sagep-auth-cli
cp .env.example .env
```

### 2.2. Editar `.env`

```env
SAGEP_AUTH_URL=http://localhost:8080

# ⚠️ Mesmo secret do servidor (passo 1.2)
SAGEP_AUTH_SECRET=seu-secret-aqui
```

### 2.3. Compilar (opcional)

```bash
go build -o sagep-auth-cli ./cmd/sagep-auth-cli
```

---

## 📦 PASSO 3: Criar Manifest na Aplicação

### 3.1. Criar `auth-manifest.yaml` na raiz do projeto

```bash
cd ~/source/BrBit/sagep-biopass-admin
```

Criar arquivo `auth-manifest.yaml`:

```yaml
application:
  code: sagep-biopass
  name: SAGEP Biopass
  description: Sistema de controle de ponto biométrico

permissions:
  - code: biopass.dashboard.view
    description: Visualizar dashboard
  - code: biopass.devices.read
    description: Listar dispositivos
  - code: biopass.devices.create
    description: Criar dispositivos
  - code: biopass.devices.update
    description: Editar dispositivos
  - code: biopass.devices.delete
    description: Remover dispositivos
  # ... adicione todas as permissões necessárias

roles:
  - code: biopass.admin
    name: Administrador BioPass
    system: true
    description: Acesso completo
    permissions:
      - biopass.*  # Wildcard: todas as permissões que começam com biopass.
  # ... adicione todas as roles base
```

---

## 🔄 PASSO 4: Sincronizar Manifest

### 4.1. Executar sync

**⚠️ IMPORTANTE:** Flags devem vir ANTES do comando `sync`

```bash
# Do diretório sagep-auth-cli
./sagep-auth-cli --manifest ~/source/BrBit/sagep-biopass-admin/auth-manifest.yaml sync

# Ou se compilou globalmente
sagep-auth-cli --manifest ~/source/BrBit/sagep-biopass-admin/auth-manifest.yaml sync

# Ou usando go run
go run ./cmd/sagep-auth-cli/main.go --manifest ~/source/BrBit/sagep-biopass-admin/auth-manifest.yaml sync
```

**✅ Saída esperada:**
```
Sincronizando aplicação: sagep-biopass
URL do auth: http://localhost:8080

Application: sagep-biopass (created)
Permissions: 20 (20 criadas, 0 atualizadas)
Roles: 3 (3 criadas, 0 atualizadas)

Sync concluído com sucesso.
```

---

## 🔐 PASSO 5: Criar Usuário Master (Opcional)

Após o sync, você precisa criar um usuário master manualmente para acessar o sistema:

### Opção A: Script SQL

```sql
-- Executar no PostgreSQL
INSERT INTO users (id, email, password, name, active, version, created_at, updated_at)
VALUES (
  gen_random_uuid(),
  'master@sagep.com.br',
  '$2a$10$...', -- Hash bcrypt da senha (gerar com Go)
  'Master Admin',
  true,
  1,
  NOW(),
  NOW()
);
```

### Opção B: Tool Go

```bash
cd ~/source/BrBit/sagep-auth
go run cmd/tools/setup_master/main.go
```

---

## ✅ Checklist

- [ ] `sagep-auth` rodando (`http://localhost:8080/health`)
- [ ] `BOOTSTRAP_SECRET` configurado no servidor
- [ ] `SAGEP_AUTH_SECRET` configurado no CLI (mesmo valor)
- [ ] `auth-manifest.yaml` criado na aplicação
- [ ] Sync executado com sucesso
- [ ] Usuário master criado (se necessário)

---

## 🆘 Troubleshooting

**Erro: "Os flags devem vir ANTES do comando"**
- ✅ Correto: `./sagep-auth-cli --manifest file.yaml sync`
- ❌ Errado: `./sagep-auth-cli sync --manifest file.yaml`

**Erro: "Assinatura HMAC inválida"**
- Verifique se `BOOTSTRAP_SECRET` e `SAGEP_AUTH_SECRET` são idênticos

**Erro: "Connection refused"**
- Verifique se o servidor `sagep-auth` está rodando

**Erro: "Timestamp muito antigo"**
- Execute o sync novamente (timestamp válido por 5 minutos)

---

**🎉 Pronto!** Sistema configurado e sincronizado.
