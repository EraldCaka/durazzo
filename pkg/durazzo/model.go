package durazzo

import (
	"fmt"
	"reflect"
	"strings"
)

// AutoMigrate creates tables based on the struct configuration for multiple models that might be inserted
func (d *Durazzo) AutoMigrate(models ...interface{}) error {
	for _, model := range models {
		modelType := reflect.TypeOf(model)
		if modelType.Kind() != reflect.Ptr || modelType.Elem().Kind() != reflect.Struct {
			return fmt.Errorf("model %v must be a pointer to a struct", modelType.Name())
		}

		tableName := strings.ToLower(modelType.Elem().Name())
		var columns []string

		for i := 0; i < modelType.Elem().NumField(); i++ {
			field := modelType.Elem().Field(i)
			column := strings.ToLower(field.Name)
			tag := field.Tag.Get("durazzo")
			sqlType := determineSQLType(field.Type, tag)

			if strings.Contains(tag, "primary_key") {
				sqlType = "SERIAL PRIMARY KEY"
			}

			if strings.Contains(tag, "unique") {
				sqlType += " UNIQUE"
			}

			columns = append(columns, fmt.Sprintf(`"%s" %s`, column, sqlType))
		}

		createQuery := fmt.Sprintf(
			`CREATE TABLE IF NOT EXISTS "%s" (%s);`,
			tableName,
			strings.Join(columns, ", "),
		)

		_, err := d.db.Exec(createQuery)
		if err != nil {
			return fmt.Errorf("failed to create table for model %v: %w", tableName, err)
		}
	}
	return nil
}

// Determine the SQL type for a field based on its Go type and struct tag
func determineSQLType(goType reflect.Type, tag string) string {
	switch goType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "INTEGER"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "BIGINT"
	case reflect.Float32, reflect.Float64:
		return "REAL"
	case reflect.String:
		if strings.Contains(tag, "size") {
			return fmt.Sprintf("VARCHAR(%s)", extractSize(tag))
		}
		return "TEXT"
	case reflect.Bool:
		return "BOOLEAN"
	default:
		return "TEXT"
	}
}

func extractSize(tag string) string {
	parts := strings.Split(tag, ":")
	if len(parts) > 1 {
		return parts[1]
	}
	return "255"
}
