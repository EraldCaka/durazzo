package durazzo

//
//import (
//	"database/sql"
//	"fmt"
//	"github.com/EraldCaka/durazzo/pkg/util"
//	"log"
//	"reflect"
//)
//
//func (d *Durazzo) Raw(query string, args ...interface{}) *RawQuery {
//
//	return &RawQuery{
//		Durazzo: d,
//		query:   query,
//		args:    args,
//	}
//}
//
//type RawQuery struct {
//	*Durazzo
//	query string
//	args  []interface{}
//	model interface{}
//}
//
//// Model sets the target model to map the results
//func (rq *RawQuery) Model(model interface{}) *RawQuery {
//	rq.model = model
//	return rq
//}
//
//// Run executes the raw query and maps the results
//func (rq *RawQuery) Run() error {
//	log.Println("Executing query:", rq.query, "with args:", rq.args)
//	log.Println(rq.args, "args")
//	rows, err := rq.Durazzo.Db.Query(rq.query, rq.args...)
//	if err != nil {
//		return fmt.Errorf("error executing raw query: %w", err)
//	}
//	defer func(rows *sql.Rows) {
//		err := rows.Close()
//		if err != nil {
//			log.Println("an error occurred with the queried rows")
//		}
//	}(rows)
//
//	if rq.model != nil {
//		modelType, _, isPointer, err := util.ResolveModelInfo(rq.model)
//		if err != nil {
//			return fmt.Errorf("error resolving model info: %w", err)
//		}
//
//		// Check if the model is a slice and handle it appropriately
//		log.Println(reflect.TypeOf(rq.model).Kind(), 1)
//		if reflect.TypeOf(rq.model).Elem().Kind() == reflect.Slice {
//			log.Println("here")
//			return util.MapRowsToSliceModel(rows, rq.model, modelType, false)
//		}
//
//		return util.MapRowsToModel(rows, rq.model, modelType, isPointer)
//	}
//
//	return nil
//}
