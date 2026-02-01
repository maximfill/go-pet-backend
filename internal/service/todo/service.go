package todo

import (
	"context"
	"log"
	"net/http"

	"github.com/maximfill/go-pet-backend/internal/repository/postgres"
)

type TodoRepository interface {
	CreateTodo(ctx context.Context, userID int64, title string, imageURL *string) (int64, error)
	GetTodosByUser(ctx context.Context, userID int64) ([]postgres.Todo, error)
	UpdateTodoCompleted(ctx context.Context, todoID int64, completed bool) error
	DeleteTodo(ctx context.Context, todoID int64) (bool, error)
}

type Echo interface {
	Echo(ctx context.Context, msg string) error
}

type Service struct {
	repo TodoRepository
	echo Echo
}

func New(repo TodoRepository, echo Echo) *Service {
	return &Service{
		repo: repo,
		echo: echo,
	}
}

func (s *Service) CreateTodo(
	ctx context.Context,
	userID int64,
	title string,
) (int64, error) {

	if s.echo != nil {
		log.Println("sending echo")
		_ = s.echo.Echo(ctx, "todo created")
	}

	imageURL, err := fetchRandomImage(ctx)
	if err != nil {
		imageURL = nil
	}

	return s.repo.CreateTodo(ctx, userID, title, imageURL)
}

func (s *Service) GetTodosByUser(
	ctx context.Context,
	userID int64,
) ([]postgres.Todo, error) {
	return s.repo.GetTodosByUser(ctx, userID)
}

func (s *Service) SetCompleted(
	ctx context.Context,
	todoID int64,
	completed bool,
) error {
	return s.repo.UpdateTodoCompleted(ctx, todoID, completed)
}

func (s *Service) DeleteTodo(
	ctx context.Context,
	todoID int64,
) (bool, error) {
	return s.repo.DeleteTodo(ctx, todoID)
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

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	url := resp.Request.URL.String()
	return &url, nil
}
