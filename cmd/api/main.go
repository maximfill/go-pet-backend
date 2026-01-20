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

	log.Fatal(http.ListenAndServe(":8080", r))
}

// Routes
// 	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK) // 200
// 		w.Write([]byte("ok"))
// 	})

// 	http.HandleFunc("/register", handler.Register) // сервер ждет когда с этого пути придет запрос
// 	http.HandleFunc("/login", handler.Login)

// 	// Server Запуск сервера (ПОСЛЕДНИЙ шаг)
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

// repo := postgres.NewUserRepository(db)
// service := userservice.New(repo)

// id, err := service.Register(
// 	context.Background(),
// 	"test@test.com",
// 	"123456",
// )
// if err != nil {
// 	log.Fatal(err)
// }

// fmt.Println("created user id:", id)

//// ===== ВРЕМЕННО: тест репозитория =====

// repo := postgres.NewUserRepository(db)

// ctx := context.Background()

//// CreateUser
// id, err := repo.CreateUser(
// 	ctx,
// 	"test@example.com",
// 	"hashed_password",
// )
// if err != nil {
// 	log.Fatal("CreateUser error:", err)
// }

// fmt.Println("created user id:", id)

//// GetUserByEmail
// user, err := repo.GetUserByEmail(ctx, "test@example.com")
// if err != nil {
// 	log.Fatal("GetUserByEmail error:", err)
// }

// fmt.Println("user from db:")
// fmt.Println("id:", user.ID)
// fmt.Println("email:", user.Email)
// fmt.Println("passwordHash:", user.PasswordHash)
// fmt.Println("createdAt:", user.CreatedAt)

// ===== КОНЕЦ временного кода =====
