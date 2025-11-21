# Melhorias de UX no CLI - Inferência Inteligente

**Data:** 2025-01-XX  
**Status:** ✅ Implementado

---

## 🎯 Objetivo

Tornar o CLI mais intuitivo e menos verboso, permitindo que o desenvolvedor informe apenas informações básicas enquanto o CLI infere automaticamente os detalhes técnicos.

---

## ✅ Mudanças Implementadas

### 1. **Application - Inferência Automática**

**Antes:**
```
Código da aplicação (slug, ex: sagep-biopass): [usuário digita]
Nome da aplicação: [usuário digita]
```

**Agora:**
```
Nome da aplicação (ex: Biopass, CRV, Core): [usuário digita apenas "Biopass"]
```

**O que o CLI faz:**
- ✅ Infere `code: sagep-biopass` automaticamente
- ✅ Infere `name: SAGEP Biopass` automaticamente
- ✅ Mostra preview e permite confirmar/editar

**Funções criadas:**
- `InferApplicationCode()`: "Biopass" → "sagep-biopass"
- `InferApplicationName()`: "Biopass" → "SAGEP Biopass"

---

### 2. **Permissions - UX Simplificada**

**Antes:**
```
Código da permissão: [usuário digita "biopass.devices.read"]
[CLI tenta inferir...]
Subject: [usuário confirma/edita]
Action: [usuário confirma/edita]
```

**Agora:**

#### Para Menus:
```
Tipo de permissão: [Menu | Recurso (entidade)]
Nome do menu: [usuário digita apenas "Dashboard"]
```

**O que o CLI faz:**
- ✅ Cria automaticamente: `code: Menu:Dashboard`
- ✅ Cria automaticamente: `subject: Menu:Dashboard`
- ✅ Cria automaticamente: `action: view`

#### Para Recursos:
```
Tipo de permissão: [Menu | Recurso (entidade)]
Nome da entidade: [usuário digita "participantes"]
Operação permitida: [read | create | update | delete | manage | view]
```

**O que o CLI faz:**
- ✅ Cria automaticamente: `code: biopass.participants.read`
- ✅ Cria automaticamente: `subject: participantes` (minúsculo, plural - como frontend espera)
- ✅ Cria automaticamente: `action: read`

**Funções criadas:**
- `InferMenuPermission()`: "Dashboard" → `code: Menu:Dashboard`, `subject: Menu:Dashboard`, `action: view`
- `InferResourcePermission()`: entidade="participantes", action="read", appCode="sagep-biopass" → `code: biopass.participants.read`, `subject: participantes`, `action: read`

---

## 📋 Fluxo de Uso Atual

### 1. Criar Application

```
🚀 Criando novo manifest...

Nome da aplicação (ex: Biopass, CRV, Core): Biopass
Descrição (opcional): Sistema de biometria

   ✅ Informações inferidas:
      Código: sagep-biopass
      Nome:   SAGEP Biopass

Confirmar informações da aplicação? (Y/n): y
```

### 2. Criar Permissions

#### Menu:
```
Tipo de permissão:
  > Menu
  Recurso (entidade)

Nome do menu (ex: Dashboard, Participantes): Dashboard

   ✅ Permissão de menu criada:
      Code:    Menu:Dashboard
      Subject: Menu:Dashboard
      Action:  view

Confirmar permissão criada? (Y/n): y
```

#### Recurso:
```
Tipo de permissão:
  Menu
  > Recurso (entidade)

Nome da entidade (ex: participantes, devices, users): participantes
Operação permitida:
  > read
  create
  update
  delete
  manage
  view

   ✅ Permissão de recurso criada:
      Code:    biopass.participants.read
      Subject: participantes
      Action:  read

Confirmar permissão criada? (Y/n): y
```

---

## 🎯 Benefícios

### Para o Desenvolvedor:

1. **Menos digitação:** Informa apenas o essencial
2. **Menos erros:** CLI garante padrões corretos
3. **Mais rápido:** Fluxo mais direto
4. **Mais intuitivo:** Perguntas em linguagem natural

### Para o Sistema:

1. **Padrões consistentes:** CLI sempre gera no formato correto
2. **Compatibilidade CASL.js:** Subjects sempre no formato que frontend espera
3. **Menos ambiguidade:** Inferência clara e previsível

---

## 🔧 Detalhes Técnicos

### Funções de Inferência

#### `InferApplicationCode(appName string) string`
- Entrada: `"Biopass"`
- Saída: `"sagep-biopass"`
- Lógica: Converte para minúsculo, adiciona prefixo "sagep-" se não tiver

#### `InferApplicationName(appName string) string`
- Entrada: `"Biopass"`
- Saída: `"SAGEP Biopass"`
- Lógica: Capitaliza primeira letra, adiciona prefixo "SAGEP " se não tiver

#### `InferMenuPermission(menuName string) (code, subject, action string)`
- Entrada: `"Dashboard"` ou `"dashboard"`
- Saída: `code="Menu:Dashboard"`, `subject="Menu:Dashboard"`, `action="view"`
- Lógica: Capitaliza primeira letra, adiciona prefixo "Menu:", action sempre "view"

#### `InferResourcePermission(entidade, action, appCode string) (code, subject, action string)`
- Entrada: `entidade="participantes"`, `action="read"`, `appCode="sagep-biopass"`
- Saída: `code="biopass.participants.read"`, `subject="participantes"`, `action="read"`
- Lógica:
  - Extrai código curto da app: "sagep-biopass" → "biopass"
  - Gera code: `{appShort}.{entidade}.{action}`
  - Subject mantém minúsculo/plural (como frontend espera)
  - Valida action é válido CASL.js

---

## 📝 Exemplos de Uso

### Exemplo Completo: Criar Manifest

```
🚀 Criando novo manifest...

Nome da aplicação: Biopass
Descrição: Sistema de biometria facial

   ✅ Informações inferidas:
      Código: sagep-biopass
      Nome:   SAGEP Biopass

Confirmar? (Y/n): y

Deseja criar permissões agora? (Y/n): y

Tipo de permissão: Menu
Nome do menu: Dashboard
   ✅ Criada: Menu:Dashboard

Tipo de permissão: Recurso
Nome da entidade: participantes
Operação: read
   ✅ Criada: biopass.participants.read

Adicionar outra permissão? (Y/n): n
```

**Resultado no YAML:**
```yaml
application:
  code: sagep-biopass
  name: SAGEP Biopass
  description: Sistema de biometria facial

permissions:
  - code: Menu:Dashboard
    subject: Menu:Dashboard
    action: view
  
  - code: biopass.participants.read
    subject: participantes
    action: read
```

---

## ✅ Checklist de Implementação

- [x] Função `InferApplicationCode()` criada
- [x] Função `InferApplicationName()` criada
- [x] Função `InferMenuPermission()` criada
- [x] Função `InferResourcePermission()` criada
- [x] Fluxo de perguntas para Application simplificado
- [x] Fluxo de perguntas para Permissions simplificado
- [x] Seleção de tipo (Menu vs Recurso) implementada
- [x] Preview das informações inferidas
- [x] Opção de confirmar/editar mantida
- [x] Código compilando sem erros

---

## 🎯 Próximos Passos (Opcional)

1. **Validação de entidades comuns:** Sugerir entidades já usadas
2. **Templates:** Permitir criar múltiplas permissions de uma vez (CRUD completo)
3. **Importar de YAML existente:** Ler manifest existente e sugerir melhorias

---

**Status:** ✅ Implementado e funcionando  
**Compatibilidade:** ✅ Mantida com versões anteriores

