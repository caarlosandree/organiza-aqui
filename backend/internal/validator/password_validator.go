package validator

import (
	"unicode"
)

// ValidatePasswordStrength valida a força da senha
// Retorna true se a senha atender aos critérios mínimos de força
func ValidatePasswordStrength(password string) bool {
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// Senha deve ter pelo menos 3 dos 4 critérios
	criteriaCount := 0
	if hasUpper {
		criteriaCount++
	}
	if hasLower {
		criteriaCount++
	}
	if hasNumber {
		criteriaCount++
	}
	if hasSpecial {
		criteriaCount++
	}

	return criteriaCount >= 3
}
