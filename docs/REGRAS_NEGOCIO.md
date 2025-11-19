# Regras de Negócio - sagep-auth-cli

## Manifest

### Idempotência
- Executar `sync` múltiplas vezes produz o mesmo resultado
- Se aplicação/role/permissão já existe, é atualizada (não duplicada)

### Aplicação
- `code` deve ser único globalmente
- Usado como identificador em todos os relacionamentos

### Permissões
- `code` deve ser único por aplicação
- Aceita qualquer formato/nomenclatura
- Wildcards (`biopass.*`) são expandidos no servidor

### Roles
- `code` deve ser único por aplicação
- `system: true` = role base (não editável via API, apenas sync)
- `system: false` = role customizada (editável via API)
- Permissões podem usar wildcards

### Usuários
- `email` deve ser único globalmente
- Se usuário existe, é atualizado (nome, senha)
- Senha em texto claro no YAML → hasheada pelo servidor
- Vinculado automaticamente à aplicação do manifest
- Roles resolvidas por código (não ID)

## Autenticação

### HMAC (Bootstrap)
- Usado quando não há aplicação/usuário ainda
- Requer `BOOTSTRAP_SECRET` no servidor
- Requer `SAGEP_AUTH_SECRET` no CLI
- Assinatura: `HMAC-SHA256(body + timestamp, secret)`

### JWT (Uso Normal)
- Usado após bootstrap inicial
- Token obtido via `/v1/authenticate`
- Carrega `application_id` do token

## Sincronização

### Ordem de Processamento
1. Aplicação (upsert por `code`)
2. Permissões (upsert por `application_id + code`)
3. Roles (upsert por `application_id + code`)
4. Role-Permissions (regenerado baseado no manifest)
5. Usuários (upsert por `email`)
6. User-Application (vinculação)
7. User-Roles (atualizado baseado nos códigos)

### Comportamento
- Criar: Se não existe, cria novo
- Atualizar: Se existe, atualiza campos (exceto IDs)
- Ignorar: Se erro em um item, continua com próximo

### Usuários
- Criar: Tabela `users` → `user_applications` → `user_roles`
- Atualizar: Se email existe, atualiza `users`, mantém vínculos
- Senha: Sempre hasheada pelo servidor (bcrypt)

## Segurança

- Senhas nunca são salvas em texto claro
- Manifest com senhas deve ser tratado como sensível
- HMAC previne replay attacks (timestamp validado)
- Roles base (`system: true`) protegidas contra edição via API

