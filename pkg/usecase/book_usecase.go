package usecase

import "project-golang-crud/domains"

type bookUsecase struct {
	repo domains.BookRepository
}

func NewBookUsecase(repo domains.BookRepository) domains.BookUsecase {
	return &bookUsecase{repo}
}
func (u *bookUsecase) Create(book *domains.Book) error {
	return u.repo.Create(book)
}
func (u *bookUsecase) Update(book *domains.Book) error {
	return u.repo.Update(book)
}
func (u *bookUsecase) Delete(id uint) error {
	return u.repo.Delete(id)
}
func (u *bookUsecase) GetByID(id uint) (*domains.Book, error) {
	return u.repo.GetByID(id)
}
func (u *bookUsecase) GetAll() ([]domains.Book, error) {
	return u.repo.GetAll()
}
