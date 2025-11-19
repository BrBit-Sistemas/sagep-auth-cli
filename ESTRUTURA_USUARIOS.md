# 📋 Estrutura YAML para Usuários no Manifest

## Proposta de Estrutura

Adicione uma seção `users:` no `auth-manifest.yaml` com a seguinte estrutura:

```yaml
application:
  code: sagep-biopass
  name: SAGEP Biopass
  description: Sistema de controle de ponto biométrico

permissions:
  # ... suas permissões aqui ...

roles:
  # ... suas roles aqui ...
  - code: biopass.admin
    name: Administrador BioPass
    system: true
    description: Acesso completo
    permissions:
      - biopass.*

users:
  # Usuário Master (acesso total)
  - email: master@sagep.com.br
    password: Master@123  # Senha em texto claro (será hasheada pelo servidor)
    name: Master Admin
    roles:
      - master  # Role code, não ID (será resolvido pelo servidor)
    active: true  # Opcional, default: true

  # Usuário comum (exemplo)
  - email: user@sagep.com.br
    password: User@123
    name: Usuário Exemplo
    roles:
      - biopass.admin  # Role code definida no manifest acima
      # Pode ter múltiplas roles
    active: true
```

## 📝 Campos Obrigatórios e Opcionais

### Campo `users[]` (array de usuários)

| Campo | Tipo | Obrigatório | Descrição |
|-------|------|-------------|-----------|
| `email` | string | ✅ | Email único do usuário (único globalmente no sistema) |
| `password` | string | ✅ | Senha em texto claro (será hasheada com bcrypt pelo servidor) |
| `name` | string | ✅ | Nome completo do usuário |
| `roles` | string[] | ✅ | Lista de **códigos de roles** (não IDs). Ex: `["master"]`, `["biopass.admin"]` |
| `active` | boolean | ❌ | Status ativo/inativo (default: `true`) |

## ⚠️ Regras Importantes

1. **Email único**: O email deve ser único globalmente. Se o usuário já existir, será atualizado (upsert).
2. **Senha em texto claro**: A senha é enviada em texto claro no manifest e será hasheada pelo servidor com bcrypt.
3. **Roles por código**: Use os **códigos das roles** (ex: `"master"`, `"biopass.admin"`), não os IDs. O servidor resolve os códigos para IDs.
4. **Vinculação automática**: O usuário é automaticamente vinculado à aplicação definida em `application.code`.
5. **Idempotência**: Executar sync múltiplas vezes com os mesmos usuários não cria duplicatas (upsert por email).

## 🔐 Fluxo de Criação

1. **Cria/atualiza usuário** na tabela `users` (upsert por email)
2. **Vincula usuário à aplicação** (cria `user_applications` se não existir)
3. **Atribui roles ao usuário** (cria/atualiza `user_roles` baseado nos códigos fornecidos)

## 📌 Exemplo Completo

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

roles:
  - code: master
    name: Master
    system: true
    description: Acesso total ao sistema
    permissions:
      - biopass.*
  
  - code: biopass.admin
    name: Administrador BioPass
    system: true
    description: Acesso completo ao BioPass
    permissions:
      - biopass.*

  - code: biopass.operator
    name: Operador BioPass
    system: true
    description: Acesso operacional
    permissions:
      - biopass.devices.read
      - biopass.devices.create

users:
  # Usuário Master
  - email: master@sagep.com.br
    password: Master@123
    name: Master Admin
    roles:
      - master

  # Usuário Admin comum
  - email: admin@sagep.com.br
    password: Admin@123
    name: Administrador
    roles:
      - biopass.admin

  # Usuário Operador
  - email: operador@sagep.com.br
    password: Operador@123
    name: Operador Sistema
    roles:
      - biopass.operator
```

## 🔄 Comportamento do Sync

### Se usuário não existe:
- ✅ Cria novo usuário
- ✅ Vincula à aplicação
- ✅ Atribui roles

### Se usuário já existe (mesmo email):
- ✅ Atualiza nome e senha (se mudou)
- ✅ Mantém vínculo com aplicação (ou cria se não existir)
- ✅ Atualiza roles (remove roles antigas não listadas, adiciona novas)

## 🚨 Segurança

- **Senhas**: Sempre use senhas fortes. O manifest deve ser tratado como **informação sensível**.
- **Versionamento**: Considere manter senhas de desenvolvimento/teste diferentes das de produção.
- **Git**: Considere usar `.gitignore` ou variáveis de ambiente para senhas em produção.

