package durazzo

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/EraldCaka/durazzo/pkg/util"
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
// MUST be a pointer of a type
func (d *Durazzo) Select(model interface{}) *SelectType {
	modelType, tableName, isPointer, err := util.ResolveModelInfo(model)

	if err != nil {
		log.Println("failed to initialize SelectType: ", err)
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

// Where adds a where condition inside the query
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

// Run executes the query asynchronously using a dedicated channel
func (st *SelectType) Run() error {
	resultChan := make(chan error, 1)
	go func() {
		startTime := time.Now()

		query, err := st.queryBuilder.BuildSelectQuery(st.tableName, st.conditions, st.limit)
		if err != nil {
			resultChan <- err
			close(resultChan)
			return
		}

		rows, err := st.Durazzo.Db.Query(query, st.args...)
		if err != nil {
			resultChan <- err
			close(resultChan)
			return
		}
		defer func(rows *sql.Rows) {
			if err := rows.Close(); err != nil {
				log.Println("an error occurred with the queried rows:", err)
			}
		}(rows)

		err = util.MapRowsToModel(rows, st.model, st.modelType, st.isPointer)
		elapsedTime := time.Since(startTime)
		log.Printf("Query : %s took %v to run\n\n", query, elapsedTime)

		resultChan <- err
		close(resultChan)
	}()
	return <-resultChan
}
