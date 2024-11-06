package domains

import (
	"time"
	"context"
)

type Publisher struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type PublisherUsecase interface {
	Create(ctx context.Context, publisher *Publisher) error
	GetAll(ctx context.Context) ([]Publisher, error)
	GetByID(ctx context.Context, id uint) (*Publisher, error)
	Update(ctx context.Context, publisher *Publisher) error
	Delete(ctx context.Context, id uint) error
	CheckNameExists(ctx context.Context, name string, id uint) error
}

type PublisherRepository interface {
	Create(ctx context.Context, model interface{}) error
	GetAll(ctx context.Context, model interface{}, tmplScan interface{}, criteria string) error
	GetByID(ctx context.Context, model interface{}, tmplScan interface{}, selectedFields []string, order string, criteria string, id uint) (interface{}, error)
	Update(ctx context.Context, model interface{}, criteria string, args ...interface{}) error
	Delete(ctx context.Context, model interface{}, criteria string, id uint) error
	CheckNameExists(ctx context.Context, name string, id uint) error
}
