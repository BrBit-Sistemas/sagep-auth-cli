# sagep-auth-cli

CLI para sincronização de manifests com o serviço `sagep-auth`.

## 📦 Instalação

```bash
git clone <repo>
cd sagep-auth-cli
go build -o sagep-auth-cli ./cmd/sagep-auth-cli
```

## ⚙️ Configuração

Configure as variáveis de ambiente:

```bash
# Obrigatório
export SAGEP_AUTH_URL=http://localhost:8080

# Para bootstrap (criação inicial)
export SAGEP_AUTH_SECRET=your-secret-here

# Para uso normal (após bootstrap)
export SAGEP_AUTH_TOKEN=your-jwt-token
```

Ou crie um arquivo `.env`:

```env
SAGEP_AUTH_URL=http://localhost:8080
SAGEP_AUTH_SECRET=your-secret-here
```

## 🚀 Comandos

### `init` - Criar manifest interativamente

Cria um novo `auth-manifest.yaml` guiando você passo a passo.

```bash
./sagep-auth-cli init
./sagep-auth-cli --manifest ./meu-manifest.yaml init
```

### `sync` - Sincronizar manifest

Envia o manifest para o servidor `sagep-auth`.

```bash
./sagep-auth-cli sync
./sagep-auth-cli --manifest ./auth-manifest.yaml sync
```

## 📚 Documentação

- **Guia Completo:** `docs/GUIA_COMPLETO.md` - Passo a passo completo
- **Regras de Negócio:** `docs/REGRAS_NEGOCIO.md` - Regras e comportamentos

## 📊 Exemplo de Saída

```bash
Sincronizando aplicação: sagep-biopass
URL do auth: http://localhost:8080

Application: sagep-biopass (created)
Permissions: 25 (25 criadas, 0 atualizadas)
Roles:       4 (4 criadas, 0 atualizadas)
Users:       2 (2 criados, 0 atualizados)

Sync concluído com sucesso.
```

## 📝 Exemplo de Manifest

```yaml
application:
  code: sagep-biopass
  name: SAGEP Biopass
  description: Sistema de controle de ponto

permissions:
  - code: biopass.devices.read
    description: Listar dispositivos

roles:
  - code: biopass.admin
    name: Administrador
    system: true
    permissions:
      - biopass.*

users:
  - email: master@sagep.com.br
    password: Master@123  # Senha em texto claro (será hasheada pelo servidor)
    name: Master Admin
    roles:
      - master
    active: true
```

Veja `auth-manifest.example.yaml` para exemplo completo.
