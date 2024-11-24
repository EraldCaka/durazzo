package durazzo

import (
	"fmt"
	"strings"
)

// DeleteType handles DELETE operations
type DeleteType struct {
	*Durazzo
	tableName  string
	conditions []string
	args       []interface{}
}

// Delete initializes a DELETE operation
func (d *Durazzo) Delete(tableName string) *DeleteType {
	return &DeleteType{
		Durazzo:    d,
		tableName:  tableName,
		conditions: []string{},
		args:       []interface{}{},
	}
}

// Where adds a condition to the DELETE query
func (dt *DeleteType) Where(field, value string) *DeleteType {
	dt.conditions = append(dt.conditions, fmt.Sprintf(`%s = $%d`, field, len(dt.args)+1))
	dt.args = append(dt.args, value)
	return dt
}

// Run executes the DELETE query
func (dt *DeleteType) Run() error {
	if len(dt.conditions) == 0 {
		return fmt.Errorf("no conditions specified for DELETE operation")
	}

	query := fmt.Sprintf(`DELETE FROM "%s" WHERE %s`, dt.tableName, strings.Join(dt.conditions, " AND "))
	_, err := dt.Durazzo.Db.Exec(query, dt.args...)
	return err
}
