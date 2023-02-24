package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Direct struct {
	ID    int64 `json:"id"`
	User1 int64 `json:"user1"`
	User2 int64 `json:"user2"`
}

type DirectModel struct {
	DB *sql.DB
}

func (d DirectModel) Insert(direct *Direct) error {
	query := `
		INSERT INTO direct(user1, user2)
		VALUES ($1, $2)
		RETURNING id`

	return d.DB.QueryRow(query, &direct.User1, &direct.User2).Scan(&direct.ID)
}

func (d DirectModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	// Construct the SQL query to delete the record.
	query := `
		DELETE FROM direct
		WHERE id = $1`

	result, err := d.DB.Exec(query, id)
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

func (d DirectModel) Get(id int64) (*Direct, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT *
		FROM direct
		WHERE id = $1`

	var direct Direct

	err := d.DB.QueryRow(query, id).Scan(
		&direct.ID,
		&direct.User1,
		&direct.User2,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &direct, nil
}

func (d DirectModel) GetAllFromUser(userId int) ([]*Direct, error) {

	query := `
		SELECT id, user1, user2
		FROM direct
		WHERE user2 = $1 or user1 = $1 
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := d.DB.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err // Update this to return an empty Metadata struct.
	}
	defer rows.Close()

	var directs []*Direct
	for rows.Next() {
		var direct Direct
		err := rows.Scan(
			&direct.ID,
			&direct.User1,
			&direct.User2,
		)
		if err != nil {
			return nil, err
		}
		directs = append(directs, &direct)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return directs, nil
}
