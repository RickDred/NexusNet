package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Storie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	AuthorID  int64     `json:"author_id"`
	Visible   bool      `json:"visible"`
	Content   string    `json:"content"`
}
type StorieModel struct {
	DB *sql.DB
}

func (s StorieModel) Insert(storie *Storie) error {
	query := `
		INSERT INTO stories(content)
		VALUES ($1)
		RETURNING id, created_at, updated_at, content, author_id`

	return s.DB.QueryRow(query, &storie.Content).Scan(&storie.ID, &storie.CreatedAt, &storie.AuthorID)
}

func (s StorieModel) Get(id int64) (*Storie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT *
		FROM stories
		WHERE id = $1`

	var storie Storie

	err := s.DB.QueryRow(query, id).Scan(&storie.ID, &storie.CreatedAt, &storie.AuthorID)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &storie, nil
}

func (p StorieModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	// Construct the SQL query to delete the record.
	query := `
		DELETE FROM stories
		WHERE id = $1`

	result, err := p.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (s StorieModel) GetAllFromUser(authorID int) ([]*Storie, error) {
	query := `
		SELECT id, created_at, author_id, content
		FROM stories
		WHERE author_id = $1
		ORDER BY created_at`
	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := s.DB.QueryContext(ctx, query, authorID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var stories []*Storie

	for rows.Next() {
		var storie Storie
		err := rows.Scan(
			&storie.ID,
			&storie.CreatedAt,
			&storie.AuthorID,
			&storie.Content,
		)
		if err != nil {
			return nil, err
		}
		stories = append(stories, &storie)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stories, nil
}
