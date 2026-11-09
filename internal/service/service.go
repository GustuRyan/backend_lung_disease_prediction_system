package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/domain"
	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailInUse         = errors.New("email is already registered")
)

type Service struct {
	repository repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repository: repo}
}

func (s *Service) Users(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	return users, s.repository.List(ctx, &users, repository.Query{})
}

func (s *Service) Register(ctx context.Context, name, email, password string) (domain.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	var existing domain.User
	err := s.repository.First(ctx, &existing, repository.Query{Where: "email = ?", Args: []any{email}})
	if err == nil {
		return domain.User{}, ErrEmailInUse
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return domain.User{}, fmt.Errorf("check email: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, fmt.Errorf("hash password: %w", err)
	}
	user := domain.User{Name: strings.TrimSpace(name), Email: email, Password: string(hash), Role: "user"}
	if err := s.repository.Create(ctx, &user); err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return domain.User{}, ErrEmailInUse
		}
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (domain.User, error) {
	var user domain.User
	err := s.repository.First(ctx, &user, repository.Query{
		Where: "email = ?", Args: []any{strings.ToLower(strings.TrimSpace(email))},
	})
	if errors.Is(err, repository.ErrNotFound) {
		return domain.User{}, ErrInvalidCredentials
	}
	if err != nil {
		return domain.User{}, fmt.Errorf("find user: %w", err)
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return domain.User{}, ErrInvalidCredentials
	}
	return user, nil
}

func (s *Service) Diseases(ctx context.Context) ([]domain.Disease, error) {
	var items []domain.Disease
	return items, s.repository.List(ctx, &items, repository.Query{})
}

func (s *Service) DiseaseByName(ctx context.Context, name string) (domain.Disease, error) {
	var item domain.Disease
	err := s.repository.First(ctx, &item, repository.Query{Where: "disease_name = ?", Args: []any{name}})
	return item, err
}

func (s *Service) Datasets(ctx context.Context) ([]domain.Dataset, error) {
	var items []domain.Dataset
	return items, s.repository.List(ctx, &items, repository.Query{Preloads: []string{"User"}})
}

func (s *Service) Models(ctx context.Context, activeOnly bool) ([]domain.Model, error) {
	var items []domain.Model
	query := repository.Query{}
	if activeOnly {
		query.Where, query.Args = "active = ?", []any{true}
	}
	return items, s.repository.List(ctx, &items, query)
}

func (s *Service) Recommendations(ctx context.Context, diseaseID *uint) ([]domain.Recommendation, error) {
	var items []domain.Recommendation
	query := repository.Query{Preloads: []string{"Disease"}}
	if diseaseID != nil {
		query.Where, query.Args = "disease_id = ?", []any{*diseaseID}
	}
	return items, s.repository.List(ctx, &items, query)
}

func (s *Service) Histories(ctx context.Context, userID uint) ([]domain.History, error) {
	var items []domain.History
	query := repository.Query{Where: "user_id = ?", Args: []any{userID}, Preloads: []string{"User", "Disease"}}
	return items, s.repository.List(ctx, &items, query)
}

func (s *Service) DatasetModels(ctx context.Context) ([]domain.DatasetModel, error) {
	var items []domain.DatasetModel
	return items, s.repository.List(ctx, &items, repository.Query{Preloads: []string{"Dataset", "MLModel"}})
}

func (s *Service) Create(ctx context.Context, value any) error {
	return s.repository.Create(ctx, value)
}

func (s *Service) Update(ctx context.Context, model any, id uint, changes map[string]any) error {
	return s.repository.UpdateByID(ctx, model, id, changes)
}

func (s *Service) Delete(ctx context.Context, model any, id uint) error {
	return s.repository.DeleteByID(ctx, model, id)
}
