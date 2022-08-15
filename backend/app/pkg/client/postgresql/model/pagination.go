package model

import (
	"github.com/Masterminds/squirrel"
	"net/http"
	"strconv"
)

type Pagination struct {
	Page  uint64
	Limit uint64
}

func NewPagination(r *http.Request) (*Pagination, error) {
	queryValues := r.URL.Query()

	page, err := strconv.ParseUint(queryValues.Get("page"), 10, 64)
	if err != nil {
		page = 1
	}
	limit, err := strconv.ParseUint(queryValues.Get("limit"), 10, 64)
	if err != nil {
		limit = 10
	}

	return &Pagination{
		Page:  page,
		Limit: limit,
	}, nil
}

func (opt Pagination) UseSelectBuilder(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
	return builder.Limit(opt.Limit).Offset(opt.Limit * (opt.Page - 1))
}
