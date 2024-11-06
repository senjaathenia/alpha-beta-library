package usecase

import (
	"context"
	"errors"
	"time"
	"project-golang-crud/domains"
)

type authorUsecase struct {
	repo domains.AuthorRepository
}

func NewAuthorUsecase(repo domains.AuthorRepository) domains.AuthorUsecase {
	return &authorUsecase{
		repo: repo,
	}
}

func (u *authorUsecase) Create(ctx context.Context, author *domains.Author) error {
	if author.Name == "" {
		return errors.New("Name is required")
	}

	err := u.repo.CheckNameExists(ctx, author.Name, author.ID)
	if err != nil {
		return err
	}

	err = u.repo.Create(ctx, author)
	if err != nil {
		return err
	}

	return nil
}

func (u *authorUsecase) GetAll(ctx context.Context) ([]domains.Author, error) {
	var tmplScan []domains.Author 
	err := u.repo.GetAll(ctx, &domains.Author{}, &tmplScan, "")
	if err != nil {
		return nil, err  
	}
	return tmplScan, nil
}

func (u *authorUsecase) GetByID(ctx context.Context, id uint) (*domains.Author, error) {
    var tmplScan domains.Author 
	
    // Panggil repository dengan parameter lengkap
    _, err := u.repo.GetByID(ctx, &domains.Author{}, &tmplScan, nil, "", "", id)
    if err != nil {
        return nil, err
    }

    // Cast hasil ke *domains.Book dan kembalikan
    return &tmplScan, nil
}

func (u *authorUsecase) Update(ctx context.Context, updatedAuthor *domains.Author) error {
    if updatedAuthor.Name == "" {
        return errors.New("Name is required")
    }

    var existingAuthor domains.Author
    _, err := u.repo.GetByID(ctx, &existingAuthor, &existingAuthor, nil, "", "", updatedAuthor.ID)
    if err != nil {
        return err
    }

  
    if err := u.repo.CheckNameExists(ctx, updatedAuthor.Name, existingAuthor.ID); err != nil {
        return errors.New("Title already exists")
    }


    return u.repo.Update(ctx, updatedAuthor, "id = ?", existingAuthor.ID)
}

func (u *authorUsecase) Delete(ctx context.Context, id uint) error {
	var existingAuthor domains.Author
	_, err := u.repo.GetByID(ctx, &existingAuthor, &existingAuthor, nil, "", "", id)
	if err != nil {
		return err
	}

	if existingAuthor.DeletedAt != nil {
		return errors.New("Book already deleted")
	}
	

	deletedAt := time.Now()
	existingAuthor.DeletedAt = &deletedAt

	err = u.repo.Delete(ctx, existingAuthor, "id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (u *authorUsecase) CheckNameExists(ctx context.Context, name string, id uint) error {
	err := u.repo.CheckNameExists(ctx, name, id)
	if err != nil {
		return err
	}
	return nil
}
