#!/bin/bash
# Script para configurar o ambiente Go para permitir acesso público ao sagep-auth-cli
# IMPORTANTE: Apenas github.com/BrBit-Sistemas/sagep-auth-cli é público.
# Todos os outros repositórios em github.com/BrBit-Sistemas são privados.

set -e

echo "🔧 Configurando ambiente Go para acesso público ao sagep-auth-cli..."
echo "   (outros repositórios BrBit-Sistemas permanecem privados)"

# Verificar configuração atual
CURRENT_GOPRIVATE=$(go env GOPRIVATE)
CURRENT_GONOPROXY=$(go env GONOPROXY)
CURRENT_GONOSUMDB=$(go env GONOSUMDB)

echo ""
echo "Configuração atual:"
echo "  GOPRIVATE: ${CURRENT_GOPRIVATE:-'(vazio)'}"
echo "  GONOPROXY:  ${CURRENT_GONOPROXY:-'(vazio)'}"
echo "  GONOSUMDB:  ${CURRENT_GONOSUMDB:-'(vazio)'}"
echo ""

# Estratégia:
# 1. Garantir que GOPRIVATE contém github.com/BrBit-Sistemas (para proteger outros repositórios privados)
# 2. Adicionar exceção via GONOPROXY/GONOSUMDB para permitir acesso público apenas ao sagep-auth-cli

REPO_PUBLIC="github.com/BrBit-Sistemas/sagep-auth-cli"
ORG_PRIVATE="github.com/BrBit-Sistemas"

# Passo 1: Garantir que GOPRIVATE contém github.com/BrBit-Sistemas
if [ -z "$CURRENT_GOPRIVATE" ]; then
    echo "📝 Adicionando ${ORG_PRIVATE} ao GOPRIVATE (protege outros repositórios privados)..."
    go env -w GOPRIVATE="$ORG_PRIVATE"
    CURRENT_GOPRIVATE="$ORG_PRIVATE"
elif ! echo "$CURRENT_GOPRIVATE" | tr ',' '\n' | grep -q "^${ORG_PRIVATE}$"; then
    # Adicionar se não estiver presente
    echo "📝 Adicionando ${ORG_PRIVATE} ao GOPRIVATE (protege outros repositórios privados)..."
    NEW_GOPRIVATE="${CURRENT_GOPRIVATE},${ORG_PRIVATE}"
    go env -w GOPRIVATE="$NEW_GOPRIVATE"
    CURRENT_GOPRIVATE="$NEW_GOPRIVATE"
else
    echo "✓ ${ORG_PRIVATE} já está no GOPRIVATE"
fi

# Passo 2: Adicionar exceção para o repositório público via GONOPROXY e GONOSUMDB
echo ""
echo "📝 Configurando exceções para ${REPO_PUBLIC}..."

# Verificar se o repositório público já está nas exceções
HAS_IN_GONOPROXY=false
HAS_IN_GONOSUMDB=false

if [ -n "$CURRENT_GONOPROXY" ]; then
    echo "$CURRENT_GONOPROXY" | tr ',' '\n' | grep -q "^${REPO_PUBLIC}$" && HAS_IN_GONOPROXY=true
fi

if [ -n "$CURRENT_GONOSUMDB" ]; then
    echo "$CURRENT_GONOSUMDB" | tr ',' '\n' | grep -q "^${REPO_PUBLIC}$" && HAS_IN_GONOSUMDB=true
fi

# Adicionar o repositório público às exceções se não estiver
if [ "$HAS_IN_GONOPROXY" = false ]; then
    if [ -z "$CURRENT_GONOPROXY" ]; then
        NEW_GONOPROXY="$REPO_PUBLIC"
    else
        NEW_GONOPROXY="${CURRENT_GONOPROXY},${REPO_PUBLIC}"
    fi
    echo "📝 Adicionando $REPO_PUBLIC ao GONOPROXY (permite acesso público)..."
    go env -w GONOPROXY="$NEW_GONOPROXY"
else
    echo "✓ $REPO_PUBLIC já está no GONOPROXY"
fi

if [ "$HAS_IN_GONOSUMDB" = false ]; then
    if [ -z "$CURRENT_GONOSUMDB" ]; then
        NEW_GONOSUMDB="$REPO_PUBLIC"
    else
        NEW_GONOSUMDB="${CURRENT_GONOSUMDB},${REPO_PUBLIC}"
    fi
    echo "📝 Adicionando $REPO_PUBLIC ao GONOSUMDB (permite checksum público)..."
    go env -w GONOSUMDB="$NEW_GONOSUMDB"
else
    echo "✓ $REPO_PUBLIC já está no GONOSUMDB"
fi

# Remover apenas github.com/brbit (antigo) se existir, mas manter github.com/BrBit-Sistemas
# pois outros repositórios são privados
if [ -n "$CURRENT_GOPRIVATE" ]; then
    # Verificar se contém github.com/brbit (antigo)
    if echo "$CURRENT_GOPRIVATE" | tr ',' '\n' | grep -q '^github.com/brbit'; then
        NEW_GOPRIVATE=$(echo "$CURRENT_GOPRIVATE" | tr ',' '\n' | \
            grep -v '^github.com/brbit$' | \
            grep -v '^github.com/brbit/' | \
            tr '\n' ',' | sed 's/,$//' | sed 's/^,//')
        
        if [ -n "$NEW_GOPRIVATE" ] && [ "$NEW_GOPRIVATE" != "$CURRENT_GOPRIVATE" ]; then
            echo "📝 Removendo github.com/brbit (antigo) do GOPRIVATE..."
            go env -w GOPRIVATE="$NEW_GOPRIVATE"
        fi
    else
        echo "✓ Nenhuma configuração antiga (github.com/brbit) encontrada no GOPRIVATE"
    fi
fi

echo ""
echo "✅ Configuração final:"
FINAL_GOPRIVATE=$(go env GOPRIVATE)
FINAL_GONOPROXY=$(go env GONOPROXY)
FINAL_GONOSUMDB=$(go env GONOSUMDB)

echo "  GOPRIVATE: ${FINAL_GOPRIVATE:-'(vazio)'}"
if echo "$FINAL_GOPRIVATE" | grep -q "github.com/BrBit-Sistemas"; then
    echo "    → github.com/BrBit-Sistemas configurado (outros repositórios são privados) ✓"
else
    echo "    ⚠️  ATENÇÃO: github.com/BrBit-Sistemas não está no GOPRIVATE!"
fi
echo "  GONOPROXY:  ${FINAL_GONOPROXY:-'(vazio)'}"
if echo "$FINAL_GONOPROXY" | grep -q "sagep-auth-cli"; then
    echo "    → sagep-auth-cli com exceção para acesso público ✓"
else
    echo "    ⚠️  ATENÇÃO: sagep-auth-cli não está no GONOPROXY!"
fi
echo "  GONOSUMDB:  ${FINAL_GONOSUMDB:-'(vazio)'}"
if echo "$FINAL_GONOSUMDB" | grep -q "sagep-auth-cli"; then
    echo "    → sagep-auth-cli com exceção para checksum público ✓"
else
    echo "    ⚠️  ATENÇÃO: sagep-auth-cli não está no GONOSUMDB!"
fi
echo ""
echo "✓ Ambiente Go configurado!"
echo "  → sagep-auth-cli: acesso público via exceção (GONOPROXY/GONOSUMDB)"
echo "  → Outros repositórios BrBit-Sistemas: protegidos via GOPRIVATE"

