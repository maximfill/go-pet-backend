package main

import (
	//"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/maximfill/go-pet-backend/internal/repository/postgres"
	todoservice "github.com/maximfill/go-pet-backend/internal/service/todo"
	userservice "github.com/maximfill/go-pet-backend/internal/service/user"
	httptransport "github.com/maximfill/go-pet-backend/internal/transport/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"google.golang.org/grpc" // ⬅ сам gRPC сервер библиотека (сервер, интерсепторы и т.д.)
	"net"                    // ⬅ TCP listener нужен, чтобы открыть порт (:50051)

	authgrpc "github.com/maximfill/go-pet-backend/internal/auth"
	todogrpc "github.com/maximfill/go-pet-backend/internal/transport/grpc/todo" // твой gRPC transport слой, сгенерированный + server.go
)

func main() {
	fmt.Println("API started")

	// 1. Конфигурация БД
	cfg := postgres.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}

	// 2. Подключение к БД
	db, err := postgres.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// 3. Миграции
	if err := postgres.Migrate(db); err != nil {
		log.Fatal(err)
	}

	// Repository
	repo := postgres.NewUserRepository(db)
	todoRepo := postgres.NewTodoRepository(db)

	// Service
	service := userservice.New(repo)
	todoService := todoservice.New(todoRepo)

	// HTTP handler
	handler := httptransport.NewUserHandler(service)
	todoHandler := httptransport.NewTodoHandler(todoService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(2 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	r.Get("/todos", todoHandler.List)
	r.Post("/register", handler.Register)
	r.Post("/login", handler.Login)

	r.Post("/todos", todoHandler.Create)
	r.Patch("/todos/{id}", todoHandler.Update)

	r.Delete("/todos/{id}", todoHandler.Delete)

	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatal(err)
		}

		grpcServer := grpc.NewServer(
			grpc.UnaryInterceptor(authgrpc.UnaryAuthInterceptor),
		)

		todogrpc.RegisterTodoServiceServer(
			grpcServer,
			todogrpc.NewServer(todoService),
		)

		log.Println("gRPC server started on :50051")
		log.Fatal(grpcServer.Serve(lis))
	}()

	log.Fatal(http.ListenAndServe(":8080", r))
}
