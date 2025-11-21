package manifest

import (
	"strings"
)

// ValidActions são as ações válidas do CASL.js
var ValidActions = []string{"read", "create", "update", "delete", "manage", "view"}

// InferSubjectAndAction tenta inferir subject e action a partir do code
// Retorna subject, action e um booleano indicando se a inferência foi bem-sucedida
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
					// Pegar a última parte e capitalizar
					resourceName := subjectParts[len(subjectParts)-1]
					subject = capitalizeFirst(resourceName)
					return subject, possibleAction, true
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
	return strings.ToUpper(s[:1]) + s[1:]
}

