package durazzo

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
)

type SelectType struct {
	*Durazzo
	modelType  reflect.Type
	tableName  string
	model      interface{}
	err        error
	conditions []string
	args       []interface{}
	limit      int
	isPointer  bool
}

// Select will take an interface of a table and will select the data from that particular table.
func (d *Durazzo) Select(model interface{}) *SelectType {
	modelType := reflect.TypeOf(model)
	var err error
	var tableName string
	var isPointer bool
	if modelType.Kind() == reflect.Pointer {

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
	} else if modelType.Kind() == reflect.Slice {
		if modelType.Elem().Kind() == reflect.Ptr {
			tableName = strings.ToLower(modelType.Elem().Elem().Name())
			modelType = modelType.Elem().Elem()
			isPointer = true
		} else {
			tableName = strings.ToLower(modelType.Elem().Name())
			modelType = modelType.Elem()
		}
	} else {
		tableName = strings.ToLower(modelType.Name())
	}

	// Initialize SelectType with the proper fields
	return &SelectType{
		Durazzo:    d,
		modelType:  modelType,
		tableName:  tableName,
		isPointer:  isPointer,
		model:      model,
		err:        err,
		conditions: []string{},
		args:       []interface{}{},
		limit:      0,
	}
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

	// Add conditions if any
	if len(d.conditions) > 0 {
		queryBuilder.WriteString(" WHERE " + strings.Join(d.conditions, " AND "))
	}

	// Add limit if set
	if d.limit > 0 {
		queryBuilder.WriteString(fmt.Sprintf(" LIMIT %d", d.limit))
	}

	return &queryBuilder, nil
}

// Run Executes the SELECT query and scans the results into the provided model
func (d *SelectType) Run() error {
	if d.err != nil {
		return d.err
	}

	query, err := d.buildQuery()
	if err != nil {
		return err
	}
	log.Println("executing query: ", query.String())
	rows, err := d.db.Query(query.String(), d.args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Ensure the provided model is a pointer to a struct or a slice
	modelValue := reflect.ValueOf(d.model)

	// If the model is not a pointer, return an error
	if modelValue.Kind() != reflect.Ptr && modelValue.Kind() != reflect.Slice {
		return fmt.Errorf("model must be a pointer to a struct or a slice")
	}

	// Dereference the model to get the actual value
	targetValue := modelValue.Elem()

	// Handle pointer to a struct
	if targetValue.Kind() == reflect.Ptr && targetValue.Elem().Kind() == reflect.Struct {
		if !rows.Next() {
			return errors.New("no rows found")
		}
		// Initialize the struct if it's nil
		if targetValue.IsNil() {
			targetValue.Set(reflect.New(targetValue.Type().Elem()))
		}

		structValue := targetValue.Elem()

		// Create field pointers for scanning
		fieldPointers := make([]interface{}, structValue.NumField())
		for i := 0; i < structValue.NumField(); i++ {
			fieldPointers[i] = structValue.Field(i).Addr().Interface()
		}

		// Scan data into the struct
		if err := rows.Scan(fieldPointers...); err != nil {
			return err
		}
		return nil
	}

	// Handle single-row case (e.g., User struct)
	if targetValue.Kind() == reflect.Struct {
		if !rows.Next() {
			return errors.New("no rows found")
		}

		fieldPointers := make([]interface{}, targetValue.NumField())
		for i := 0; i < targetValue.NumField(); i++ {
			fieldPointers[i] = targetValue.Field(i).Addr().Interface()
		}

		if err := rows.Scan(fieldPointers...); err != nil {
			return err
		}
		return nil
	}

	// Handle multi-row case ([]User or []*User)
	if targetValue.Kind() == reflect.Slice {
		for rows.Next() {
			elem := reflect.New(d.modelType)

			if d.isPointer {
				fieldPointers := make([]interface{}, d.modelType.NumField())
				structElem := elem.Elem()
				for i := 0; i < d.modelType.NumField(); i++ {
					fieldPointers[i] = structElem.Field(i).Addr().Interface()
				}
				if err := rows.Scan(fieldPointers...); err != nil {
					return err
				}
				targetValue.Set(reflect.Append(targetValue, structElem.Addr()))

			} else {
				fieldPointers := make([]interface{}, d.modelType.NumField())
				structElem := elem.Elem()

				for i := 0; i < d.modelType.NumField(); i++ {
					fieldPointers[i] = structElem.Field(i).Addr().Interface()
				}
				if err := rows.Scan(fieldPointers...); err != nil {
					return err
				}
				targetValue.Set(reflect.Append(targetValue, structElem))
			}
		}
		return nil
	}

	return fmt.Errorf("unsupported model type: %s", targetValue.Kind())
}
