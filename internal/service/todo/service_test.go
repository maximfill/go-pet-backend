package todo

import (
	"context"
	"errors"
	"testing"

	"github.com/maximfill/go-pet-backend/internal/repository/postgres"
)

/*
fakeRepo — единый фейковый репозиторий.
*/
type fakeRepo struct {
	// delete
	deleteResult bool
	deleteError  error
	deleteCalledWithID int64

	// create
	createResultID int64
	createError    error
	createCalledWithUserID int64
	createCalledWithTitle  string

	// read
	readResult []postgres.Todo
	readError  error
	readCalledWithUserID int64
}

// ====== Create ======

func (f *fakeRepo) CreateTodo(
	ctx context.Context,
	userID int64,
	title string,
	imageURL *string,
) (int64, error) {
	f.createCalledWithUserID = userID
	f.createCalledWithTitle = title
	return f.createResultID, f.createError
}

// ====== Read ======

func (f *fakeRepo) GetTodosByUser(
	ctx context.Context,
	userID int64,
) ([]postgres.Todo, error) {
	f.readCalledWithUserID = userID
	return f.readResult, f.readError
}

// ====== Update ======

func (f *fakeRepo) UpdateTodoCompleted(
	ctx context.Context,
	todoID int64,
	completed bool,
) error {
	return nil
}

// ====== Delete ======

func (f *fakeRepo) DeleteTodo(
	ctx context.Context,
	todoID int64,
) (bool, error) {
	f.deleteCalledWithID = todoID
	return f.deleteResult, f.deleteError
}

//
// ================== TESTS ==================
//

// ---------- Delete ----------

func TestService_DeleteTodo_Success(t *testing.T) {
	repo := &fakeRepo{
		deleteResult: true,
	}

	service := New(repo, nil)

	ok, err := service.DeleteTodo(context.Background(), 42)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatalf("expected true, got false")
	}
	if repo.deleteCalledWithID != 42 {
		t.Fatalf("expected id=42, got %d", repo.deleteCalledWithID)
	}
}

func TestService_DeleteTodo_Error(t *testing.T) {
	repo := &fakeRepo{
		deleteError: errors.New("db error"),
	}

	service := New(repo, nil)

	_, err := service.DeleteTodo(context.Background(), 1)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

// ---------- Create ----------

func TestService_CreateTodo_Success(t *testing.T) {
	repo := &fakeRepo{
		createResultID: 10,
	}

	service := New(repo, nil)

	id, err := service.CreateTodo(context.Background(), 5, "hello")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 10 {
		t.Fatalf("expected id=10, got %d", id)
	}
	if repo.createCalledWithUserID != 5 {
		t.Fatalf("expected userID=5, got %d", repo.createCalledWithUserID)
	}
	if repo.createCalledWithTitle != "hello" {
		t.Fatalf("expected title=hello, got %s", repo.createCalledWithTitle)
	}
}

// ---------- Read ----------

func TestService_GetTodosByUser_Success(t *testing.T) {
	expected := []postgres.Todo{
		{ID: 1, Title: "one"},
		{ID: 2, Title: "two"},
	}

	repo := &fakeRepo{
		readResult: expected,
	}

	service := New(repo, nil)

	result, err := service.GetTodosByUser(context.Background(), 7)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 todos, got %d", len(result))
	}
	if repo.readCalledWithUserID != 7 {
		t.Fatalf("expected userID=7, got %d", repo.readCalledWithUserID)
	}
}