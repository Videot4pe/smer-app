package smer

import (
	"backend/pkg/client/postgresql"
	db "backend/pkg/client/postgresql/model"
	"backend/pkg/logging"
	"backend/pkg/utils"
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"math"
)

type Storage struct {
	queryBuilder sq.StatementBuilderType
	client       postgresql.Client
	logger       *logging.Logger
	ctx          context.Context
}

func NewSmerStorage(ctx context.Context, client postgresql.Client, logger *logging.Logger) *Storage {
	return &Storage{
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		client:       client,
		logger:       logger,
		ctx:          ctx,
	}
}

const (
	scheme = "public"
	table  = "smers"
)

func (s *Storage) queryLogger(sql, table string, args []interface{}) *logging.Logger {
	return s.logger.ExtraFields(map[string]interface{}{
		"sql":   sql,
		"table": table,
		"args":  args,
	})
}

func (s *Storage) All(userId uint16, filters []*db.Filter, pagination *db.Pagination, sorts ...*db.Sort) ([]Smer, *utils.Meta, error) {
	query := s.queryBuilder.Select(
		"id",
		"user_id",
		"situation",
		"thoughts",
		"emotions",
		"reactions",
		"created_at",
		"updated_at",
	).From(scheme + "." + table).Where(sq.Eq{"user_id": userId})

	for _, filter := range filters {
		s.logger.Trace(filter)
		query = filter.UseSelectBuilder(query)
	}
	for _, sort := range sorts {
		s.logger.Trace(sort)
		query = sort.UseSelectBuilder(query)
	}
	if pagination != nil {
		query = pagination.UseSelectBuilder(query)
	}

	sql, args, err := query.ToSql()
	s.logger.Trace(sql)
	s.logger.Trace(args)
	logger := s.queryLogger(sql, table, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return nil, nil, err
	}

	rows, err := s.client.Query(s.ctx, sql, args...)
	if err != nil {
		err = db.ErrDoQuery(err)
		logger.Error(err)
		return nil, nil, err
	}

	defer rows.Close()

	list := make([]Smer, 0)

	for rows.Next() {
		p := Smer{}
		if err = rows.Scan(
			&p.Id, &p.UserId, &p.Situation, &p.Thoughts, &p.Emotions, &p.Reactions, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			err = db.ErrScan(err)
			logger.Error(err)
			return nil, nil, err
		}

		list = append(list, p)
	}

	var count uint64

	err = s.client.QueryRow(s.ctx, fmt.Sprintf("SELECT COUNT(*) FROM %v.%v", scheme, table)).Scan(&count)
	if err != nil {
		err = db.ErrDoQuery(err)
		logger.Error(err)
		return nil, nil, err
	}
	meta := &utils.Meta{
		TotalItems: count,
		TotalPages: uint64(math.Ceil(float64(count) / float64(pagination.Limit))),
	}

	return list, meta, nil
}

func (s *Storage) Create(smer Smer, userId uint16) (uint16, error) {

	lastInsertId := uint16(0)

	query := s.queryBuilder.Insert(scheme+"."+table).Columns(
		"user_id",
		"situation",
		"thoughts",
		"emotions",
		"reactions",
	).Values(userId, smer.Situation, smer.Thoughts, smer.Emotions, smer.Reactions).Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	logger := s.queryLogger(sql, table, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return lastInsertId, err
	}

	err = s.client.QueryRow(s.ctx, sql, args...).Scan(&lastInsertId)

	if err != nil {
		logger.Error(err)
		return lastInsertId, err
	}

	return lastInsertId, nil
}

func (s *Storage) GetById(userId uint16, id uint16) (*Smer, error) {

	var smer Smer

	query := s.queryBuilder.Select(
		"id",
		"user_id",
		"situation",
		"thoughts",
		"emotions",
		"reactions",
		"created_at",
		"updated_at",
	).From(scheme + "." + table).Where(sq.Eq{"id": id}).Where(sq.Eq{"user_id": userId})

	sql, args, err := query.ToSql()
	logger := s.queryLogger(sql, table, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return nil, err
	}

	row := s.client.QueryRow(s.ctx, sql, args...)

	if err = row.Scan(
		&smer.Id, &smer.UserId, &smer.Situation, &smer.Thoughts, &smer.Emotions, &smer.Reactions, &smer.CreatedAt, &smer.UpdatedAt,
	); err != nil {
		err = db.ErrScan(err)
		logger.Error(err)
		return nil, err
	}

	return &smer, nil
}

func (s *Storage) Update(userId uint16, id uint16, smer Smer) error {
	query := s.queryBuilder.Update(scheme+"."+table).
		Set("situation", smer.Situation).
		Set("thoughts", smer.Thoughts).
		Set("emotions", smer.Emotions).
		Set("reactions", smer.Reactions).
		Where(sq.Eq{"id": id}).Where(sq.Eq{"user_id": userId})

	sql, args, err := query.ToSql()
	logger := s.queryLogger(sql, table, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return err
	}

	logger.Trace("do query")
	_, err = s.client.Exec(s.ctx, sql, args...)

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (s *Storage) Delete(userId uint16, id uint16) error {

	query := s.queryBuilder.Delete("users").
		Where(sq.Eq{"id": id}).Where(sq.Eq{"user_id": userId})

	sql, args, err := query.ToSql()
	logger := s.queryLogger(sql, table, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return err
	}

	logger.Trace("do query")
	_, err = s.client.Exec(s.ctx, sql, args...)

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
