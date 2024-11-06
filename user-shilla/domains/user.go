package domains

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID     `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	Username  string     `gorm:"unique;not null" json:"username"`
	Email     string     `gorm:"unique;not null" json:"email"`
	Password  string     `gorm:"not null" json:"-"`
	Role string 		`gorm:"not null" json:"role"`
	CreatedAt time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
	BookLoans []BookLoans  `gorm:"foreignKey:UserID"` // Pastikan ini ada jika kamu menggunakan relasi
	BookRequests []BookRequest `json:"book_requests" gorm:"foreignKey:UserID"` // Relasi dengan BookRequest}
}
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetAll(ctx context.Context) ([]Book, error)
	CreateBookRequest(ctx context.Context, user *User) error
}

type UserUsecase interface {
	Register(ctx context.Context, username, email, password, role string) (*User, error)
	Update(ctx context.Context, id string, username, email, password string) error
	Delete(ctx context.Context, id string) (*User, error)
	Validate(ctx context.Context, username, password string) (string, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetAll(ctx context.Context) ([]Book, error)
	CreateBookRequest(ctx context.Context, request *BookRequest) error
}


type Response struct {
	Message string        `json:"message"`
	Data    interface{}   `json:"data"`
	Errors  []ErrorDetail `json:"errors"`
	Code    int           `json:"code"`
}

type ErrorDetail struct {
	Message   string `json:"message"`
	Parameter string `json:"parameter"`
}

type DeleteRequest struct {
	ID string `json:"id"` // ID yang diterima dari request body
}

var ErrUserNotFound = errors.New("user not found")
