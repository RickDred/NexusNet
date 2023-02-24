package data

import (
	"context"
	"database/sql"
	"time"
)

type Feel struct {
	Mood      string    `json:"mood"`
	CreatedAt time.Time `jsom:"created_at"`
}

type FeelModel struct {
	DB *sql.DB
}

func (f FeelModel) Insert(userid int64) error {
	query := `
	INSERT INTO feels (user_id, mood)
	VALUES ($1, $2)
	RETURNING id, created_at`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := f.DB.QueryRowContext(ctx, query, userid, "ready to enjoy").Scan()
	if err != nil {
		return err
	}
	return nil
}
