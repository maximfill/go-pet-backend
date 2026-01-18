package todo

import (
	"context"
	"net/http"

	"github.com/maximfill/go-pet-backend/internal/repository/postgres"
)

type Service struct {
	todo *postgres.TodoRepository
}

func New(todo *postgres.TodoRepository) *Service {
	return &Service{todo: todo}
}

func (s *Service) CreateTodo(
	ctx context.Context,
	userID int64,
	title string,
) (int64, error) {

	imageURL, err := fetchRandomImage(ctx)
	if err != nil {
		imageURL = nil
	}

	return s.todo.CreateTodo(ctx, userID, title, imageURL)
}

func fetchRandomImage(ctx context.Context) (*string, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://picsum.photos/200/300",
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req) // выход в интернет
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	url := resp.Request.URL.String() // финальный url
	return &url, nil
}

func (s *Service) GetTodosByUser(
	ctx context.Context,
	userID int64,
) ([]postgres.Todo, error) {
	return s.todo.GetTodosByUser(ctx, userID)
}

func (s *Service) SetCompleted(
	ctx context.Context,
	todoID int64,
	completed bool,
) error {
	return s.todo.UpdateTodoCompleted(ctx, todoID, completed)
}

func (s *Service) DeleteTodo(
	ctx context.Context,
	todoID int64,
) error {
	return s.todo.DeleteTodo(ctx, todoID)
}
