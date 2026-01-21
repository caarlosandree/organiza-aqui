package repository

import (
	"fmt"
	"strings"
)

// QueryBuilder ajuda a construir queries SQL de forma segura
type QueryBuilder struct {
	baseQuery   string
	whereClause []string
	orderBy     string
	limitValue  *int
	offsetValue *int
	args        []interface{}
}

// NewQueryBuilder cria um novo QueryBuilder
func NewQueryBuilder(baseQuery string) *QueryBuilder {
	return &QueryBuilder{
		baseQuery:   baseQuery,
		whereClause: make([]string, 0),
		args:        make([]interface{}, 0),
	}
}

// WhereEqual adiciona uma condição WHERE com igualdade
func (qb *QueryBuilder) WhereEqual(column string, value interface{}) *QueryBuilder {
	if value != nil {
		qb.whereClause = append(qb.whereClause, fmt.Sprintf("%s = $%d", column, len(qb.args)+1))
		qb.args = append(qb.args, value)
	}
	return qb
}

// WhereNotEqual adiciona uma condição WHERE com diferença
func (qb *QueryBuilder) WhereNotEqual(column string, value interface{}) *QueryBuilder {
	if value != nil {
		qb.whereClause = append(qb.whereClause, fmt.Sprintf("%s != $%d", column, len(qb.args)+1))
		qb.args = append(qb.args, value)
	}
	return qb
}

// WhereIn adiciona uma condição WHERE IN
func (qb *QueryBuilder) WhereIn(column string, values []interface{}) *QueryBuilder {
	if len(values) > 0 {
		placeholders := make([]string, len(values))
		startPos := len(qb.args) + 1
		for i, val := range values {
			placeholders[i] = fmt.Sprintf("$%d", startPos+i)
			qb.args = append(qb.args, val)
		}
		qb.whereClause = append(qb.whereClause, fmt.Sprintf("%s IN (%s)", column, strings.Join(placeholders, ", ")))
	}
	return qb
}

// WhereGreaterOrEqual adiciona uma condição WHERE >=
func (qb *QueryBuilder) WhereGreaterOrEqual(column string, value interface{}) *QueryBuilder {
	if value != nil {
		qb.whereClause = append(qb.whereClause, fmt.Sprintf("%s >= $%d", column, len(qb.args)+1))
		qb.args = append(qb.args, value)
	}
	return qb
}

// WhereLessOrEqual adiciona uma condição WHERE <=
func (qb *QueryBuilder) WhereLessOrEqual(column string, value interface{}) *QueryBuilder {
	if value != nil {
		qb.whereClause = append(qb.whereClause, fmt.Sprintf("%s <= $%d", column, len(qb.args)+1))
		qb.args = append(qb.args, value)
	}
	return qb
}

// WhereBetween adiciona uma condição WHERE BETWEEN
func (qb *QueryBuilder) WhereBetween(column string, start, end interface{}) *QueryBuilder {
	if start != nil && end != nil {
		pos1 := len(qb.args) + 1
		pos2 := len(qb.args) + 2
		qb.whereClause = append(qb.whereClause, fmt.Sprintf("%s >= $%d AND %s <= $%d", column, pos1, column, pos2))
		qb.args = append(qb.args, start, end)
	}
	return qb
}

// WhereIsNull adiciona uma condição WHERE IS NULL
func (qb *QueryBuilder) WhereIsNull(column string) *QueryBuilder {
	qb.whereClause = append(qb.whereClause, fmt.Sprintf("%s IS NULL", column))
	return qb
}

// WhereIsNotNull adiciona uma condição WHERE IS NOT NULL
func (qb *QueryBuilder) WhereIsNotNull(column string) *QueryBuilder {
	qb.whereClause = append(qb.whereClause, fmt.Sprintf("%s IS NOT NULL", column))
	return qb
}

// WhereArrayContains adiciona uma condição WHERE para arrays (PostgreSQL = ANY)
func (qb *QueryBuilder) WhereArrayContains(column string, value interface{}) *QueryBuilder {
	if value != nil {
		qb.whereClause = append(qb.whereClause, fmt.Sprintf("$%d = ANY(%s)", len(qb.args)+1, column))
		qb.args = append(qb.args, value)
	}
	return qb
}

// WhereArrayOverlaps adiciona uma condição WHERE para arrays (PostgreSQL &&)
func (qb *QueryBuilder) WhereArrayOverlaps(column string, values interface{}) *QueryBuilder {
	if values != nil {
		qb.whereClause = append(qb.whereClause, fmt.Sprintf("%s && $%d", column, len(qb.args)+1))
		qb.args = append(qb.args, values)
	}
	return qb
}

// WhereCustom adiciona uma condição WHERE customizada
// A condição deve usar placeholders $n onde n começa em len(qb.args)+1
// Exemplo: WhereCustom("column = $1", value) onde $1 será substituído automaticamente
func (qb *QueryBuilder) WhereCustom(condition string, args ...interface{}) *QueryBuilder {
	if len(args) == 0 {
		qb.whereClause = append(qb.whereClause, condition)
		return qb
	}
	
	// Substituir placeholders na condição
	startPos := len(qb.args) + 1
	for i, arg := range args {
		oldPlaceholder := fmt.Sprintf("$%d", i+1)
		newPlaceholder := fmt.Sprintf("$%d", startPos+i)
		condition = strings.ReplaceAll(condition, oldPlaceholder, newPlaceholder)
		qb.args = append(qb.args, arg)
	}
	
	qb.whereClause = append(qb.whereClause, condition)
	return qb
}

// OrderBy adiciona uma cláusula ORDER BY
func (qb *QueryBuilder) OrderBy(orderBy string) *QueryBuilder {
	qb.orderBy = orderBy
	return qb
}

// Limit adiciona uma cláusula LIMIT
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	if limit > 0 {
		qb.limitValue = &limit
	}
	return qb
}

// Offset adiciona uma cláusula OFFSET
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	if offset > 0 {
		qb.offsetValue = &offset
	}
	return qb
}

// Build constrói a query SQL final
func (qb *QueryBuilder) Build() (string, []interface{}) {
	query := qb.baseQuery
	
	// Adicionar WHERE clause
	if len(qb.whereClause) > 0 {
		query += " WHERE " + strings.Join(qb.whereClause, " AND ")
	}
	
	// Adicionar ORDER BY
	if qb.orderBy != "" {
		query += " ORDER BY " + qb.orderBy
	}
	
	// Adicionar LIMIT
	if qb.limitValue != nil {
		pos := len(qb.args) + 1
		query += fmt.Sprintf(" LIMIT $%d", pos)
		qb.args = append(qb.args, *qb.limitValue)
	}
	
	// Adicionar OFFSET
	if qb.offsetValue != nil {
		pos := len(qb.args) + 1
		query += fmt.Sprintf(" OFFSET $%d", pos)
		qb.args = append(qb.args, *qb.offsetValue)
	}
	
	return query, qb.args
}

// Args retorna os argumentos da query
func (qb *QueryBuilder) Args() []interface{} {
	return qb.args
}

// toStringSlice converte []string para []interface{}
func toStringSlice(slice []string) []interface{} {
	result := make([]interface{}, len(slice))
	for i, s := range slice {
		result[i] = s
	}
	return result
}
