package service

import (
	"context"
	"errors"
	"testing"

	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/domain"
	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/repository"
)

type fakeRepository struct {
	user      domain.User
	firstErr  error
	createErr error
}

func (f *fakeRepository) List(context.Context, any, repository.Query) error { return nil }

func (f *fakeRepository) First(_ context.Context, destination any, _ repository.Query) error {
	if f.firstErr != nil {
		return f.firstErr
	}
	user, ok := destination.(*domain.User)
	if ok {
		*user = f.user
	}
	return nil
}

func (f *fakeRepository) Create(_ context.Context, value any) error {
	if f.createErr != nil {
		return f.createErr
	}
	if user, ok := value.(*domain.User); ok {
		user.ID = 1
		f.user = *user
	}
	return nil
}

func (f *fakeRepository) UpdateByID(context.Context, any, uint, map[string]any) error {
	return nil
}

func (f *fakeRepository) DeleteByID(context.Context, any, uint) error { return nil }

func TestRegisterNormalizesEmailAndHashesPassword(t *testing.T) {
	repo := &fakeRepository{firstErr: repository.ErrNotFound}
	svc := New(repo)

	user, err := svc.Register(context.Background(), " Lung Care ", " USER@Example.COM ", "secure-password")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if user.Email != "user@example.com" {
		t.Fatalf("Register() email = %q, want normalized email", user.Email)
	}
	if user.Password == "secure-password" || user.Password == "" {
		t.Fatal("Register() did not hash password")
	}
}

func TestLoginRejectsUnknownUser(t *testing.T) {
	svc := New(&fakeRepository{firstErr: repository.ErrNotFound})
	_, err := svc.Login(context.Background(), "missing@example.com", "password")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("Login() error = %v, want ErrInvalidCredentials", err)
	}
}

func TestRegisterMapsDuplicateEmail(t *testing.T) {
	repo := &fakeRepository{firstErr: repository.ErrNotFound, createErr: repository.ErrConflict}
	_, err := New(repo).Register(context.Background(), "User", "user@example.com", "secure-password")
	if !errors.Is(err, ErrEmailInUse) {
		t.Fatalf("Register() error = %v, want ErrEmailInUse", err)
	}
}
