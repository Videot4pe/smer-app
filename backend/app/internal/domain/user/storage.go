package user

import (
	"backend/pkg/auth"
	"backend/pkg/client/postgresql"
	"backend/pkg/logging"
	"context"

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
	scheme      = "public"
	table       = "users"
	tokensTable = "tokens"
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
	token := ""

	// Creating user
	query := s.queryBuilder.Insert(table).
		Columns("email", "username", "name", "surname", "patronymic", "is_active", "is_verified", "is_oauth", "password").
		Values(user.Email, user.Username, user.Name, user.Surname, user.Patronymic, true, isOAuth, isOAuth, hashedPassword).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	logger := s.queryLogger(sql, table, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return lastInsertId, token, err
	}

	logger.Trace("Creating user")
	err = s.client.QueryRow(s.ctx, sql, args...).Scan(&lastInsertId)

	if err != nil {
		logger.Error(err)
		return lastInsertId, token, err
	}

	jwt := auth.LinkJwt{
		Data: auth.LinkJwtData{
			Id: lastInsertId,
		},
	}

	token, err = auth.Encode(&jwt, 10)

	tokenQuery := s.queryBuilder.Insert(tokensTable).
		Columns("user_id", "token", "token_type").
		Values(lastInsertId, token, "ACTIVATE")

	sql, args, err = tokenQuery.ToSql()
	logger = s.queryLogger(sql, tokensTable, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return lastInsertId, token, err
	}

	logger.Trace("Creating activation token\n")
	_, err = s.client.Exec(s.ctx, sql, args...)

	if err != nil {
		logger.Error(err)
		return lastInsertId, token, err
	}

	return lastInsertId, token, nil
}

func (s *Storage) GetById(id uint16) (*User, error) {

	var user User

	query := s.queryBuilder.Select("id", "email", "username", "name", "surname", "patronymic", "is_active", "avatar_id").
		From(table).
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	logger := s.queryLogger(sql, table, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return nil, err
	}

	logger.Trace("Getting user by id")
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
		From(table).
		Where(sq.Eq{"email": email})

	sql, args, err := query.ToSql()
	logger := s.queryLogger(sql, table, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return 0, false, err
	}

	logger.Trace("Getting user by credentials")
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
	query := s.queryBuilder.Update(table).
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

	logger.Trace("Updating user")
	_, err = s.client.Exec(s.ctx, sql, args...)

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (s *Storage) Activate(token string) error {
	_, linkJwt, err := auth.Decode(&auth.LinkJwt{}, token)

	if err != nil {
		s.logger.Error("Activation token decode error\n", err)
		s.removeToken(token, "ACTIVATE")
		return err
	}

	userId := linkJwt.Data.Id
	query := s.queryBuilder.Update(table).
		Set("is_verified", true).
		Set("is_active", true).
		Where(sq.Eq{"id": userId})

	sql, args, err := query.ToSql()
	logger := s.queryLogger(sql, table, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return err
	}

	logger.Trace("Activating user account")
	_, err = s.client.Exec(s.ctx, sql, args...)

	if err != nil {
		logger.Error(err)
		return err
	}

	s.removeToken(token, "ACTIVATE")

	return nil
}

func (s *Storage) Delete(id uint16) error {
	// TODO удалять 'остатки' пользователя
	// в виде токенов и прочего

	query := s.queryBuilder.Update(table).
		Set("is_active", false).
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	logger := s.queryLogger(sql, table, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return err
	}

	logger.Trace("Deleting user")
	_, err = s.client.Exec(s.ctx, sql, args...)

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (s *Storage) GetByEmail(email string) (uint16, bool, error) {
	var user User

	query := s.queryBuilder.Select("id", "is_verified").
		From(scheme + "." + table).
		Where(sq.Eq{"email": email})

	sql, args, err := query.ToSql()
	logger := s.queryLogger(sql, table, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return 0, false, err
	}

	logger.Trace("Getting user")
	err = s.client.QueryRow(s.ctx, sql, args...).Scan(&user.Id, &user.IsVerified)

	if err != nil {
		err = db.ErrScan(err)
		logger.Error(err)
		return 0, false, err
	}

	return user.Id, user.IsVerified, nil
}

func (s *Storage) PasswordReset(userId uint16) (string, error) {
	s.removeTokenByUserId(userId, "RESET_PASS")

	jwt := auth.LinkJwt{
		Data: auth.LinkJwtData{
			Id: userId,
		},
	}

	token, err := auth.Encode(&jwt, 10)

	tokenQuery := s.queryBuilder.Insert(tokensTable).
		Columns("user_id", "token", "token_type").
		Values(userId, token, "RESET_PASS")

	sql, args, err := tokenQuery.ToSql()
	logger := s.queryLogger(sql, tokensTable, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return token, err
	}

	logger.Trace("Creating reset token")
	_, err = s.client.Exec(s.ctx, sql, args...)

	if err != nil {
		logger.Error(err)
		return token, err
	}

	return token, nil
}

func (s *Storage) ChangePassword(token string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	_, _, err = auth.Decode(&auth.LinkJwt{}, token)
	if err != nil {
		s.removeToken(token, "RESET_PASS")
		return err
	}

	userId, err := s.getUserIdByToken(token, "RESET_PASS")
	if err != nil {
		return err
	}

	query := s.queryBuilder.Update(table).
		Set("password", hashedPassword).
		Where(sq.Eq{"id": userId})

	sql, args, err := query.ToSql()
	logger := s.queryLogger(sql, table, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return err
	}

	logger.Trace("Updating password")
	_, err = s.client.Exec(s.ctx, sql, args...)

	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Trace("Removing reset-password token")
	s.removeToken(token, "RESET_PASS")

	return nil
}

func (s *Storage) IsRefreshTokenActual(token string) (uint16, error) {
	query := s.queryBuilder.Select("user_id").
		From(scheme + "." + token).
		Where(sq.Eq{"token": token, "token_type": "AUTH"})

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
	removeQuery := s.queryBuilder.Delete(tokensTable).
		Where(sq.Eq{"user_id": userId, "token_type": "AUTH"})

	sql, args, err := removeQuery.ToSql()
	logger := s.queryLogger(sql, tokensTable, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return err
	}

	logger.Trace("Deleting refresh (auth) token")
	_, err = s.client.Exec(s.ctx, sql, args...)

	if err != nil {
		logger.Error(err)
		return err
	}

	insertQuery := s.queryBuilder.Insert(tokensTable).
		Columns("user_id", "token", "token_type").
		Values(userId, token, "AUTH")

	sql, args, err = insertQuery.ToSql()
	logger = s.queryLogger(sql, token, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return err
	}

	logger.Trace("Creating new refresh (auth) token")
	_, err = s.client.Exec(s.ctx, sql, args...)

	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

// TODO token_type string -> enum (?)
func (s *Storage) getUserIdByToken(token string, token_type string) (uint16, error) {
	var userId uint16 = 0

	tokenQuery := s.queryBuilder.
		Select("user_id").
		From(tokensTable).
		Where(sq.Eq{"token": token, "token_type": token_type})

	sql, args, err := tokenQuery.ToSql()
	logger := s.queryLogger(sql, tokensTable, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return userId, err
	}

	logger.Trace("Searching user by reset token")
	err = s.client.QueryRow(s.ctx, sql, args...).Scan(&userId)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return userId, err
	}

	return userId, nil
}

func (s *Storage) removeToken(token string, tokenType string) {
	removeQuery := s.queryBuilder.Delete(tokensTable).Where(sq.Eq{"token": token, "token_type": tokenType})

	sql, args, err := removeQuery.ToSql()
	logger := s.queryLogger(sql, tokensTable, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
	}

	_, err = s.client.Exec(s.ctx, sql, args...)

	if err != nil {
		logger.Error(err)
	}
}

func (s *Storage) removeTokenByUserId(user_id uint16, tokenType string) {
	removeQuery := s.queryBuilder.Delete(tokensTable).Where(sq.Eq{"user_id": user_id, "token_type": tokenType})

	sql, args, err := removeQuery.ToSql()
	logger := s.queryLogger(sql, tokensTable, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
	}

	_, err = s.client.Exec(s.ctx, sql, args...)

	if err != nil {
		logger.Error(err)
	}
}
