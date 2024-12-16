package repo

import (
	"avito-test-task/internal/domain"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PostgresUserRepo struct {
	db           IPool
	retryAdapter IPostgresRetryAdapter
}

func NewPostrgesUserRepo(db IPool, retryAdapter IPostgresRetryAdapter) *PostgresUserRepo {
	return &PostgresUserRepo{
		db:           db,
		retryAdapter: retryAdapter,
	}
}

func (p *PostgresUserRepo) Create(ctx context.Context, user *domain.User, lg *zap.Logger) error {
	lg.Info("create user", zap.String("user_id", user.UserID.String()))

	query := `insert into users(user_id, mail, password, role) values ($1, $2, $3, $4)`
	_, err := p.db.Exec(ctx, query, user.UserID, user.Mail, user.Password, user.Role)
	if err != nil {
		lg.Warn("postgres create user error", zap.Error(err))
		return err
	}

	return nil
}

func (p *PostgresUserRepo) DeleteByID(ctx context.Context, id uuid.UUID, lg *zap.Logger) error {
	lg.Info("delete user", zap.String("user_id", id.String()))

	query := `delete from users where user_id=$1`
	_, err := p.db.Exec(ctx, query, id)
	if err != nil {
		lg.Warn("postgres delete user error", zap.Error(err))
		return err
	}

	return nil
}

func (p *PostgresUserRepo) Update(ctx context.Context, newUserData *domain.User, lg *zap.Logger) error {
	lg.Info("update user", zap.String("user_id", newUserData.UserID.String()))

	query := `update users set user_id=$1,	
			mail=$2,
			password=$3,
			role=$4 where user_id=$1`
	_, err := p.db.Exec(ctx, query, newUserData.UserID, newUserData.Mail,
		newUserData.Password, newUserData.Role)
	if err != nil {
		lg.Warn("postgres update user error", zap.Error(err))
		return err
	}

	return nil
}

func (p *PostgresUserRepo) GetByID(ctx context.Context, id uuid.UUID, lg *zap.Logger) (domain.User, error) {
	var user domain.User
	lg.Info("get user by id", zap.String("user_id", id.String()))

	query := `select * from users where user_id=$1`
	rows := p.db.QueryRow(ctx, query, id)

	err := rows.Scan(&user.UserID, &user.Mail, &user.Password, &user.Role)
	if err != nil {
		lg.Warn("postgres get by id user error", zap.Error(err))
		return domain.User{}, err
	}

	return user, nil
}

func (p *PostgresUserRepo) GetAll(ctx context.Context, offset int, limit int, lg *zap.Logger) ([]domain.User, error) {
	lg.Info("get users", zap.Int("offset", offset), zap.Int("limit", limit))

	query := `select * from users order by user_id limit $1 offset $2`
	rows, err := p.db.Query(ctx, query, limit, offset)
	defer rows.Close()
	if err != nil {
		lg.Warn("user repo: get all error", zap.Error(err))
		return nil,
			fmt.Errorf("user repo: get all error: %v", err.Error())
	}

	var (
		users []domain.User
		user  domain.User
	)
	for rows.Next() {
		err = rows.Scan(&user.UserID, &user.Mail, &user.Password, &user.Role)
		if err != nil {
			lg.Warn("postgres user get all error: scan user error")
			continue
		}
		users = append(users, user)
	}

	return users, err
}

func (p *PostgresUserRepo) CreateCode(ctx context.Context, user *domain.User, codeHash string, lg *zap.Logger) error {
	lg.Info("create or update code", zap.String("user_id", user.UserID.String()))

	query := `UPDATE codes SET code = $1 WHERE user_id = $2`
	result, err := p.db.Exec(ctx, query, codeHash, user.UserID)
	if err != nil {
		lg.Warn("postgres update code error", zap.Error(err))
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		query := `INSERT INTO codes (user_id, code) VALUES ($1, $2)`
		_, err = p.db.Exec(ctx, query, user.UserID, codeHash)
		if err != nil {
			lg.Warn("postgres insert code error", zap.Error(err))
			return err
		}
	}

	return nil
}

func (p *PostgresUserRepo) GetHashCode(ctx context.Context, user *domain.User, lg *zap.Logger) (string, error) {
	lg.Info("get hash code", zap.String("user_id", user.UserID.String()))

	query := `SELECT code FROM codes WHERE user_id = $1`
	var code string
	err := p.db.QueryRow(ctx, query, user.UserID).Scan(&code)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("hash code not found for user %s", user.UserID)
		}
		lg.Warn("postgres get hash code error", zap.Error(err))
		return "", err
	}

	return code, nil
}
