package domains

// Author struct represents the Author model
type Author struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name       string    `json:"name"`            // JSON tag harus diawali dengan json:
}

// AuthorRepository defines the methods that a repository implementation should provide
type AuthorRepository interface {
	Create(author *Author) error
	Update(author *Author) error
	Delete(id uint) error
	GetByID(id uint) (*Author, error)
	GetAll() ([]Author, error)
}

// BookUsecase defines the methods for the business logic layer (usecase)
type AuthorUsecase interface {
	Create(author *Author) error
	Update(author *Author) error
	Delete(id uint) error
	GetByID(id uint) (*Author, error)
	GetAll() ([]Author, error)
}
