package usecase

import "project-golang-crud/domains"

type publisherUsecase struct {
	repo domains.PublisherRepository
}

func NewPublisherUsecase(repo domains.PublisherRepository) domains.PublisherUsecase {
	return &publisherUsecase{repo}
}
func (u *publisherUsecase) Create(publisher *domains.Publisher) error {
	return u.repo.Create(publisher)
}
func (u *publisherUsecase) Update(publisher *domains.Publisher) error {
	return u.repo.Update(publisher)
}
func (u *publisherUsecase) Delete(id uint) error {
	return u.repo.Delete(id)
}
func (u *publisherUsecase) GetByID(id uint) (*domains.Publisher, error) {
	return u.repo.GetByID(id)
}
func (u *publisherUsecase) GetAll() ([]domains.Publisher, error) {
	return u.repo.GetAll()
}