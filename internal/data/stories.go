package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Story struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	AuthorID  int64     `json:"author_id"`
	Visible   bool      `json:"visible"`
	Content   string    `json:"content"`
}
type StoryModel struct {
	DB *sql.DB
}

func (s StoryModel) Insert(story *Story) error {
	query := `
		INSERT INTO stories(content, author_id, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id`

	return s.DB.QueryRow(query, &story.Content, &story.AuthorID).Scan(&story.ID)
}

func (s StoryModel) Get(id int64) (*Story, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT *
		FROM stories
		WHERE id = $1`

	var story Story

	err := s.DB.QueryRow(query, id).Scan(&story.ID, &story.CreatedAt, &story.AuthorID)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &story, nil
}

func (p StoryModel) Delete(id int64) error {
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

func (s StoryModel) GetAllFromUser(authorID int) ([]*Story, error) {
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

	var stories []*Story

	for rows.Next() {
		var story Story
		err := rows.Scan(
			&story.ID,
			&story.CreatedAt,
			&story.AuthorID,
			&story.Content,
		)
		if err != nil {
			return nil, err
		}
		stories = append(stories, &story)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stories, nil
}
