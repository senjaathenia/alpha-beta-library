package repository

import (
	"project-golang-crud/domains"

	"gorm.io/gorm"
)

type AuthorsRepository struct {
	db *gorm.DB
}
func NewAuthorRepository(db *gorm.DB) domains.AuthorRepository {
	return &AuthorsRepository{db}
}
func (r *AuthorsRepository) Create(author *domains.Author) error {
	return r.db.Create(author).Error
}
func (r *AuthorsRepository) Update(author *domains.Author) error {
	return r.db.Save(author).Error
}
func (r *AuthorsRepository) Delete(id uint) error {
	return r.db.Delete(&domains.Author{}, id).Error
}
func (r *AuthorsRepository) GetByID(id uint) (*domains.Author, error) {
	var author domains.Author
	err := r.db.First(&author, id).Error
	return &author, err
}
func (r *AuthorsRepository) GetAll() ([]domains.Author, error) {
	var authors []domains.Author
	err := r.db.Find(&authors).Error
	return authors, err
}