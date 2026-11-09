package repository

import (
	"context"
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("record not found")
	ErrConflict = errors.New("record already exists")
)

type Query struct {
	Where    string
	Args     []any
	Preloads []string
}

type Repository interface {
	List(ctx context.Context, destination any, query Query) error
	First(ctx context.Context, destination any, query Query) error
	Create(ctx context.Context, value any) error
	UpdateByID(ctx context.Context, model any, id uint, changes map[string]any) error
	DeleteByID(ctx context.Context, model any, id uint) error
}

type Store struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) List(ctx context.Context, destination any, query Query) error {
	db := applyQuery(s.db.WithContext(ctx), query)
	return db.Find(destination).Error
}

func (s *Store) First(ctx context.Context, destination any, query Query) error {
	db := applyQuery(s.db.WithContext(ctx), query)
	if err := db.First(destination).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	} else {
		return err
	}
}

func (s *Store) Create(ctx context.Context, value any) error {
	return normalizeError(s.db.WithContext(ctx).Create(value).Error)
}

func (s *Store) UpdateByID(ctx context.Context, model any, id uint, changes map[string]any) error {
	result := s.db.WithContext(ctx).Model(model).Where("id = ?", id).Updates(changes)
	if result.Error != nil {
		return normalizeError(result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) DeleteByID(ctx context.Context, model any, id uint) error {
	result := s.db.WithContext(ctx).Where("id = ?", id).Delete(model)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func applyQuery(db *gorm.DB, query Query) *gorm.DB {
	if query.Where != "" {
		db = db.Where(query.Where, query.Args...)
	}
	for _, preload := range query.Preloads {
		db = db.Preload(preload)
	}
	return db
}

func normalizeError(err error) error {
	var mysqlError *mysql.MySQLError
	if errors.As(err, &mysqlError) && mysqlError.Number == 1062 {
		return ErrConflict
	}
	return err
}
