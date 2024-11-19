package durazzo

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"
)

// SelectType handles SELECT queries and is created by Durazzo
type SelectType struct {
	*Durazzo
	modelType    reflect.Type
	tableName    string
	model        interface{}
	conditions   []string
	args         []interface{}
	limit        int
	isPointer    bool
	queryBuilder QueryBuilder
}

// QueryBuilder defines methods to construct SQL queries
type QueryBuilder interface {
	BuildSelectQuery(tableName string, conditions []string, limit int) (string, error)
}

type SQLQueryBuilder struct{}

func (qb *SQLQueryBuilder) BuildSelectQuery(tableName string, conditions []string, limit int) (string, error) {
	if tableName == "" {
		return "", errors.New("table name cannot be empty")
	}

	var queryBuilder strings.Builder
	queryBuilder.WriteString(fmt.Sprintf(`SELECT * FROM "%s"`, tableName))

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE " + strings.Join(conditions, " AND "))
	}

	if limit > 0 {
		queryBuilder.WriteString(fmt.Sprintf(" LIMIT %d", limit))
	}

	return queryBuilder.String(), nil
}

// Select initializes a SELECT operation from Durazzo it receives a pointer of an interface
// MUST be a pointer
func (d *Durazzo) Select(model interface{}) *SelectType {
	modelType, tableName, isPointer, err := resolveModelInfo(model)

	if err != nil {
		log.Fatalf("failed to initialize SelectType: %v", err)
	}

	return &SelectType{
		Durazzo:      d,
		modelType:    modelType,
		tableName:    tableName,
		model:        model,
		conditions:   []string{},
		args:         []interface{}{},
		limit:        0,
		isPointer:    isPointer,
		queryBuilder: &SQLQueryBuilder{},
	}
}

// Where adds a condition to the query (e.g., "fieldname = ?")
func (st *SelectType) Where(field, value string) *SelectType {
	st.conditions = append(st.conditions, fmt.Sprintf(`%s = $%d`, field, len(st.args)+1))
	st.args = append(st.args, value)
	return st
}

// Limit sets the limit for the query
func (st *SelectType) Limit(limit int) *SelectType {
	st.limit = limit
	return st
}

func (st *SelectType) Run() error {
	startTime := time.Now()

	query, err := st.queryBuilder.BuildSelectQuery(st.tableName, st.conditions, st.limit)
	if err != nil {
		return err
	}

	rows, err := st.Durazzo.Db.Query(query, st.args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	err = mapRowsToModel(rows, st.model, st.modelType, st.isPointer)

	elapsedTime := time.Since(startTime)
	log.Printf("Query : %s took %v to run\n", query, elapsedTime)

	return err
}

// mapRowsToModel maps database rows to the provided model
func mapRowsToModel(rows *sql.Rows, model interface{}, modelType reflect.Type, isPointer bool) error {
	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() != reflect.Ptr && modelValue.Kind() != reflect.Slice {
		return errors.New("model must be a pointer to a struct or slice")
	}

	targetValue := modelValue.Elem()

	if targetValue.Kind() == reflect.Struct || (targetValue.Kind() == reflect.Ptr && targetValue.Elem().Kind() == reflect.Struct) || (targetValue.Kind() == reflect.Ptr && targetValue.Elem().Kind() == reflect.Invalid) {
		if !rows.Next() {
			log.Println(errors.New("no rows found"))
			return nil
		}
		if targetValue.Kind() == reflect.Ptr {

			if targetValue.IsNil() {

				targetValue.Set(reflect.New(targetValue.Type().Elem()))
			}
			targetValue = targetValue.Elem()
		}

		if targetValue.Kind() != reflect.Struct {
			return fmt.Errorf("targetValue must be a struct or a pointer to a struct, got %s", targetValue.Kind())
		}

		return scanRow(rows, targetValue)
	}

	if targetValue.Kind() == reflect.Slice {
		for rows.Next() {

			elem := reflect.New(modelType)
			if err := scanRow(rows, elem.Elem()); err != nil {
				return err
			}
			if isPointer {
				targetValue.Set(reflect.Append(targetValue, elem))
			} else {
				targetValue.Set(reflect.Append(targetValue, elem.Elem()))
			}
		}
		return nil
	}

	return fmt.Errorf("unsupported model type: %s", targetValue.Kind())
}

// scanRow scans a single row into a struct or a pointer to a struct
func scanRow(rows *sql.Rows, targetValue reflect.Value) error {

	if targetValue.Kind() == reflect.Ptr {
		if targetValue.IsNil() {

			targetValue.Set(reflect.New(targetValue.Type().Elem()))
		}
		targetValue = targetValue.Elem()
	}

	if targetValue.Kind() != reflect.Struct {
		return fmt.Errorf("targetValue must be a struct or a pointer to a struct, got %s", targetValue.Kind())
	}

	fieldPointers := make([]interface{}, targetValue.NumField())
	for i := 0; i < targetValue.NumField(); i++ {
		fieldPointers[i] = targetValue.Field(i).Addr().Interface()
	}

	return rows.Scan(fieldPointers...)
}

// resolveModelInfo extracts model information for table and type resolution
func resolveModelInfo(model interface{}) (reflect.Type, string, bool, error) {
	modelType := reflect.TypeOf(model)
	var tableName string
	var isPointer bool

	switch {
	case modelType.Kind() == reflect.Ptr:
		if modelType.Elem().Kind() == reflect.Struct {

			tableName = strings.ToLower(modelType.Elem().Name())
			modelType = modelType.Elem()
		} else if modelType.Elem().Elem().Name() == "" {
			tableName = strings.ToLower(modelType.Elem().Elem().Elem().Name())
			modelType = modelType.Elem().Elem().Elem()
			isPointer = true
		} else {
			tableName = strings.ToLower(modelType.Elem().Elem().Name())
			modelType = modelType.Elem().Elem()
		}
	case modelType.Kind() == reflect.Slice:
		if modelType.Elem().Kind() == reflect.Ptr {
			tableName = strings.ToLower(modelType.Elem().Elem().Name())
			modelType = modelType.Elem().Elem()
			isPointer = true
		} else {
			tableName = strings.ToLower(modelType.Elem().Name())
			modelType = modelType.Elem()
		}
	default:
		return nil, "", false, fmt.Errorf("unsupported model type: %s", modelType.Kind())
	}

	return modelType, tableName, isPointer, nil
}
