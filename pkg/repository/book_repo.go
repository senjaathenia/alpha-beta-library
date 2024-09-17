package repository

import (
	"project-golang-crud/domains"

	"gorm.io/gorm"
)

type BooksRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) domains.BookRepository {
	return &BooksRepository{db}
}
func (r *BooksRepository) Create(book *domains.Book) error {
	return r.db.Create(book).Error
}
func (r *BooksRepository) Update(book *domains.Book) error {
	return r.db.Save(book).Error
}
func (r *BooksRepository) Delete(id uint) error {
	return r.db.Delete(&domains.Book{}, id).Error
}
func (r *BooksRepository) GetByID(id uint) (*domains.Book, error) {
	var book domains.Book
	err := r.db.First(&book, id).Error
	return &book, err
}
func (r *BooksRepository) GetAll() ([]domains.Book, error) {
	var books []domains.Book
	err := r.db.Find(&books).Error
	return books, err
}
