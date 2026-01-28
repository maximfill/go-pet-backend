package todo

import (
	"context"
	"log"
	"net/http"

	"github.com/maximfill/go-pet-backend/internal/repository/postgres"
)

// internal/service/todo/service.go

type Echo interface {
	Echo(ctx context.Context, msg string) error
}

type Service struct {
	repo *postgres.TodoRepository
	echo Echo
}

func New(repo *postgres.TodoRepository, echo Echo) *Service {
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

	//  вызов внешнего сервиса
	if s.echo != nil {
		log.Println("sending echo:", "что отправил то и вернул")
		_ = s.echo.Echo(ctx, "todo created")
	}

	// бизнес-логика
	imageURL, err := fetchRandomImage(ctx)
	if err != nil {
		imageURL = nil // допустимо
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
