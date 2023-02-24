package data

import (
	"context"
	"database/sql"
	"time"
)

type Message struct {
	ID        int64     `json:"id"`
	DirectId  int64     `json:"direct_id"`
	CreatedAt time.Time `json:"created_at"`
	Content   string    `json:"content"`
	SenderId  int64     `json:"sender_id"`
}

type MessageModel struct {
	DB *sql.DB
}

func (m MessageModel) Insert(message *Message) error {
	query := `
		INSERT INTO messages(direct_id, content, sender_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	return m.DB.QueryRow(query, &message.DirectId, &message.Content, &message.SenderId).Scan(&message.ID, &message.CreatedAt)
}

func (m MessageModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	// Construct the SQL query to delete the record.
	query := `
		DELETE FROM messages
		WHERE id = $1`

	result, err := m.DB.Exec(query, id)
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

func (m MessageModel) GetAllFromDirect(directId int) ([]*Message, error) {

	query := `
		SELECT id, direct_id, created_at, content, sender_id
		FROM messages
		WHERE direct_id = $1
		ORDER BY id
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, directId)
	if err != nil {
		return nil, err // Update this to return an empty Metadata struct.
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		var message Message
		err := rows.Scan(
			&message.ID,
			&message.DirectId,
			&message.CreatedAt,
			&message.Content,
			&message.SenderId,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
