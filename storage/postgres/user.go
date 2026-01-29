package postgres

import (
	"context"
	"database/sql"
	"github.com/MihirSahani/Project-27/storage/entity"
)

func NewPostgresUserStorageManager() *PostgresUserStorageManager {
	return &PostgresUserStorageManager{}
}

type PostgresUserStorageManager struct {
}

func (p *PostgresUserStorageManager) CreateUser(ctx context.Context, db *sql.Tx, user *entity.User) (*entity.User, error) {
	query := `
		INSERT INTO users (email, password_hash, first_name, last_name)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	err := db.QueryRowContext(ctx, query,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
	).Scan(&user.Id, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (p *PostgresUserStorageManager) GetUserByEmail(ctx context.Context, db *sql.Tx, email string) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, created_at
		FROM users
		WHERE email = $1
	`
	user := &entity.User{}
	err := db.QueryRowContext(ctx, query, email).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (p *PostgresUserStorageManager) GetUserByID(ctx context.Context, db *sql.Tx, id int64) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, created_at
		FROM users
		WHERE id = $1
	`
	user := &entity.User{}
	err := db.QueryRowContext(ctx, query, id).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (p *PostgresUserStorageManager) UpdateUser(ctx context.Context, db *sql.Tx, user *entity.User) (*entity.User, error) {
	query := `
		UPDATE users
		SET email = $1, password_hash = $2, first_name = $3, last_name = $4, updated_at = NOW()
		WHERE id = $5
	`
	result, err := db.ExecContext(ctx, query,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.Id,
	)

	rowsAffected, err := result.RowsAffected() 
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	return nil, err
}

func (p *PostgresUserStorageManager) DeleteUser(ctx context.Context, db *sql.Tx, id int64) (*entity.User, error) {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	result, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	return nil, err
}

func (p *PostgresUserStorageManager) ActivateUser(ctx context.Context, db *sql.Tx, id int64) (*entity.User, error) {
	query := `
		UPDATE users
		SET is_active = true
		WHERE id = $1
	`
	_, err := db.ExecContext(ctx, query, id)
	return nil,err
}
