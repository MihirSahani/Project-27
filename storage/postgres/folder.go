package postgres

import (
	"context"
	"database/sql"

	"github.com/MihirSahani/Project-27/storage/entity"
)

type PostgresFolderStorageManager struct{}

func NewPostgresFolderStorageManager() *PostgresFolderStorageManager {
	return &PostgresFolderStorageManager{}
}

func (p *PostgresFolderStorageManager) CreateFolder(ctx context.Context, db *sql.Tx, folder *entity.Folder) (*entity.Folder, error) {
	query := `
		INSERT INTO folders (name, user_id)
		VALUES ($1, $2)
		RETURNING id, name, user_id, created_at, updated_at
	`
	err := db.QueryRowContext(ctx, query,
		folder.Name,
		folder.UserId,
	).Scan(&folder.Id, &folder.Name, &folder.UserId, &folder.CreatedAt, &folder.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return folder, nil
}

func (p *PostgresFolderStorageManager) DeleteFolder(ctx context.Context, db *sql.Tx, id int64) (*entity.Folder, error) {
	query := `
		DELETE FROM folders
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

	return nil, nil
}

func (p *PostgresFolderStorageManager) GetAllFolders(ctx context.Context, db *sql.Tx, id int64) ([]*entity.Folder, error) {
	query := `
		SELECT id, name, user_id, created_at, updated_at
		FROM folders
		WHERE user_id = $1
	`
	rows, err := db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	folders := []*entity.Folder{}
	for rows.Next() {
		folder := &entity.Folder{}
		err := rows.Scan(
			&folder.Id,
			&folder.Name,
			&folder.UserId,
			&folder.CreatedAt,
			&folder.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		folders = append(folders, folder)
	}
	return folders, nil
}

func (p *PostgresFolderStorageManager) GetNotesInFolder(ctx context.Context, db *sql.Tx, folderId int64, userId int64) ([]*entity.Note, error) {
	query := `
		SELECT id, title, updated_at 
		FROM notes 
		WHERE folder_id=$1 AND user_id=$2
	`
	var notes []*entity.Note
	rows, err := db.QueryContext(ctx, query, folderId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		note := entity.Note{}
		err := rows.Scan(&note.Id, &note.Title, &note.UpdatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, &note)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return notes, nil
}
