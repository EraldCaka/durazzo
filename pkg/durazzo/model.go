package durazzo

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// AutoMigrate creates a table based on the struct configuration
func (d *Durazzo) AutoMigrate(model interface{}) error {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() != reflect.Ptr || modelType.Elem().Kind() != reflect.Struct {
		return errors.New("model must be a pointer to a struct")
	}

	tableName := strings.ToLower(modelType.Elem().Name())
	var columns []string

	for i := 0; i < modelType.Elem().NumField(); i++ {
		field := modelType.Elem().Field(i)
		column := strings.ToLower(field.Name)
		tag := field.Tag.Get("durazzo")
		sqlType := "TEXT"

		if tag == "primary_key" {
			sqlType = "SERIAL PRIMARY KEY"
		} else if strings.Contains(tag, "size") {
			sqlType = fmt.Sprintf("VARCHAR(%s)", extractSize(tag))
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
	return err
}

// Extract size for VARCHAR type from struct tag
func extractSize(tag string) string {
	parts := strings.Split(tag, ":")
	if len(parts) > 1 {
		return parts[1]
	}
	return "255"
}
