package database

import (
	"strings"

	"gorm.io/gorm/schema"
)

// NamingStrategy implements custom naming strategy for GORM
type NamingStrategy struct {
	schema.NamingStrategy
}

// TableName converts struct name to table name
func (ns NamingStrategy) TableName(str string) string {
	// Convert CamelCase to snake_case and add plural
	return ns.toSnakeCase(str) + "s"
}

// ColumnName converts field name to column name
func (ns NamingStrategy) ColumnName(table, column string) string {
	return ns.toSnakeCase(column)
}

// JoinTableName converts relation to join table name
func (ns NamingStrategy) JoinTableName(joinTable string) string {
	return ns.toSnakeCase(joinTable)
}

// RelationshipFKName generates foreign key name for relation
func (ns NamingStrategy) RelationshipFKName(rel schema.Relationship) string {
	return ns.toSnakeCase(rel.Name) + "_id"
}

// CheckerName generates checker name
func (ns NamingStrategy) CheckerName(table, column string) string {
	return table + "_" + column + "_check"
}

// IndexName generates index name
func (ns NamingStrategy) IndexName(table, column string) string {
	return "idx_" + table + "_" + column
}

// toSnakeCase converts CamelCase to snake_case
func (ns NamingStrategy) toSnakeCase(str string) string {
	var result strings.Builder
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r - 'A' + 'a')
	}
	return result.String()
}