package usecase

import (
	"context"
	"errors"
	"time"
	"project-golang-crud/domains"
)

type publisherUsecase struct {
	repo domains.PublisherRepository
}

func NewPublisherUsecase(repo domains.PublisherRepository) domains.PublisherUsecase {
	return &publisherUsecase{
		repo: repo,
	}
}

func (u *publisherUsecase) Create(ctx context.Context, publisher *domains.Publisher) error {
	if publisher.Name == "" {
		return errors.New("Name is required")
	}

	err := u.repo.CheckNameExists(ctx, publisher.Name, publisher.ID)
	if err != nil {
		return err
	}

	err = u.repo.Create(ctx, publisher)
	if err != nil {
		return err
	}

	return nil
}

func (u *publisherUsecase) GetAll(ctx context.Context) ([]domains.Publisher, error) {
	var tmplScan []domains.Publisher 
	err := u.repo.GetAll(ctx, &domains.Publisher{}, &tmplScan, "")
	if err != nil {
		return nil, err  
	}
	return tmplScan, nil
}

func (u *publisherUsecase) GetByID(ctx context.Context, id uint) (*domains.Publisher, error) {
    var tmplScan domains.Publisher 
	
    // Panggil repository dengan parameter lengkap
    _, err := u.repo.GetByID(ctx, &domains.Publisher{}, &tmplScan, nil, "", "", id)
    if err != nil {
        return nil, err
    }

    // Cast hasil ke *domains.Book dan kembalikan
    return &tmplScan, nil
}

func (u *publisherUsecase) Update(ctx context.Context, updatedPublisher *domains.Publisher) error {
    if updatedPublisher.Name == "" {
        return errors.New("Name is required")
    }

    var existingPublisher domains.Publisher
    _, err := u.repo.GetByID(ctx, &existingPublisher, &existingPublisher, nil, "", "", updatedPublisher.ID)
    if err != nil {
        return err
    }

  
    if err := u.repo.CheckNameExists(ctx, updatedPublisher.Name, existingPublisher.ID); err != nil {
        return errors.New("Name already exists")
    }


    return u.repo.Update(ctx, updatedPublisher, "id = ?", existingPublisher.ID)
}

func (u *publisherUsecase) Delete(ctx context.Context, id uint) error {
	var existingPublisher domains.Publisher
	_, err := u.repo.GetByID(ctx, &existingPublisher, &existingPublisher, nil, "", "", id)
	if err != nil {
		return err
	}

	if existingPublisher.DeletedAt != nil {
		return errors.New("Publisher already deleted")
	}
	

	deletedAt := time.Now()
	existingPublisher.DeletedAt = &deletedAt

	err = u.repo.Delete(ctx, existingPublisher, "id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (u *publisherUsecase) CheckNameExists(ctx context.Context, name string, id uint) error {
	err := u.repo.CheckNameExists(ctx, name, id)
	if err != nil {
		return err
	}
	return nil
}
