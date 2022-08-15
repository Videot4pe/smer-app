package postgresql

import "backend/pkg/logging"

type QueryLogger interface {
	queryLogger(sql, table string, args []interface{}) *logging.Logger
}
