# Guia Completo - sagep-auth-cli

## O que é?

CLI para integrar aplicações com o `sagep-auth`, sincronizando aplicação, permissões, roles base e usuários iniciais.

## Passo a Passo

### 1. Configurar Ambiente

```bash
# Clonar repositório
git clone <repo>
cd sagep-auth-cli

# Compilar
go build -o sagep-auth-cli ./cmd/sagep-auth-cli

# Configurar .env
cat > .env << EOF
SAGEP_AUTH_URL=http://localhost:8080
SAGEP_AUTH_SECRET=$(openssl rand -base64 32)
EOF
```

⚠️ **Importante:** Configure o mesmo `SAGEP_AUTH_SECRET` no servidor `sagep-auth`.

### 2. Criar Manifest

#### Opção A: Interativo (Recomendado)

```bash
./sagep-auth-cli init
```

O wizard pergunta:
- Código e nome da aplicação
- Deseja criar usuários? (Master e comuns)
- Permissões
- Roles base

#### Opção B: Manual

Copie `auth-manifest.example.yaml` e edite conforme necessário.

### 3. Sincronizar

```bash
./sagep-auth-cli sync
```

O CLI:
1. Carrega o manifest
2. Autentica (HMAC ou JWT)
3. Envia para `/v1/applications/sync`
4. Exibe resultado

### 4. Verificar

Acesse o servidor `sagep-auth` e verifique:
- Aplicação criada
- Permissões sincronizadas
- Roles base criadas
- Usuários criados e vinculados

## Estrutura do Manifest

### `application`
- `code`: Identificador único (ex: `sagep-biopass`)
- `name`: Nome amigável
- `description`: Descrição opcional

### `permissions`
- `code`: Código único (ex: `biopass.devices.read`)
- `description`: Descrição opcional

### `roles`
- `code`: Código único (ex: `biopass.admin`)
- `name`: Nome amigável
- `system`: `true` = role base (protegida), `false` = customizada
- `permissions`: Lista de códigos ou wildcards (ex: `biopass.*`)

### `users`
- `email`: Email único
- `password`: Senha em texto claro (será hasheada pelo servidor)
- `name`: Nome completo
- `roles`: Lista de códigos de roles
- `active`: Status (default: `true`)

## Autenticação

### Bootstrap (Inicial)

Usa HMAC com secret compartilhado:

```bash
export SAGEP_AUTH_SECRET=your-secret
./sagep-auth-cli sync
```

### Uso Normal

Após criar usuário, use JWT:

```bash
export SAGEP_AUTH_TOKEN=your-jwt-token
./sagep-auth-cli sync
```

## Fluxo Completo

```
1. Desenvolvedor cria nova app (ex: sagep-biopass)
   ↓
2. Executa: ./sagep-auth-cli init
   ↓
3. Preenche informações (app, permissões, roles, usuários)
   ↓
4. Executa: ./sagep-auth-cli sync (HMAC)
   ↓
5. Servidor cria: aplicação, permissões, roles, usuários
   ↓
6. Usuário faz login com credenciais do manifest
   ↓
7. Próximos syncs usam JWT (uso normal)
```

## Troubleshooting

**Erro: "SAGEP_AUTH_TOKEN ou SAGEP_AUTH_SECRET é obrigatório"**
- Configure `.env` ou variáveis de ambiente

**Erro: "401 Unauthorized"**
- Verifique se `SAGEP_AUTH_SECRET` está correto
- Para JWT, verifique se o token é válido

**Usuários não criados**
- Verifique se o manifest tem a seção `users:`
- Verifique logs do servidor

