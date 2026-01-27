package user

import (
	"context"
	"errors"

	"github.com/maximfill/go-pet-backend/internal/repository/postgres"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Service struct {
	users *postgres.UserRepository
}

func New(users *postgres.UserRepository) *Service {
	return &Service{users: users}
}

func (s *Service) Register(ctx context.Context, email string, password string) (int64, error) {
	_, err := s.users.GetUserByEmail(ctx, email)
	if err == nil {
		return 0, ErrUserAlreadyExists
	}

	hash, err := hashPassword(password)
	if err != nil {
		return 0, err
	}

	return s.users.CreateUser(ctx, email, hash)
}

func (s *Service) Login(ctx context.Context, email string, password string) (int64, error) {
	user, err := s.users.GetUserByEmail(ctx, email)
	if err != nil {
		return 0, ErrInvalidCredentials
	}

	if !checkPassword(password, user.PasswordHash) {
		return 0, ErrInvalidCredentials
	}

	return user.ID, nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
	return err == nil
}
