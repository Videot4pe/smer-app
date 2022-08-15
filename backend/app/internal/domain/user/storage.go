package user

import (
	"backend/pkg/client/postgresql"
	"backend/pkg/logging"
	"context"
	"github.com/dchest/uniuri"
	"golang.org/x/crypto/bcrypt"

	db "backend/pkg/client/postgresql/model"
	sq "github.com/Masterminds/squirrel"
)

type Storage struct {
	queryBuilder sq.StatementBuilderType
	client       postgresql.Client
	logger       *logging.Logger
	ctx          context.Context
}

func NewUserStorage(ctx context.Context, client postgresql.Client, logger *logging.Logger) *Storage {
	return &Storage{
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		client:       client,
		logger:       logger,
		ctx:          ctx,
	}
}

const (
	scheme       = "public"
	table        = "users"
	refreshTable = "refresh_tokens"
)

func (s *Storage) queryLogger(sql, table string, args []interface{}) *logging.Logger {
	return s.logger.ExtraFields(map[string]interface{}{
		"sql":   sql,
		"table": table,
		"args":  args,
	})
}

func (s *Storage) All(filter *db.Filter, pagination *db.Pagination, sorts ...db.Sort) ([]User, error) {
	query := s.queryBuilder.Select("id", "email", "username", "name", "surname", "patronymic", "is_active", "created_at", "updated_at").
		From(scheme + "." + table)

	if filter != nil {
		query = filter.UseSelectBuilder(query)
	}
	if pagination != nil {
		query = pagination.UseSelectBuilder(query)
	}
	for _, sort := range sorts {
		query = sort.UseSelectBuilder(query)
	}

	sql, args, err := query.ToSql()
	logger := s.queryLogger(sql, table, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return nil, err
	}

	logger.Trace("do query")
	rows, err := s.client.Query(s.ctx, sql, args...)
	if err != nil {
		err = db.ErrDoQuery(err)
		logger.Error(err)
		return nil, err
	}

	defer rows.Close()

	list := make([]User, 0)

	for rows.Next() {
		p := User{}
		if err = rows.Scan(
			&p.Id, &p.Email, &p.Username, &p.Name, &p.Surname, &p.Patronymic, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			err = db.ErrScan(err)
			logger.Error(err)
			return nil, err
		}

		list = append(list, p)
	}

	return list, nil
}

func (s *Storage) Create(user User, isOAuth bool) (uint16, string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	lastInsertId := uint16(0)

	hash := uniuri.NewLen(15)

	query := s.queryBuilder.Insert("users").
		Columns("email", "username", "name", "surname", "patronymic", "is_active", "is_verified", "is_oauth", "password", "token_hash").
		Values(user.Email, user.Username, user.Name, user.Surname, user.Patronymic, true, isOAuth, isOAuth, hashedPassword, hash).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	logger := s.queryLogger(sql, table, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return lastInsertId, hash, err
	}

	logger.Trace("do query")
	err = s.client.QueryRow(s.ctx, sql, args...).Scan(&lastInsertId)

	if err != nil {
		logger.Error(err)
		return lastInsertId, hash, err
	}

	return lastInsertId, hash, nil
}

func (s *Storage) GetById(id uint16) (*User, error) {

	var user User

	query := s.queryBuilder.Select("id", "email", "username", "name", "surname", "patronymic", "is_active", "avatar_id").
		From("users").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	logger := s.queryLogger(sql, table, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return nil, err
	}

	logger.Trace("do query")
	row := s.client.QueryRow(s.ctx, sql, args...)

	if err = row.Scan(&user.Id, &user.Email, &user.Username, &user.Name, &user.Surname, &user.Patronymic, &user.IsActive, &user.AvatarId); err != nil {
		err = db.ErrScan(err)
		logger.Error(err)
		return nil, err
	}

	return &user, nil
}

func (s *Storage) GetByCredentials(email, password string) (uint16, bool, error) {

	var user User

	query := s.queryBuilder.Select("id", "password", "is_active", "is_verified").
		From("users").
		Where(sq.Eq{"email": email})

	sql, args, err := query.ToSql()
	logger := s.queryLogger(sql, table, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return 0, false, err
	}

	logger.Trace("do query")
	row := s.client.QueryRow(s.ctx, sql, args...)

	if err = row.Scan(&user.Id, &user.Password, &user.IsActive, &user.IsVerified); err != nil {
		err = db.ErrScan(err)
		logger.Error(err)
		return 0, false, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		err = db.ErrScan(err)
		logger.Error(err)
		return 0, false, err
	}

	return user.Id, user.IsVerified, nil
}

func (s *Storage) Update(id uint16, user User) error {

	query := s.queryBuilder.Update("users").
		//Set("email", user.Email).
		Set("username", user.Username).
		Set("name", user.Name).
		Set("surname", user.Surname).
		Set("patronymic", user.Patronymic).
		Set("avatar_id", user.AvatarId).
		Where(sq.Eq{"id": id})

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

func (s *Storage) Activate(hash string) error {
	query := s.queryBuilder.Update("users").
		Set("is_verified", true).
		Set("is_active", true).
		Where(sq.Eq{"is_verified": false, "token_hash": hash})

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

func (s *Storage) Delete(id uint16) error {

	query := s.queryBuilder.Update("users").
		Set("is_active", false).
		Where(sq.Eq{"id": id})
	//.Delete("users").

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

func (s *Storage) GetByEmailAndGenerateHash(email string) (uint16, bool, string, error) {

	var user User

	query := s.queryBuilder.Select("id", "is_verified").
		From(scheme + "." + table).
		Where(sq.Eq{"email": email})

	hash := uniuri.NewLen(15)

	sql, args, err := query.ToSql()
	logger := s.queryLogger(sql, table, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return 0, false, hash, err
	}

	row := s.client.QueryRow(s.ctx, sql, args...)

	if err = row.Scan(&user.Id, &user.IsVerified); err != nil {
		err = db.ErrScan(err)
		logger.Error(err)
		return 0, false, hash, err
	}

	setTokenQuery := s.queryBuilder.Update(scheme+"."+table).Set("token_hash", hash).Where(sq.Eq{"id": user.Id})

	sql, args, err = setTokenQuery.ToSql()
	logger = s.queryLogger(sql, table, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return user.Id, user.IsVerified, hash, err
	}

	logger.Trace("do query")
	_, err = s.client.Exec(s.ctx, sql, args...)

	if err != nil {
		logger.Error(err)
		return user.Id, user.IsVerified, hash, err
	}

	return user.Id, user.IsVerified, hash, nil
}

func (s *Storage) PasswordReset(id uint16) error {

	query := s.queryBuilder.Delete("users").
		Where(sq.Eq{"id": id})

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

func (s *Storage) IsRefreshTokenActual(token string) (uint16, error) {
	query := s.queryBuilder.Select("user_id").
		From(scheme + "." + refreshTable).
		Where(sq.Eq{"token": token})

	sql, args, err := query.ToSql()
	logger := s.queryLogger(sql, table, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return 0, err
	}

	var userId uint16

	row := s.client.QueryRow(s.ctx, sql, args...)

	if err = row.Scan(&userId); err != nil {
		err = db.ErrScan(err)
		logger.Error(err)
		return 0, err
	}

	return userId, nil
}

func (s *Storage) UpdateRefreshToken(token string, userId uint16) error {

	removeQuery := s.queryBuilder.Delete("refresh_tokens").
		Where(sq.Eq{"user_id": userId})

	sql, args, err := removeQuery.ToSql()
	logger := s.queryLogger(sql, refreshTable, args)
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

	insertQuery := s.queryBuilder.Insert("refresh_tokens").
		Columns("user_id", "token").
		Values(userId, token)

	sql, args, err = insertQuery.ToSql()
	logger = s.queryLogger(sql, refreshTable, args)
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
