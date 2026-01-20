package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Todo struct {
	ID        int64
	UserID    int64
	Title     string
	Completed bool
	ImageURL  sql.NullString
	CreatedAt time.Time
}

type TodoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) CreateTodo(
	ctx context.Context,
	userID int64,
	title string,
	imageURL *string,
) (int64, error) {
	var id int64

	err := r.db.QueryRowContext(
		ctx,
		`
		INSERT INTO todos (user_id, title, image_url)
		VALUES ($1, $2, $3)
		RETURNING id
		`,
		userID,
		title,
		imageURL,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf(
			"create todo failed: user_id=%d title=%q: %w",
			userID,
			title,
			err,
		)
	}

	return id, nil
}

func (r *TodoRepository) GetTodosByUser(
	ctx context.Context,
	userID int64,
) ([]Todo, error) {

	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT id, user_id, title, completed, image_url, created_at
		FROM todos
		WHERE user_id = $1
		ORDER BY created_at DESC
		`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo

	for rows.Next() {
		var t Todo

		if err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Title,
			&t.Completed,
			&t.ImageURL,
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}

		todos = append(todos, t)
	}

	return todos, nil
}

func (r *TodoRepository) UpdateTodoCompleted(
	ctx context.Context,
	id int64,
	completed bool,
) error {

	_, err := r.db.ExecContext(
		ctx,
		`
		UPDATE todos
		SET completed = $1
		WHERE id = $2
		`,
		completed,
		id,
	)

	return err
}

func (r *TodoRepository) DeleteTodo(
	ctx context.Context,
	id int64,
) (bool, error) {

	res, err := r.db.ExecContext(
		ctx,
		`
		DELETE FROM todos
		WHERE id = $1
		`,
		id,
	)

	if err != nil {
		return false, fmt.Errorf("delete todo failed: id=%d: %w", id, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	log.Printf("[DB] rows affected=%d for id=%d", rows, id)
	return rows > 0, nil
}
