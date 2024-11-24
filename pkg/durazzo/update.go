package durazzo

import (
	"fmt"
	"strings"
)

// UpdateType handles UPDATE operations
type UpdateType struct {
	*Durazzo
	tableName  string
	updates    []string
	conditions []string
	args       []interface{}
}

// Update initializes an UPDATE operation
func (d *Durazzo) Update(tableName string) *UpdateType {
	return &UpdateType{
		Durazzo:    d,
		tableName:  tableName,
		updates:    []string{},
		conditions: []string{},
		args:       []interface{}{},
	}
}

// Set adds a field-value pair to be updated
func (ut *UpdateType) Set(field string, value interface{}) *UpdateType {
	ut.updates = append(ut.updates, fmt.Sprintf(`%s = $%d`, field, len(ut.args)+1))
	ut.args = append(ut.args, value)
	return ut
}

// Where adds a condition to the UPDATE query
func (ut *UpdateType) Where(field, value string) *UpdateType {
	ut.conditions = append(ut.conditions, fmt.Sprintf(`%s = $%d`, field, len(ut.args)+1))
	ut.args = append(ut.args, value)
	return ut
}

// Run executes the UPDATE query
func (ut *UpdateType) Run() error {
	if len(ut.updates) == 0 {
		return fmt.Errorf("no updates specified for UPDATE operation")
	}
	if len(ut.conditions) == 0 {
		return fmt.Errorf("no conditions specified for UPDATE operation")
	}

	query := fmt.Sprintf(`UPDATE "%s" SET %s WHERE %s`, ut.tableName, strings.Join(ut.updates, ", "), strings.Join(ut.conditions, " AND "))
	_, err := ut.Durazzo.Db.Exec(query, ut.args...)
	return err
}
