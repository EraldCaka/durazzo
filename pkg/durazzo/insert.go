package durazzo

import (
	"errors"
	"fmt"
	"github.com/EraldCaka/durazzo/pkg/util"
	"log"
	"reflect"
	"strings"
)

// InsertType handles INSERT operations
type InsertType struct {
	*Durazzo
	tableName string
	model     interface{}
}

// Insert initializes an INSERT operation
func (d *Durazzo) Insert(model interface{}) *InsertType {
	_, tableName, _, err := util.ResolveModelInfo(model)
	if err != nil {
		log.Fatalf("failed to initialize InsertType: %v", err)
	}

	return &InsertType{
		Durazzo:   d,
		tableName: tableName,
		model:     model,
	}
}

// Run executes the INSERT query
func (it *InsertType) Run() error {
	columns, values, placeholders, err := prepareInsertData(it.model)
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s)`, it.tableName, strings.Join(columns, ", "), strings.Join(placeholders, ", "))
	_, err = it.Durazzo.Db.Exec(query, values...)
	return err
}

// prepareInsertData prepares the columns, values, and placeholders for an INSERT statement
func prepareInsertData(model interface{}) ([]string, []interface{}, []string, error) {
	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}

	if modelValue.Kind() != reflect.Struct {
		return nil, nil, nil, errors.New("model must be a struct or a pointer to a struct")
	}

	var columns []string
	var values []interface{}
	var placeholders []string

	for i := 0; i < modelValue.NumField(); i++ {
		field := modelValue.Field(i)
		if !field.CanInterface() {
			continue
		}
		columns = append(columns, modelValue.Type().Field(i).Name)
		values = append(values, field.Interface())
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(placeholders)+1))
	}

	return columns, values, placeholders, nil
}
