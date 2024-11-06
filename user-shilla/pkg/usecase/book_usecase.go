package usecase

import (
	"context"
	"errors"
	"project-golang-crud/domains"
	"time"	
)

type bookUsecase struct {
	repo domains.BookRepository
}

func NewBookUsecase(repo domains.BookRepository) domains.BookUsecase {
	return &bookUsecase{
		repo: repo,
	}
}

func (u *bookUsecase) GetAuthorNameByID(ctx context.Context, authorID uint) (string, error) {
    author, err := u.repo.GetAuthorByID(ctx, authorID)
    if err != nil {
        return "", err
    }
    return author.Name, nil
}

func (u *bookUsecase) GetPublisherNameByID(ctx context.Context, publisherID uint) (string, error) {
    publisher, err := u.repo.GetPublisherByID(ctx, publisherID)
    if err != nil {
        return "", err
    }
    return publisher.Name, nil
}


// Create menambah buku baru
func (u *bookUsecase) Create(ctx context.Context, book *domains.Book) error {
	if book.Title == "" {
		return errors.New("Title is required")
	}

	// Cek apakah title sudah ada
	err := u.repo.CheckTitleExists(ctx, book.Title, book.ID)
	if err != nil {
		return err
	}

	// Cek apakah author ada
	if err := u.repo.AuthorExists(ctx, int(book.AuthorID)); err != nil {
		return err
	}

	// Cek apakah publisher ada
	if err := u.repo.PublisherExists(ctx, int(book.PublisherID)); err != nil {
		return err
	}

	// Mengambil nama author dan publisher
	authorName, err := u.repo.GetAuthorNameByID(ctx, book.AuthorID)
	if err != nil {
		return err
	}
	book.AuthorName = authorName

	publisherName, err := u.repo.GetPublisherNameByID(ctx, book.PublisherID)
	if err != nil {
		return err
	}
	book.PublisherName = publisherName

	// Menyimpan book
	err = u.repo.Create(ctx, book)
	if err != nil {
		return err
	}

	return nil
}


// GetAll mengambil semua buku
func (u *bookUsecase) GetAll(ctx context.Context) ([]domains.Book, error) {
	var tmplScan []domains.Book 
	err := u.repo.GetAll(ctx, &domains.Book{}, &tmplScan, "")
	if err != nil {
		return nil, err  
	}
	return tmplScan, nil
}

// GetByID mengambil buku berdasarkan ID
func (u *bookUsecase) GetByID(ctx context.Context, id uint) (*domains.Book, error) {
	var tmplScan domains.Book 
	
	_, err := u.repo.GetByID(ctx, &domains.Book{}, &tmplScan, nil, "", "", id)
	if err != nil {
		return nil, err
	}

	return &tmplScan, nil
}

// Update memperbarui buku
func (u *bookUsecase) Update(ctx context.Context, updatedBook *domains.Book) error {
    if updatedBook.Title == "" {
        return errors.New("Title is required")
    }

    var existingBook domains.Book
    _, err := u.repo.GetByID(ctx, &existingBook, &existingBook, nil, "", "", updatedBook.ID)
    if err != nil {
        return err
    }

    if err := u.repo.CheckTitleExists(ctx, updatedBook.Title, existingBook.ID); err != nil {
        return errors.New("Title already exists")
    }

    // Ambil dan tetapkan nama penulis
    authorName, err := u.GetAuthorNameByID(ctx, updatedBook.AuthorID) // Make sure AuthorID is of type uint
    if err != nil {
        return err
    }
    updatedBook.AuthorName = authorName

    // Ambil dan tetapkan nama penerbit
    publisherName, err := u.GetPublisherNameByID(ctx, updatedBook.PublisherID)
    if err != nil {
        return err
    }
    updatedBook.PublisherName = publisherName

    if err := u.repo.AuthorExists(ctx, int(updatedBook.AuthorID)); err != nil {
        return err
    }

    if err := u.repo.PublisherExists(ctx, int(updatedBook.PublisherID)); err != nil {
        return err
    }

    return u.repo.Update(ctx, updatedBook, "id = ?", updatedBook.ID)
}



// Delete menghapus buku secara soft delete
func (u *bookUsecase) Delete(ctx context.Context, id uint) error {
	var existingBook domains.Book
	_, err := u.repo.GetByID(ctx, &existingBook, &existingBook, nil, "", "", id)
	if err != nil {
		return err
	}

	if existingBook.DeletedAt != nil {
		return errors.New("Book already deleted")
	}

	deletedAt := time.Now()
	existingBook.DeletedAt = &deletedAt

	err = u.repo.Delete(ctx, existingBook, "id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

// AuthorExists memeriksa apakah pengarang ada
func (u *bookUsecase) AuthorExists(ctx context.Context, authorID int) error {
	return u.repo.AuthorExists(ctx, authorID)
}

// PublisherExists memeriksa apakah penerbit ada
func (u *bookUsecase) PublisherExists(ctx context.Context, publisherID int) error {
	return u.repo.PublisherExists(ctx, publisherID)
}

// CheckTitleExists memeriksa apakah judul buku ada
func (u *bookUsecase) CheckTitleExists(ctx context.Context, title string, id uint) error {
	return u.repo.CheckTitleExists(ctx, title, id)
}

// CheckStock memeriksa ketersediaan stok buku
func (u *bookUsecase) CheckStock(ctx context.Context, bookID uint) (bool, error) {
	var book domains.Book
	_, err := u.repo.GetByID(ctx, &book, &book, nil, "", "", bookID)
	if err != nil {
		return false, err
	}

	return book.Stock > 0, nil
}

// DecreaseStock mengurangi stok buku saat dipinjam
func (u *bookUsecase) DecreaseStock(ctx context.Context, bookID uint) error {
	var book domains.Book
	_, err := u.repo.GetByID(ctx, &book, &book, nil, "", "", bookID)
	if err != nil {
		return err
	}

	if book.Stock <= 0 {
		return errors.New("stok buku tidak cukup")
	}

	book.Stock-- // Kurangi stok
	return u.repo.Update(ctx, &book, "id = ?", book.ID)
}
