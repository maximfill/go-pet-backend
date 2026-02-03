package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"

	"github.com/maximfill/go-pet-backend/internal/auth"
	postman "github.com/maximfill/go-pet-backend/internal/clients"
	"github.com/maximfill/go-pet-backend/internal/repository/postgres"
	todoservice "github.com/maximfill/go-pet-backend/internal/service/todo"
	userservice "github.com/maximfill/go-pet-backend/internal/service/user"
	grpcTodo "github.com/maximfill/go-pet-backend/internal/transport/grpc/todo"
	httptransport "github.com/maximfill/go-pet-backend/internal/transport/http"
)

func main() {
	fmt.Println("API started")

	// ===== DB =====
	db, err := postgres.New(postgres.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := postgres.Migrate(db); err != nil {
		log.Fatal(err)
	}

	// ===== Repos =====
	userRepo := postgres.NewUserRepository(db)
	todoRepo := postgres.NewTodoRepository(db)

	// ===== External gRPC =====
	postmanConn, err := postman.NewConn("grpc.postman-echo.com:443")
	if err != nil {
		log.Fatal(err)
	}
	echoClient := postman.NewEchoClient(postmanConn)

	// ===== Services =====
	userService := userservice.New(userRepo)
	todoService := todoservice.New(
		todoRepo,
		echoClient,
	)

	// ===== HTTP =====
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(2 * time.Second))

	userHandler := httptransport.NewUserHandler(userService)
	todoHandler := httptransport.NewTodoHandler(todoService)

	// ---------- PUBLIC ----------
	r.Post("/register", userHandler.Register)
	r.Post("/login", userHandler.Login)

	// ---------- PROTECTED ----------
	r.Group(func(r chi.Router) {
		r.Use(httptransport.AuthMiddleware)

		r.Post("/todos", todoHandler.Create)
		r.Get("/todos", todoHandler.List)
		r.Delete("/todos/{id}", todoHandler.Delete)
	})

	// ===== gRPC =====
	go func() {
		lis, _ := net.Listen("tcp", ":50051")

		grpcServer := grpc.NewServer(
			grpc.UnaryInterceptor(auth.UnaryAuthInterceptor),
		)

		grpcTodo.RegisterTodoServiceServer(
			grpcServer,
			grpcTodo.NewServer(todoService),
		)

		log.Println("gRPC started on :50051")
		log.Fatal(grpcServer.Serve(lis))
	}()

	log.Fatal(http.ListenAndServe(":8080", r))
}
