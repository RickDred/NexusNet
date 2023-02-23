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
	CreatedAt   time.Time `json:"created_at,omitempty"`
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

func (p PostModel) GetAll(title string, authorID int, authorName string, filters Filters) ([]*Post, Metadata, error) {
	// Update the SQL query to include the window function which counts the total
	// (filtered) records.
	// (filtered) records.
	query := fmt.Sprintf(`
		SELECT p.id, p.created_at, p.title, p.description, p.author_id, p.updated_at
		FROM posts p, users u
		WHERE (LOWER(p.title) = LOWER($1) OR $1 = '')
		AND ((u.id) = ($2) OR $2 = 0)
		AND (LOWER(u.name) = LOWER($3) OR $3 = '')
		ORDER BY %s %s, id ASC
		LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.sortDirection())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []any{title, authorID, authorName, filters.limit(), filters.offset()}
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

func (p PostModel) Update(post *Post) error {
	query := `
	UPDATE posts
	SET title = $1, description = $2, updated_at = NOW()
	WHERE id = $3
	RETURNING created_at, author_id, id`
	args := []any{
		post.Title,
		post.Description,
		post.ID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := p.DB.QueryRowContext(ctx, query, args...).Scan(&post.ID)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail

		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}
