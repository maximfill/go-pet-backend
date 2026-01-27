package todo

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	auth "github.com/maximfill/go-pet-backend/internal/auth"
	todoservice "github.com/maximfill/go-pet-backend/internal/service/todo"
)

type Server struct {
	UnimplementedTodoServiceServer
	service *todoservice.Service
}

func NewServer(service *todoservice.Service) *Server {
	return &Server{service: service}
}

func (s *Server) CreateTodo(
	ctx context.Context,
	req *CreateTodoRequest,
) (*CreateTodoResponse, error) {

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no user")
	}

	id, err := s.service.CreateTodo(ctx, userID, req.Title)
	if err != nil {
		return nil, err
	}

	return &CreateTodoResponse{
		Id: id,
	}, nil
}

func (s *Server) ListTodos(
	ctx context.Context,
	req *ListTodosRequest,
) (*ListTodosResponse, error) {

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no user")
	}

	todos, err := s.service.GetTodosByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	resp := &ListTodosResponse{}
	for _, t := range todos {
		item := &Todo{
			Id:        t.ID,
			Title:     t.Title,
			Completed: t.Completed,
		}
		if t.ImageURL.Valid {
			item.ImageUrl = t.ImageURL.String
		}
		resp.Todos = append(resp.Todos, item)
	}

	return resp, nil
}
