package postgres

import (
	"context"
	"database/sql"

	"github.com/MihirSahani/Project-27/storage/entity"
)

type PostgresNoteStorageManager struct {}

func NewPostgresNoteStorageManager() *PostgresNoteStorageManager {
	return &PostgresNoteStorageManager{}
}

func (p *PostgresNoteStorageManager) CreateNote(ctx context.Context, db *sql.Tx, note *entity.Note) (*entity.Note, error) {
	query := `
		INSERT INTO notes (title, content, user_id, folder_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, content, user_id, folder_id, created_at, updated_at
	`

	err := db.QueryRowContext(ctx, query,
		note.Title,
		note.Content,
		note.UserId,
		note.FolderId,
	).Scan(&note.Id, &note.Title, &note.Content, &note.UserId, &note.FolderId, &note.CreatedAt, &note.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return note, nil
}

func (p *PostgresNoteStorageManager) DeleteNote(ctx context.Context, db *sql.Tx, id int64) (*entity.Note, error) {
	query := `
		DELETE FROM notes
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

func (p *PostgresNoteStorageManager) GetNoteByID(ctx context.Context, db *sql.Tx, id int64) (*entity.Note, error) {
	query := `
		SELECT id, title, content, created_at, updated_at
		FROM notes
		WHERE id = $1
	`
	note := &entity.Note{}
	err := db.QueryRowContext(ctx, query, id).Scan(
		&note.Id,
		&note.Title,
		&note.Content,
		&note.CreatedAt,
		&note.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (p *PostgresNoteStorageManager) UpdateNote(ctx context.Context, db *sql.Tx, note *entity.Note) (*entity.Note, error) {
	query := `
		UPDATE notes
		SET title = $1, content = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, title, content, user_id, folder_id, created_at, updated_at
	`
	err := db.QueryRowContext(ctx, query,
		note.Title,
		note.Content,
		note.Id,
	).Scan(&note.Id, &note.Title, &note.Content, &note.UserId, &note.FolderId, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return note, nil
}
