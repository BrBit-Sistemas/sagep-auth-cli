package manifest

import (
	"strings"
)

// ValidActions são as ações válidas do CASL.js
var ValidActions = []string{"read", "create", "update", "delete", "manage", "view"}

// InferMenuPermission cria uma permission de menu a partir do nome do menu
// Entrada: "Dashboard" ou "dashboard" → Saída: code="Menu:Dashboard", subject="Menu:Dashboard", action="view"
func InferMenuPermission(menuName string) (code, subject, action string) {
	// Normalizar: capitalizar primeira letra
	menuName = strings.TrimSpace(menuName)
	if len(menuName) == 0 {
		return "", "", ""
	}
	
	// Capitalizar primeira letra
	normalized := strings.ToUpper(menuName[:1]) + strings.ToLower(menuName[1:])
	
	code = "Menu:" + normalized
	subject = code
	action = "view"
	
	return code, subject, action
}

// InferResourcePermission cria uma permission de recurso a partir de entidade e ação
// Entrada: entidade="participantes", action="read", appCode="sagep-biopass"
// Saída: code="biopass.participants.read", subject="participantes", action="read"
func InferResourcePermission(entidade, action, appCode string) (code, subject, actionOut string) {
	entidade = strings.TrimSpace(strings.ToLower(entidade))
	action = strings.TrimSpace(strings.ToLower(action))
	appCode = strings.TrimSpace(strings.ToLower(appCode))
	
	if entidade == "" || action == "" || appCode == "" {
		return "", "", ""
	}
	
	// Validar action
	valid := false
	for _, validAction := range ValidActions {
		if action == validAction {
			valid = true
			break
		}
	}
	if !valid {
		return "", "", ""
	}
	
	// Extrair código curto da aplicação (ex: "sagep-biopass" → "biopass")
	appShort := extractAppShortCode(appCode)
	
	// Gerar code: {appShort}.{entidade}.{action}
	code = appShort + "." + entidade + "." + action
	
	// Subject é a entidade no formato que o frontend espera (minúsculo, plural)
	subject = entidade
	
	// Action é a ação
	actionOut = action
	
	return code, subject, actionOut
}

// extractAppShortCode extrai código curto da aplicação
// Ex: "sagep-biopass" → "biopass"
// Ex: "sagep-crv" → "crv"
// Ex: "biopass" → "biopass"
func extractAppShortCode(appCode string) string {
	parts := strings.Split(appCode, "-")
	if len(parts) > 1 {
		// Se tem hífen, pegar última parte
		return parts[len(parts)-1]
	}
	// Se não tem hífen, usar como está
	return appCode
}

// InferApplicationCode gera código da aplicação a partir do nome
// Entrada: "Biopass" → Saída: "sagep-biopass"
func InferApplicationCode(appName string) string {
	appName = strings.TrimSpace(appName)
	if len(appName) == 0 {
		return ""
	}
	
	// Converter para minúsculo e substituir espaços por hífens
	slug := strings.ToLower(appName)
	slug = strings.ReplaceAll(slug, " ", "-")
	
	// Se não começa com "sagep-", adicionar
	if !strings.HasPrefix(slug, "sagep-") {
		slug = "sagep-" + slug
	}
	
	return slug
}

// InferApplicationName normaliza o nome da aplicação
// Entrada: "biopass" ou "Biopass" → Saída: "SAGEP Biopass"
func InferApplicationName(appName string) string {
	appName = strings.TrimSpace(appName)
	if len(appName) == 0 {
		return ""
	}
	
	// Capitalizar primeira letra
	normalized := strings.ToUpper(appName[:1]) + strings.ToLower(appName[1:])
	
	// Se não começa com "SAGEP", adicionar
	if !strings.HasPrefix(strings.ToUpper(normalized), "SAGEP") {
		normalized = "SAGEP " + normalized
	}
	
	return normalized
}

// InferSubjectAndAction tenta inferir subject e action a partir do code
// Retorna subject, action e um booleano indicando se a inferência foi bem-sucedida
// MANTIDO PARA COMPATIBILIDADE - mas agora temos funções mais específicas acima
func InferSubjectAndAction(code string) (subject string, action string, ok bool) {
	code = strings.TrimSpace(code)
	if code == "" {
		return "", "", false
	}

	// 1. Padrão Menu:{Nome} → subject="Menu:{Nome}", action="view"
	if strings.HasPrefix(code, "Menu:") {
		return code, "view", true
	}

	// 2. Padrão simples: {Subject}.{Action} (ex: "Device.read")
	parts := strings.SplitN(code, ".", 2)
	if len(parts) == 2 {
		possibleSubject := parts[0]
		possibleAction := parts[1]

		// Verificar se action é válido CASL (não contém mais pontos)
		for _, validAction := range ValidActions {
			if possibleAction == validAction {
				// Capitalizar primeira letra do subject
				subject = capitalizeFirst(possibleSubject)
				return subject, possibleAction, true
			}
		}
	}

	// 3. Padrão com múltiplos pontos: {app}.{resource}.{action} (ex: "biopass.devices.read")
	// Tentar extrair da última parte
	lastDotIndex := strings.LastIndex(code, ".")
	if lastDotIndex > 0 {
		possibleAction := code[lastDotIndex+1:]
		
		// Verificar se a última parte é uma action válida
		for _, validAction := range ValidActions {
			if possibleAction == validAction {
				// Extrair subject (tudo antes do último ponto)
				possibleSubject := code[:lastDotIndex]
				
				// Tentar extrair nome do recurso (última parte antes da action)
				subjectParts := strings.Split(possibleSubject, ".")
				if len(subjectParts) > 0 {
					// Pegar a última parte e manter minúsculo (frontend espera minúsculo)
					resourceName := subjectParts[len(subjectParts)-1]
					// NÃO capitalizar - frontend espera minúsculo
					return resourceName, possibleAction, true
				}
			}
		}
	}

	// 4. Se não conseguiu inferir, retornar false
	return "", "", false
}

// capitalizeFirst capitaliza a primeira letra de uma string
func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}
