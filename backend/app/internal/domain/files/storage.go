package files

import (
	"backend/pkg/client/postgresql"
	db "backend/pkg/client/postgresql/model"
	"backend/pkg/logging"
	"context"
	sq "github.com/Masterminds/squirrel"
)

type Storage struct {
	queryBuilder sq.StatementBuilderType
	client       postgresql.Client
	logger       *logging.Logger
	ctx          context.Context
}

const (
	scheme = "public"
	table  = "files"
)

func NewFilesStorage(ctx context.Context, client postgresql.Client, logger *logging.Logger) *Storage {
	return &Storage{
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		client:       client,
		logger:       logger,
		ctx:          ctx,
	}
}

func (s *Storage) queryLogger(sql, table string, args []interface{}) *logging.Logger {
	return s.logger.ExtraFields(map[string]interface{}{
		"sql":   sql,
		"table": table,
		"args":  args,
	})
}

func (s *Storage) Create(path string, name string) (uint16, error) {
	query := s.queryBuilder.Insert(scheme+"."+table).Columns("path", "name").Values(path, name).Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	logger := s.queryLogger(sql, table, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return 0, err
	}

	rows, err := s.client.Query(s.ctx, sql, args...)
	if err != nil {
		err = db.ErrDoQuery(err)
		logger.Error(err)
		return 0, err
	}

	defer rows.Close()

	var file File
	for rows.Next() {
		if err = rows.Scan(&file.Id); err != nil {
			err = db.ErrScan(err)
			logger.Error(err)
			return 0, err
		}
	}

	return file.Id, nil
}

func (s *Storage) GetById(id uint16) (string, error) {

	var file File

	query := s.queryBuilder.Select("id", "path", "name").
		From(scheme + "." + table).
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	logger := s.queryLogger(sql, table, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return "", err
	}

	logger.Trace("do query")
	row := s.client.QueryRow(s.ctx, sql, args...)

	if err = row.Scan(&file.Id, &file.Path, &file.Name); err != nil {
		err = db.ErrScan(err)
		logger.Error(err)
		return "", err
	}

	return file.Path, nil
}
