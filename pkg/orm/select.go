package orm

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
)

type SelectType struct {
	*Durazzo
	modelType reflect.Type
	tableName string
	model     interface{}
	err       error
}

// Select will take an interface of a table and will select the data from that particular table.
func (d *Durazzo) Select(model interface{}) *SelectType {
	modelType := reflect.TypeOf(model)
	var err error
	var tableName string

	log.Println(modelType.Kind())
	if modelType.Kind() == reflect.Pointer {
		if modelType.Elem().Kind() == reflect.Slice || modelType.Elem().Kind() == reflect.Map {
			err = fmt.Errorf("pointer should be only assigned before a structure not a slice or a map")
			log.Println(err)
		}

		log.Println(modelType)
		tableName = strings.ToLower(modelType.Elem().Name())

		modelType = modelType.Elem()

	} else if modelType.Kind() == reflect.Slice {
		if modelType.Elem().Kind() == reflect.Ptr {
			log.Println(strings.ToLower(modelType.Elem().Elem().Name()))
			tableName = strings.ToLower(modelType.Elem().Elem().Name())
			modelType = modelType.Elem().Elem()
		} else {
			log.Println(modelType.Elem().Name())
			tableName = strings.ToLower(modelType.Elem().Name())
			modelType = modelType.Elem()

		}
	} else {
		tableName = strings.ToLower(modelType.Name())
	}

	return &SelectType{d, modelType, tableName, model, err}
}

// Where adds a condition to the query (e.g., "fieldname = ?")
func (d *SelectType) Where(field, value string) *SelectType {
	d.conditions = append(d.conditions, fmt.Sprintf(`%s = $%d`, field, len(d.args)+1))
	d.args = append(d.args, value)
	return d
}

// Limit sets the limit for the query
func (d *SelectType) Limit(limit int) *SelectType {
	d.limit = limit
	return d
}

// First Returns the first element of the
func (d *SelectType) First() *SelectType {
	// TODO : implement
	return d
}

// BuildQuery dynamically builds the SQL query for SELECT
func (d *SelectType) buildQuery() (*strings.Builder, error) {
	var queryBuilder strings.Builder

	if d.tableName == "" {
		return &strings.Builder{}, errors.New("table name cannot be empty")
	}
	queryBuilder.WriteString(fmt.Sprintf(`SELECT * FROM "%s"`, d.tableName))

	return &queryBuilder, nil
}

// Run Executes the SELECT query and scans the results into the provided model
func (d *SelectType) Run() error {
	if d.err != nil {
		return nil
	}
	query, err := d.buildQuery()
	if err != nil {
		return err
	}

	rows, err := d.db.Query(query.String())
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		elem := reflect.New(d.modelType).Elem()

		fieldPointers := make([]interface{}, elem.NumField())
		for i := range fieldPointers {
			fieldPointers[i] = elem.Field(i).Addr().Interface()
		}

		if err := rows.Scan(fieldPointers...); err != nil {
			return err
		}
		rowData := make(map[string]interface{})
		for i := 0; i < elem.NumField(); i++ {
			fieldName := d.modelType.Field(i).Name
			fieldValue := elem.Field(i).Interface()
			rowData[fieldName] = fieldValue
		}
		log.Printf("Row: %+v\n", rowData)
	}

	return rows.Err()
}
