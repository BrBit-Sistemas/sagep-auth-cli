# Formato do tenant_id no Manifest

## 📋 Visão Geral

O campo `tenant_id` no manifest pode ter dois formatos diferentes, dependendo do nível de acesso necessário do usuário.

## ✅ Formatos Suportados

### 1. UnidadeId (Guid)

**Formato:** UUID como string
**Uso:** Para usuários vinculados a uma unidade específica
**Exemplo:**
```yaml
users:
  - email: operador@unidade.com
    tenant_id: "550e8400-e29b-41d4-a716-446655440000"  # UnidadeId (Guid)
```

**Comportamento:**
- Backend interpreta como UnidadeId
- Usuário vê apenas dados da sua unidade
- Pode ter acesso de Regional (mas tenant_id ainda é da Unidade)

### 2. SecretariaTenantId (string)

**Formato:** String identificadora
**Uso:** Para usuários Master/Admin de Secretaria
**Exemplo:**
```yaml
users:
  - email: master@sagep.com.br
    tenant_id: "sc-sejuc"  # SecretariaTenantId (string)
```

**Comportamento:**
- Backend interpreta como SecretariaTenantId
- Usuário vê dados de todas unidades da secretaria
- Deve ter role `master`, `core_admin` ou `core_gestor_estrutura`

### 3. Sem tenant_id (Omitido)

**Formato:** Campo omitido ou `null`
**Uso:** Para usuários globais (sem multi-tenancy)
**Exemplo:**
```yaml
users:
  - email: global@example.com
    # tenant_id omitido
```

**Comportamento:**
- Usuário não usa multi-tenancy
- Campo `tenant_id` fica `NULL` no banco

## 🎯 Quando Usar Cada Formato

| Formato | Quando Usar | Exemplo |
|---------|-------------|---------|
| **UnidadeId (Guid)** | Usuário de unidade específica | `"550e8400-e29b-41d4-a716-446655440000"` |
| **SecretariaTenantId (string)** | Master/Admin de Secretaria | `"sc-sejuc"` |
| **Omitido** | Usuário global | (campo não presente) |

## 📝 Exemplos Completos

### Exemplo 1: Master de Secretaria
```yaml
users:
  - email: master@sagep.com.br
    password: Master@123
    name: Master Admin
    tenant_id: "sc-sejuc"  # SecretariaTenantId (string)
    roles:
      - master
```

### Exemplo 2: Usuário de Unidade
```yaml
users:
  - email: operador@unidade.com
    password: Operador@123
    name: Operador da Unidade
    tenant_id: "550e8400-e29b-41d4-a716-446655440000"  # UnidadeId (Guid)
    roles:
      - biopass.user
```

### Exemplo 3: Usuário Global
```yaml
users:
  - email: global@example.com
    password: Global@123
    name: Usuário Global
    # tenant_id omitido
    roles:
      - system_admin
```

## ⚠️ Observações Importantes

1. **Apenas para novos usuários:** `tenant_id` no manifest só é aplicado na criação de novos usuários. Usuários existentes não têm seu `tenant_id` atualizado via sync.

2. **Valor correto:** O `SecretariaTenantId` deve ser o mesmo valor do campo `TenantId` da entidade `Secretaria` no `core-api`.

3. **Compatibilidade:** Backend tenta interpretar ambos os formatos automaticamente (Guid primeiro, depois string).

## 🔗 Referências

- [Regras de Negócio - tenant_id Format](../../sagep-auth/docs/business-rules/auth/tenant-id-format.md)
- [Multi-Tenancy Hierárquico](../../sagep-core-api/docs/MASTER_ROLE_MULTI_TENANCY.md)





