package domains

import (
	"time"
	"context"
)

type Author struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type AuthorUsecase interface {
	Create(ctx context.Context, author *Author) error
	GetAll(ctx context.Context) ([]Author, error)
	GetByID(ctx context.Context, id uint) (*Author, error)
	Update(ctx context.Context, author *Author) error
	Delete(ctx context.Context, id uint) error
	CheckNameExists(ctx context.Context, name string, id uint) error
}

type AuthorRepository interface {
	Create(ctx context.Context, model interface{}) error
	GetAll(ctx context.Context, model interface{}, tmplScan interface{}, criteria string) error
	GetByID(ctx context.Context, model interface{}, tmplScan interface{}, selectedFields []string, order string, criteria string, id uint) (interface{}, error)
	Update(ctx context.Context, model interface{}, criteria string, args ...interface{}) error
	Delete(ctx context.Context, model interface{}, criteria string, id uint) error
	CheckNameExists(ctx context.Context, name string, id uint) error
}