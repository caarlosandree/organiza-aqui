package service

import (
	"strings"
)

// Lexorank é uma implementação simples do algoritmo Lexorank para ordenação de itens
// Baseado em: https://github.com/atlassian/lexorank

const (
	// LexorankMin é o valor mínimo de um rank
	LexorankMin = "0|"
	// LexorankMax é o valor máximo de um rank
	LexorankMax = "z|"
	// LexorankMid é o valor médio de um rank
	LexorankMid = "m|"
)

// GenerateLexorank gera um novo lexorank entre dois ranks
// Se afterRank é vazio, gera um rank no início
// Se beforeRank é vazio, gera um rank no final
func GenerateLexorank(afterRank, beforeRank string) string {
	if afterRank == "" {
		afterRank = LexorankMin
	}
	if beforeRank == "" {
		beforeRank = LexorankMax
	}

	// Se os ranks são iguais, precisamos recalcular
	if afterRank == beforeRank {
		return GenerateLexorankBetween(afterRank, beforeRank)
	}

	return GenerateLexorankBetween(afterRank, beforeRank)
}

// GenerateLexorankBetween gera um rank entre dois ranks
func GenerateLexorankBetween(after, before string) string {
	// Normalizar os ranks (remover o separador |)
	after = strings.TrimSuffix(after, "|")
	before = strings.TrimSuffix(before, "|")

	// Encontrar o primeiro caractere diferente
	i := 0
	for i < len(after) && i < len(before) && after[i] == before[i] {
		i++
	}

	// Se after é prefixo de before, adicionar um caractere
	if i == len(after) {
		if i < len(before) {
			// Pegar o próximo caractere de before e dividir
			nextChar := before[i]
			if nextChar > '0' {
				// Inserir um caractere entre after e before
				return after + string(nextChar-1) + "|"
			}
		}
		// Adicionar 'a' no final
		return after + "a|"
	}

	// Se before é prefixo de after, adicionar um caractere
	if i == len(before) {
		// Adicionar 'a' no final de after
		return after + "a|"
	}

	// Pegar os caracteres diferentes
	afterChar := after[i]
	beforeChar := before[i]

	// Se a diferença é maior que 1, podemos inserir no meio
	if beforeChar-afterChar > 1 {
		midChar := (afterChar + beforeChar) / 2
		return after[:i] + string(midChar) + "|"
	}

	// Se a diferença é 1, precisamos adicionar mais caracteres
	// Pegar o restante de after
	rest := after[i:]
	if len(rest) > 0 {
		// Tentar adicionar um caractere no meio
		return after[:i] + string(afterChar) + "m|"
	}

	// Último recurso: adicionar 'a' no final
	return after + "a|"
}

// GetNextLexorank retorna o próximo lexorank após um rank dado
func GetNextLexorank(rank string) string {
	rank = strings.TrimSuffix(rank, "|")
	return rank + "a|"
}

// GetPreviousLexorank retorna o lexorank anterior a um rank dado
func GetPreviousLexorank(rank string) string {
	rank = strings.TrimSuffix(rank, "|")
	if len(rank) == 0 {
		return LexorankMin
	}
	// Decrementar o último caractere
	lastChar := rank[len(rank)-1]
	if lastChar > '0' {
		return rank[:len(rank)-1] + string(lastChar-1) + "|"
	}
	// Se não pode decrementar, retornar o mínimo
	return LexorankMin
}
