package domains

import (
	"time"
	"errors"
)

type User struct {
	ID        string `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	Username  string `gorm:"unique;not null" json:"username"`
	Email     string `gorm:"unique;not null" json:"email"`
	Password  string `gorm:"not null" json:"-"`
	CreatedAt time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}

type UserRepository interface{
	Create(user *User) error
	Update(user *User) error
	Delete(id string) error
	GetByUsername(username string) (*User, error)
	GetByID(id string) (*User, error) 
}

type UserUsecase interface{
	Register(username, email, password string) (*User,  error)
	Update(id string, username, email, password string)error
	Delete(id string) (*User, error)
	Validate(username, password string) (string, error)
	GetByUsername(username string) (*User, error)
	GetByID(id string) (*User, error) 
}
type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
    Errors []ErrorDetail `json:"errors,omitempty"`
	Code    int         `json:"code"`
}
type ErrorDetail struct {
    Message  string `json:"message"`
    Parameter string `json:"parameter"`
}

type DeleteRequest struct {
    ID string `json:"id"` // ID yang diterima dari request body
}
var ErrUserNotFound = errors.New("user not found")


