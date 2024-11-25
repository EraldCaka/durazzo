package util

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
)

// ResolveModelInfo extracts model information for table and type resolution
func ResolveModelInfo(model interface{}) (reflect.Type, string, bool, error) {
	modelType := reflect.TypeOf(model)
	var tableName string
	var isPointer bool

	switch {
	case modelType.Kind() == reflect.Ptr:
		if isPrimitiveType(modelType.Elem().Kind()) {
			tableName = strings.ToLower(modelType.Elem().Name())
			modelType = modelType.Elem()
			return modelType, tableName, isPointer, nil
		}
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
	case modelType.Kind() == reflect.Struct:
		tableName = strings.ToLower(modelType.Name())

	default:
		return nil, "", false, fmt.Errorf("unsupported model type: %s", modelType.Kind())
	}

	return modelType, tableName, isPointer, nil
}

// MapRowsToModel maps database rows to the provided model
func MapRowsToModel(rows *sql.Rows, model interface{}, modelType reflect.Type, isPointer bool) error {
	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() != reflect.Ptr && modelValue.Kind() != reflect.Slice {
		return errors.New("model must be a pointer to a struct or slice")
	}

	targetValue := modelValue.Elem()

	if isPrimitiveType(targetValue.Kind()) || (targetValue.Kind() == reflect.Ptr && isPrimitiveType(targetValue.Elem().Kind())) {
		if !rows.Next() {
			return errors.New("no rows found")
		}

		var value interface{}
		switch targetValue.Kind() {
		case reflect.Ptr:
			value = reflect.New(targetValue.Elem().Type()).Interface()
		default:
			value = reflect.New(targetValue.Type()).Interface()
		}

		if err := rows.Scan(value); err != nil {
			return err
		}

		if targetValue.Kind() == reflect.Ptr {
			targetValue.Set(reflect.ValueOf(value))
		} else {
			targetValue.Set(reflect.ValueOf(value).Elem())
		}
		return nil
	}

	if targetValue.Kind() == reflect.Struct ||
		(targetValue.Kind() == reflect.Ptr && targetValue.Elem().Kind() == reflect.Struct) ||
		(targetValue.Kind() == reflect.Ptr && targetValue.Elem().Kind() == reflect.Invalid) {
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

		return ScanRow(rows, targetValue)
	}

	if targetValue.Kind() == reflect.Slice {

		for rows.Next() {
			elem := reflect.New(modelType)

			if err := ScanRow(rows, elem.Elem()); err != nil {
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

// ScanRow scans a single row into a struct or a pointer to a struct
func ScanRow(rows *sql.Rows, targetValue reflect.Value) error {
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

func MapRowsToSliceModel(rows *sql.Rows, model interface{}, modelType reflect.Type) error {
	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() != reflect.Ptr {
		return errors.New("model must be a pointer to a slice")
	}

	targetValue := modelValue.Elem()
	if targetValue.Kind() != reflect.Slice {
		return fmt.Errorf("targetValue must be a slice, got %s", targetValue.Kind())
	}
	for rows.Next() {
		elem := reflect.New(modelType)
		if err := ScanRow(rows, elem.Elem()); err != nil {
			return err
		}
		targetValue.Set(reflect.Append(targetValue, elem.Elem()))
	}
	return nil
}

// isPrimitiveType checks if a type is a primitive Go type.
func isPrimitiveType(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.String, reflect.Bool, reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}
