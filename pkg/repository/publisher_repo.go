package repository

import (
	"project-golang-crud/domains"

	"gorm.io/gorm"
)

type PublisherRepository struct {
	db *gorm.DB
}
func NewPublisherRepository(db *gorm.DB) domains.PublisherRepository {
	return &PublisherRepository{db}
}
func (r * PublisherRepository) Create(publisher *domains.Publisher) error {
	return r.db.Create(publisher).Error
}
func (r *PublisherRepository) Update(publisher *domains.Publisher) error {
	return r.db.Save(publisher).Error
}
func (r *PublisherRepository) Delete(id uint) error {
	return r.db.Delete(&domains.Publisher{}, id).Error
}
func (r *PublisherRepository) GetByID(id uint) (*domains.Publisher, error) {
	var publisher domains.Publisher
	err := r.db.First(&publisher, id).Error
	return &publisher, err
}
func (r *PublisherRepository) GetAll() ([]domains.Publisher, error) {
	var publishers []domains.Publisher
	err := r.db.Find(&publishers).Error
	return publishers, err
}