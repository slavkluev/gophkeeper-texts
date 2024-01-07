package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"texts/internal/domain/models"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) GetAll(ctx context.Context, userUID uint64) ([]models.Text, error) {
	const op = "storage.sqlite.GetAll"

	stmt, err := s.db.Prepare("SELECT id, text, info FROM texts WHERE user_uid = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.QueryContext(ctx, userUID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var texts []models.Text

	for rows.Next() {
		var text models.Text
		err = rows.Scan(&text.ID, &text.Text, &text.Info)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		texts = append(texts, text)
	}

	return texts, nil
}

func (s *Storage) SaveText(ctx context.Context, text string, info string, userUID uint64) (uint64, error) {
	const op = "storage.sqlite.SaveText"

	stmt, err := s.db.Prepare("INSERT INTO texts(text, info, user_uid) VALUES(?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, text, info, userUID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return uint64(id), nil
}

func (s *Storage) UpdateText(ctx context.Context, id uint64, text string, info string, userUID uint64) error {
	const op = "storage.sqlite.UpdateText"

	stmt, err := s.db.Prepare("UPDATE texts SET text = ?, info = ? WHERE id = ? AND user_uid = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, text, info, id, userUID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
