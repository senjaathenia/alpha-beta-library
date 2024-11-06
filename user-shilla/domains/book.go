package domains

import (
	"context"
	"time"
)

type Book struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Title          string    `gorm:"unique;not null" json:"title"`
	AuthorID       uint      `gorm:"not null" json:"author_id"`
	AuthorName     string    `json:"author_name"`     // Menambahkan field AuthorName
	PublisherID    uint      `gorm:"not null" json:"publisher_id"`
	PublisherName  string    `json:"publisher_name"`  // Menambahkan field PublisherName
	Summary        string    `gorm:"not null" json:"summary"` 
	Stock          int       `gorm:"not null" json:"stock"`
	MaxStock       int       `gorm:"not null" json:"max_stock"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      *time.Time `gorm:"index" json:"deleted_at"`
}

type BookUsecase interface {
	Create(ctx context.Context, book *Book) error
	Update(ctx context.Context, book *Book) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*Book, error)
	GetAll(ctx context.Context) ([]Book, error)
	AuthorExists(ctx context.Context, authorID int) error
	PublisherExists(ctx context.Context, publisherID int) error
	CheckTitleExists(ctx context.Context, title string, id uint) error
	GetAuthorNameByID(ctx context.Context, authorID uint) (string, error) // Mendapatkan nama author berdasarkan ID
	GetPublisherNameByID(ctx context.Context, publisherID uint) (string, error) // Mendapatkan nama publisher berdasarkan ID
}

type BookRepository interface {
	GetByID(ctx context.Context, model interface{}, tmplScan interface{}, selectedFields []string, order string, criteria string, id uint) (interface{}, error)
	GetAll(ctx context.Context, model interface{}, tmplScan interface{}, criteria string) error
	Create(ctx context.Context, model interface{}) error
	Update(ctx context.Context, model interface{}, criteria string, args ...interface{}) error
	Delete(ctx context.Context, model interface{}, criteria string, id uint) error
	AuthorExists(ctx context.Context, authorID int) error
	PublisherExists(ctx context.Context, publisherID int) error
	CheckTitleExists(ctx context.Context, title string, id uint) error
	DecreaseStock(ctx context.Context, id uint) error // Mengurangi stok buku
	CheckStock(ctx context.Context, id uint) (bool, error) // Memeriksa ketersediaan stok buku
	GetAuthorNameByID(ctx context.Context, authorID uint) (string, error) // Mendapatkan nama author berdasarkan ID
	GetPublisherNameByID(ctx context.Context, publisherID uint) (string, error) // Mendapatkan nama publisher berdasarkan ID
    GetAuthorByID(ctx context.Context, authorID uint) (Author, error)
    GetPublisherByID(ctx context.Context, publisherID uint) (Publisher, error)

}
