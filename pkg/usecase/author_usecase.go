package usecase

import "project-golang-crud/domains"

type authorUsecase struct {
	repo domains.AuthorRepository
}

func NewAuthorUsecase(repo domains.AuthorRepository) domains.AuthorUsecase {
	return &authorUsecase{repo}
}
func (u *authorUsecase) Create(author *domains.Author)error {
	return u.repo.Create(author)
}
func (u *authorUsecase) Update(author *domains.Author) error{
	return u.repo.Update(author)
}
func (u *authorUsecase) Delete(id uint) error {
	return u.repo.Delete(id)
}
func (u *authorUsecase) GetByID(id uint) (*domains.Author, error) {
	return u.repo.GetByID(id)
}
func (u *authorUsecase) GetAll() ([]domains.Author, error) {
	return u.repo.GetAll()
}