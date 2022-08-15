package model

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Masterminds/squirrel"
)

type Sort struct {
	column string
	order  string
}

func NewSorts(r *http.Request) ([]*Sort, error) {
	queryValues := r.URL.Query()

	var sorts []*Sort
	for i := 0; ; i++ {
		keyColumn := fmt.Sprintf("sort[%d][column]", i)
		column := queryValues.Get(keyColumn)

		keyOrder := fmt.Sprintf("sort[%d][order]", i)
		order := queryValues.Get(keyOrder)

		if column == "" {
			break
		}

		sorts = append(sorts, NewSort(column, order))
	}

	return sorts, nil
}

func NewSort(column, order string) *Sort {
	order = strings.ToUpper(order)
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	return &Sort{
		column: column,
		order:  order,
	}
}

func (opt Sort) UseSelectBuilder(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
	return builder.OrderBy(opt.column + " " + opt.order)
}
