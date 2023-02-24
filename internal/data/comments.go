package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Comment struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	AuthorID  int64     `json:"author_id"`
	PostID    int64     `json:"post_id"`
	Content   string    `json:"content"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CommentModel struct {
	DB *sql.DB
}

func (m CommentModel) Insert(comment *Comment) error {
	query := `
		INSERT INTO comments(created_at, updated_at, author_id, post_id, content)
		VALUES (NOW(), NOW(), $1, $2, $3)
		RETURNING id`

	return m.DB.QueryRow(query, &comment.AuthorID, &comment.PostID, &comment.Content).Scan(&comment.ID)
}

func (c CommentModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	// Construct the SQL query to delete the record.
	query := `
		DELETE FROM comments
		WHERE id = $1`

	result, err := c.DB.Exec(query, id)
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

func (c CommentModel) Update(comment *Comment) error {
	query := `
	UPDATE comments
	SET content = $1, updated_at = NOW()
	WHERE id = $2
	RETURNING created_at, author_id, updated_at, post_id`
	args := []any{
		comment.Content,
		comment.ID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := c.DB.QueryRowContext(ctx, query, args...).Scan(&comment.AuthorID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (c CommentModel) Get(id int64) (*Comment, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, author_id, post_id, content, created_at, updated_at
		FROM comments
		WHERE id = $1`

	var comment Comment

	err := c.DB.QueryRow(query, id).Scan(
		&comment.ID,
		&comment.CreatedAt,
		&comment.PostID,
		&comment.Content,
		&comment.UpdatedAt,
		&comment.AuthorID,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &comment, nil
}

func (c CommentModel) GetAllFromPost(postID int) ([]*Comment, error) {
	query := `
		SELECT id, created_at, author_id, content, post_id, updated_at
		FROM comments
		WHERE post_id = $1
		ORDER BY created_at`
	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := c.DB.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var comments []*Comment

	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&comment.ID,
			&comment.CreatedAt,
			&comment.AuthorID,
			&comment.Content,
			&comment.PostID,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
