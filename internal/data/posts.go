package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Post struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-,omitempty"`
	Title       string    `json:"title,omitempty"`
	UpdatedAt   time.Time `json:"updated_at"`
	Description string    `json:"description"`
	AuthorID    int64     `json:"author_id"`
}

type PostModel struct {
	DB *sql.DB
}

func (m PostModel) Insert(post *Post) error {
	query := `
		INSERT INTO posts(title, author_id, description)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`

	return m.DB.QueryRow(query, &post.Title, &post.AuthorID, &post.Description).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
}

func (p PostModel) Get(id int64) (*Post, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT *
		FROM posts
		WHERE id = $1`

	var post Post

	err := p.DB.QueryRow(query, id).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.Title,
		&post.Description,
		&post.UpdatedAt,
		&post.AuthorID,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (p PostModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	// Construct the SQL query to delete the record.
	query := `
		DELETE FROM posts
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

func (p PostModel) GetAll(title string, filters Filters) ([]*Post, Metadata, error) {
	// Update the SQL query to include the window function which counts the total
	// (filtered) records.
	// (filtered) records.
	query := fmt.Sprintf(`
		SELECT id, created_at, title, description, author_id, updated_at
		FROM posts
		WHERE (LOWER(title) = LOWER($1) OR $1 = '')
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []any{title, filters.limit(), filters.offset()}
	rows, err := p.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err // Update this to return an empty Metadata struct.
	}
	defer rows.Close()
	// Declare a totalRecords variable.
	totalRecords := 0
	var posts []*Post
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.ID,
			&post.CreatedAt,
			&post.Title,
			&post.Description,
			&post.AuthorID,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, Metadata{}, err // Update this to return an empty Metadata struct.
		}
		posts = append(posts, &post)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err // Update this to return an empty Metadata struct.
	}
	// Generate a Metadata struct, passing in the total record count and pagination
	// parameters from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	// Include the metadata struct when returning.
	return posts, metadata, nil
}
