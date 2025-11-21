# Análise do YAML - Padrões e Recomendações

**Data:** 2025-01-XX  
**Objetivo:** Padronizar subjects para compatibilidade com frontend (CASL.js)

---

## ✅ Confirmações

### 1. Master com `permissions: []` ✅

**Status:** **CORRETO E CONFIRMADO**

O backend **realmente converte** `permissions: []` para:
```json
{
  "abilities": [
    { "action": "manage", "subject": "all" }
  ]
}
```

**Código de referência:** `get_user_info_usecase.go:86-92`

✅ **Não precisa mudar nada!**

---

## ⚠️ Problema Identificado

### Subject de `biopass.participants.read`

**Situação atual:**
- **Code:** `biopass.participants.read`
- **Subject no YAML:** `Participants` (singular, maiúsculo)
- **Frontend verifica:** `hasPermission('read', 'participantes')` (plural, minúsculo)

**O que acontece:**

1. **Se fornecido explicitamente** `subject: Participants`:
   - Backend usa exatamente: `subject: "Participants"`
   - CASL.js recebe: `{action: "read", subject: "Participants"}`
   - Frontend verifica: `ability.can('read', 'participantes')` ❌ **NÃO FUNCIONA!**

2. **Se deixar vazio** (backend infere):
   - Backend parseia: `biopass.participants.read`
   - Extrai: `participants` (plural, minúsculo)
   - **Capitaliza automaticamente:** `Participants` (singular, maiúsculo)
   - CASL.js recebe: `{action: "read", subject: "Participants"}`
   - Frontend verifica: `ability.can('read', 'participantes')` ❌ **NÃO FUNCIONA!**

**Problema:** O backend capitaliza e pode singularizar, mas o frontend espera minúsculo e plural.

---

## ✅ Solução

### Use `subject` explícito conforme o frontend espera

**O backend usa EXATAMENTE o que você fornecer no `subject`**, sem modificações.

**YAML Correto:**

```yaml
permissions:
  - code: biopass.participants.read
    subject: participantes  # ✅ Minúsculo, plural (como o frontend espera)
    action: read
    description: Acesso de leitura aos participantes
```

**Resultado:**
- Backend salva: `subject: "participantes"`
- CASL.js recebe: `{action: "read", subject: "participantes"}`
- Frontend verifica: `ability.can('read', 'participantes')` ✅ **FUNCIONA!**

---

## 📋 Padrão Recomendado

### Subjects devem seguir o padrão do frontend

**Regra de ouro:** O `subject` no YAML deve ser **exatamente** o que o frontend verifica no CASL.js.

#### Padrão por Tipo:

1. **Recursos (entidades):**
   ```yaml
   - code: biopass.participants.read
     subject: participantes  # ✅ Minúsculo, plural (se frontend usa plural)
     action: read
   
   - code: biopass.devices.read
     subject: devices  # ✅ Minúsculo, plural (se frontend usa plural)
     action: read
   ```

2. **Menus:**
   ```yaml
   - code: Menu:Dashboard
     subject: Menu:Dashboard  # ✅ Mantém formato Menu:{Nome}
     action: view
   
   - code: Menu:Participantes
     subject: Menu:Participantes  # ✅ Mantém formato Menu:{Nome}
     action: view
   ```

3. **Singular vs Plural:**
   - Se frontend usa `participantes` → use `participantes`
   - Se frontend usa `Participant` → use `Participant`
   - **Verifique o código do frontend!**

---

## 🔍 Como Verificar o Padrão do Frontend

### Método 1: Verificar código do frontend

Procure por verificações CASL no código:
```typescript
// Exemplo no frontend
ability.can('read', 'participantes')  // ✅ Frontend espera 'participantes'
ability.can('read', 'devices')        // ✅ Frontend espera 'devices'
```

### Método 2: Testar o endpoint `/me`

Após fazer sync, verifique o que o backend retorna:
```bash
curl -X GET http://auth-url/me \
  -H "Authorization: Bearer <token>"
```

Veja o campo `abilities`:
```json
{
  "abilities": [
    { "action": "read", "subject": "participantes" }  // ✅ Este é o que o frontend recebe
  ]
}
```

O `subject` que aparece aqui deve ser **exatamente** o que o frontend verifica.

---

## ✅ YAML Corrigido

```yaml
application:
  code: sagep-biopass
  name: Sagep Biopass
  description: Sistema de biometria facial do ecossistema SAGEP

permissions:
  # Menus
  - code: Menu:Dashboard
    subject: Menu:Dashboard
    action: view
    description: Acesso ao menu dashboard
  
  - code: Menu:Participantes
    subject: Menu:Participantes
    action: view
    description: Acesso ao menu participantes
  
  # Recursos - USAR EXATAMENTE O QUE O FRONTEND ESPERA
  - code: biopass.participants.read
    subject: participantes  # ✅ Minúsculo, plural (conforme frontend verifica)
    action: read
    description: Acesso de leitura aos participantes

roles:
  # Master - permissions vazio (backend converte automaticamente para "manage all")
  - code: master
    name: Master
    system: true
    description: Acesso total ao sistema
    permissions: []  # ✅ Correto - backend converte para {action: "manage", subject: "all"}
  
  # Outras roles
  - code: biopass.user
    name: Usuário
    system: true
    description: Acesso de usuário comum
    permissions:
      - Menu:Dashboard
      - Menu:Participantes
      - biopass.participants.read

users:
  - email: alan@brbitsistemas.com.br
    password: Bb2025!@
    name: Alan Rezende
    active: true
    roles:
      - master
  
  - email: alanrezendeee@gmail.com
    password: Bb2025!@
    name: Alan Rezende
    active: true
    roles: []
```

---

## 📝 Respostas às Ponderações

### 1. **Inconsistência no subject de participantes**

❌ **Problema:** `subject: Participants` (singular, maiúsculo)  
✅ **Solução:** `subject: participantes` (plural, minúsculo - como frontend espera)

**Motivo:** O backend usa **exatamente** o que você fornece no `subject`. Se o frontend verifica `'participantes'` (minúsculo, plural), use exatamente isso.

### 2. **Falta de permissão para o menu**

✅ **Já está correto no YAML!**  
- `Menu:Participantes` está definido na seção `permissions`
- `biopass.user` referencia corretamente

### 3. **Padrão de subjects**

✅ **Recomendação:** 
- Use **exatamente** o que o frontend verifica no CASL.js
- Se frontend usa `participantes` → use `participantes`
- Se frontend usa `Participant` → use `Participant`
- **Verifique o código do frontend ou teste o endpoint `/me`**

### 4. **Conversão do master**

✅ **Confirmado:** Backend converte `permissions: []` para `{action: "manage", subject: "all"}`  
✅ **Código:** `get_user_info_usecase.go:86-92`

### 5. **Code vs Subject**

✅ **São independentes:**
- `code`: Identificador único (ex: `biopass.participants.read`)
- `subject`: Recurso para CASL.js (ex: `participantes`)
- O `code` não precisa corresponder ao `subject`

---

## 🎯 Resumo Executivo

### ✅ O que está correto:

1. Master com `permissions: []` → Backend converte para `manage + all` ✅
2. Menu:Participantes está definido ✅
3. Estrutura geral do YAML ✅

### ⚠️ O que precisa ajustar:

1. **Subject de participantes:**
   - ❌ Atual: `Participants` (singular, maiúsculo)
   - ✅ Correto: `participantes` (plural, minúsculo - como frontend espera)

### 📋 Checklist:

- [ ] Verificar no código do frontend qual padrão é usado (minúsculo/maiúsculo, singular/plural)
- [ ] Usar `subject` **exatamente** como o frontend verifica no CASL.js
- [ ] Testar endpoint `/me` para confirmar que subjects estão corretos
- [ ] Ajustar YAML conforme padrão identificado

---

**Recomendação Final:** Use `subject: participantes` (minúsculo, plural) para garantir compatibilidade com o frontend que verifica `ability.can('read', 'participantes')`.

