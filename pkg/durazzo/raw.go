package durazzo

import (
	"database/sql"
	"fmt"
	"github.com/EraldCaka/durazzo/pkg/util"
	"log"
	"reflect"
	"regexp"
	"strings"
)

// Raw creates a RawQuery object and ensures identifiers are quoted
func (d *Durazzo) Raw(query string, args ...interface{}) *RawQuery {
	return &RawQuery{
		Durazzo: d,
		query:   autoQuoteIdentifiers(query),
		args:    args,
	}
}

type RawQuery struct {
	*Durazzo
	query string
	args  []interface{}
	model interface{}
}

// Model sets the target model to map the results
func (rq *RawQuery) Model(model interface{}) *RawQuery {
	rq.model = model
	return rq
}

// Run executes the raw query and maps the results
func (rq *RawQuery) Run() error {
	rows, err := rq.Durazzo.Db.Query(rq.query, rq.args...)
	if err != nil {
		return fmt.Errorf("error executing raw query: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("an error occurred with the queried rows")
		}
	}(rows)

	if rq.model != nil {
		modelType, _, isPointer, err := util.ResolveModelInfo(rq.model)
		if err != nil {
			return fmt.Errorf("error resolving model info: %w", err)
		}

		if reflect.TypeOf(rq.model).Elem().Kind() == reflect.Slice {
			err = util.MapRowsToSliceModel(rows, rq.model, modelType)
			return err
		}
		err = util.MapRowsToModel(rows, rq.model, modelType, isPointer)
		return err
	}

	return nil
}

// autoQuoteIdentifiers adds quotes to table and column names in the query
func autoQuoteIdentifiers(query string) string {
	identifierRegex := regexp.MustCompile(`\b[a-zA-Z_][a-zA-Z0-9_]*\b`)

	reservedKeywords := map[string]bool{
		"select": true, "from": true, "where": true, "insert": true,
		"update": true, "delete": true, "into": true, "values": true,
		"set": true, "and": true, "or": true, "not": true, "order": true, "by": true,
		"asc": true, "desc": true, "limit": true, "exists": true, "join": true,
		"on": true, "group": true, "having": true, "like": true, "between": true,
		"distinct": true, "null": true, "is": true, "in": true, "case": true,
		"when": true, "then": true, "else": true, "end": true, "union": true,
		"intersect": true, "except": true, "all": true, "any": true, "some": true,
		"as": true, "column": true, "row": true, "table": true, "schema": true,
		"alter": true, "drop": true, "truncate": true, "rename": true, "add": true,
		"constraint": true, "foreign": true, "primary": true, "key": true,
		"unique": true, "index": true, "check": true, "default": true, "cascade": true,
		"references": true, "selective": true, "with": true, "rollback": true, "commit": true,
		"transaction": true, "savepoint": true, "procedure": true, "trigger": true,
		"function": true, "view": true, "materialized": true, "privileges": true,
		"grant": true, "revoke": true, "session": true, "audit": true, "offset": true,
		"inner": true, "outer": true, "left": true, "right": true, "count": true, "avg": true,
	}

	return identifierRegex.ReplaceAllStringFunc(query, func(match string) string {
		if _, isReserved := reservedKeywords[strings.ToLower(match)]; isReserved {
			return match
		}
		return fmt.Sprintf(`"%s"`, match)
	})
}
