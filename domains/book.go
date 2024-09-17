package domains

import "time"

// Book struct represents the Book model
type Book struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`            // JSON tag harus diawali dengan json:
	AuthorID    uint      `json:"author_id"`        // Perbaikan tag JSON
	PublisherID uint      `json:"publisher_id"`     // Perbaikan tag JSON
	CreatedAt   time.Time `json:"created_at"`       // Perbaikan tag JSON
	UpdatedAt   time.Time `json:"updated_at"`       // Perbaikan tag JSON
}

// BookRepository defines the methods that a repository implementation should provide
type BookRepository interface {
	Create(book *Book) error
	Update(book *Book) error
	Delete(id uint) error
	GetByID(id uint) (*Book, error)
	GetAll() ([]Book, error)
}

// BookUsecase defines the methods for the business logic layer (usecase)
type BookUsecase interface {
	Create(book *Book) error
	Update(book *Book) error
	Delete(id uint) error
	GetByID(id uint) (*Book, error)
	GetAll() ([]Book, error)
}
