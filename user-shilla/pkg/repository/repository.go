package repository

import (
	"context"
	"errors"
	"fmt"
	"project-golang-crud/domains"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GenericRepository struct {
	db *gorm.DB
}

func NewGenericRepository(db *gorm.DB) *GenericRepository {
	return &GenericRepository{db: db}
}

func (r *GenericRepository) GetByID(ctx context.Context, model interface{}, tmplScan interface{}, selectedFields []string, order string, criteria string, id uint) (interface{}, error) {
	query := r.db.WithContext(ctx).Model(model)

	// //Selected
	if len(selectedFields) > 0 {
		query = query.Select(selectedFields)
	}

	// //Criteria
	if criteria != "" {
		query = query.Where(criteria, id)
	} else {
		query = query.Where("id = ?", id)
	}

	// //Order
	if order != "" {
		query = query.Order(order)
	}

	err := query.First(tmplScan).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return tmplScan, nil
}

func (r *GenericRepository) GetAuthorByID(ctx context.Context, id uint) (domains.Author, error) {
    var author domains.Author
    err := r.db.WithContext(ctx).First(&author, id).Error
    return author, err
}

func (r *GenericRepository) GetPublisherByID(ctx context.Context, id uint) (domains.Publisher, error) {
    var publisher domains.Publisher
    err := r.db.WithContext(ctx).First(&publisher, id).Error
    return publisher, err
}

func (r *GenericRepository) GetAll(ctx context.Context, model interface{}, tmplScan interface{}, criteria string) error {
	err := r.db.WithContext(ctx).Model(model).Where("deleted_at is NULL").Where(criteria).Find(tmplScan).Error
	return err
}

func (r *GenericRepository) Create(ctx context.Context, model interface{}) error {
	err := r.db.WithContext(ctx).Create(model).Error

	if err != nil {
		if strings.Contains(err.Error(), "Duplicate key value violates unique constraint") {
			return fmt.Errorf("Record already exist: %v", err)
		}
		return err
	}
	return nil
}

func (r *GenericRepository) GetAuthorNameByID(ctx context.Context, authorID uint) (string, error) {
	var author domains.Author
	if err := r.db.WithContext(ctx).Where("id = ?", authorID).First(&author).Error; err != nil {
		return "", err
	}
	return author.Name, nil
}

func (r *GenericRepository) GetPublisherNameByID(ctx context.Context, publisherID uint) (string, error) {
	var publisher domains.Publisher
	if err := r.db.WithContext(ctx).Where("id = ?", publisherID).First(&publisher).Error; err != nil {
		return "", err
	}
	return publisher.Name, nil
}


func (r *GenericRepository) Update(ctx context.Context, model interface{}, criteria string, args ...interface{}) error {
    // Dapatkan tipe dan nilai dari model
    val := reflect.ValueOf(model).Elem()
    typ := val.Type()

    // Ambil ID dari model untuk menemukan record yang akan diupdate
    idField := val.FieldByName("ID")
    if !idField.IsValid() {
        return errors.New("Model must have an ID field")
    }

    var existingModel interface{}
    // Buat instance dari model yang sesuai
    switch typ.Name() {
    case "Book":
        existingModel = &domains.Book{}
    case "Publisher":
        existingModel = &domains.Publisher{}
    case "Author":
        existingModel = &domains.Author{}
    case "Loan":
        existingModel = &domains.BookLoans{}
    case "Request":
        existingModel = &domains.BookRequest{}
    default:
        return errors.New("Unsupported model type")
    }

    // Ambil record yang ada berdasarkan ID
    err := r.db.WithContext(ctx).Where("id = ?", idField.Interface()).First(existingModel).Error
    if err != nil {
        return gorm.ErrRecordNotFound
    }

    // Set UpdatedAt ke waktu saat ini
    updatedAtField := val.FieldByName("UpdatedAt")
    if updatedAtField.IsValid() {
        updatedAtField.Set(reflect.ValueOf(time.Now()))
    }

    // Siapkan query untuk memperbarui
    query := r.db.WithContext(ctx).Model(existingModel)

    // Tambahkan kondisi jika ada
    if criteria != "" {
        query = query.Where(criteria, args...)
    }

    // Loop untuk mengecek dan update field yang disertakan
    for i := 0; i < typ.NumField(); i++ {
        field := val.Field(i)
        if field.CanInterface() && !field.IsZero() { // Cek apakah field bukan nol
            query = query.Set(typ.Field(i).Name, field.Interface())
        }
    }

    // Gunakan Omit untuk tidak mengubah CreatedAt
    err = query.Omit("created_at").Updates(existingModel).Error
    if err != nil {
        return err
    }

    return nil
}


func (r *GenericRepository) Delete(ctx context.Context, model interface{}, criteria string, id uint) error {
	err := r.db.WithContext(ctx).Model(model).Where("id = ?", id).Updates(map[string]interface{}{"deleted_at": time.Now().UTC()}).Error
	return err
}

func (r *GenericRepository) AuthorExists(ctx context.Context, authorID int) error {
	var author domains.Author
	err := r.db.Model(&domains.Author{}).Where("id = ?", authorID).First(&author).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Author does not exist")
		}
		return err
	}
	return nil
}

func (r *GenericRepository) PublisherExists(ctx context.Context, publisherID int) error {
	var publisher domains.Publisher
	err := r.db.Model(&domains.Publisher{}).Where("id = ?", publisherID).First(&publisher).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Publisher does not exist")
		}
		return err
	}
	return nil
}
func (r *GenericRepository) CheckTitleExists(ctx context.Context, title string, id uint) error {
	var book domains.Book
	// Mencari buku dengan judul yang sama, kecuali yang memiliki ID tertentu
	err := r.db.WithContext(ctx).Model(&domains.Book{}).Where("title = ? AND id != ?", title, id).First(&book).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Jika tidak ditemukan buku dengan judul yang sama, kembalikan nil
			return nil
		}
		// Kembalikan error lain jika ada
		return err
	}
	// Jika buku dengan judul yang sama ditemukan, kembalikan error
	return errors.New("book with this title already exists")
}

func (r *GenericRepository) CheckNameExists(ctx context.Context, name string, id uint) error {
	var author domains.Author
	// Mencari penulis dengan nama yang sama, kecuali yang memiliki ID tertentu
	err := r.db.WithContext(ctx).Model(&domains.Author{}).Where("name = ? AND id != ?", name, id).First(&author).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Jika tidak ditemukan penulis dengan nama yang sama, kembalikan nil
			return nil
		}
		// Kembalikan error lain jika ada
		return err
	}
	// Jika pen	ulis dengan nama yang sama ditemukan, kembalikan error
	return errors.New("author with this name already exists")
}

// Menambahkan metode di repository
func (r *GenericRepository) CheckBookExists(ctx context.Context, bookID int) error {
	var book domains.Book
	err := r.db.WithContext(ctx).Where("id = ?", bookID).First(&book).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Book does not exist")
		}
		return err
	}
	return nil
}

func (r *GenericRepository) CheckUserExists(ctx context.Context, userID uuid.UUID) error {
	var user domains.User // Pastikan ada struct User di domains
	err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("User does not exist")
		}
		return err
	}
	return nil
}

func (r *GenericRepository) GetByUsername(ctx context.Context, username string) (*domains.User, error) {
	var user domains.User
	err := r.db.WithContext(ctx).Where("username = ? AND deleted_at IS NULL", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Tidak ditemukan, return nil
		}
		return nil, err // Kembalikan error lain jika terjadi
	}
	return &user, nil
}

func (r *GenericRepository) CheckUserIDByUsername(ctx context.Context, userID uuid.UUID, username string) error {
	var user domains.User
	err := r.db.WithContext(ctx).Where("id = ? AND username = ?", userID, username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("User ID and username do not match")
		}
		return err
	}
	return nil
}

func (r *GenericRepository) GetByUsernameLoans(ctx context.Context, username string) (*domains.UserWithLoans, error) {
    var userWithLoans domains.UserWithLoans

    // Mengambil user berdasarkan username
    err := r.db.WithContext(ctx).
        Preload("BookLoans"). // Load related BookLoans records
        Where("username = ? AND deleted_at IS NULL", username).
        First(&userWithLoans.User).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, domains.ErrUserNotFound // Mengembalikan error khusus untuk user tidak ditemukan
        }
        return nil, err // Kembalikan error lain jika terjadi
    }

    // Memeriksa apakah user ditemukan
    if userWithLoans.User.ID == uuid.Nil {
        return nil, domains.ErrUserNotFound // Atau return nil, nil untuk tidak ditemukan
    }

    // Mengambil BookLoans berdasarkan user_id
    err = r.db.WithContext(ctx).
        Where("user_id = ?", userWithLoans.User.ID). // Mengambil BookLoans berdasarkan user_id
        Find(&userWithLoans.BookLoans).Error
    if err != nil {
        return nil, err // Kembalikan error jika terjadi
    }

    return &userWithLoans, nil
}

func (r *GenericRepository) GetByUsernameRequest(ctx context.Context, username string) (*domains.UserWithRequest, error) {
    var UserWithRequest domains.UserWithRequest

    // Mengambil user berdasarkan username
    err := r.db.WithContext(ctx).
        Preload("BookRequests"). // Load related BookLoans records
        Where("username = ? AND deleted_at IS NULL", username).
        First(&UserWithRequest.User).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, domains.ErrUserNotFound // Mengembalikan error khusus untuk user tidak ditemukan
        }
        return nil, err // Kembalikan error lain jika terjadi
    }

    // Memeriksa apakah user ditemukan
    if UserWithRequest.User.ID == uuid.Nil {
        return nil, domains.ErrUserNotFound // Atau return nil, nil untuk tidak ditemukan
    }

    // Mengambil BookLoans berdasarkan user_id
    err = r.db.WithContext(ctx).
        Where("user_id = ?", UserWithRequest.User.ID). // Mengambil BookLoans berdasarkan user_id
        Find(&UserWithRequest.BookRequest).Error
    if err != nil {
        return nil, err // Kembalikan error jika terjadi
    }

    return &UserWithRequest, nil
}



// CheckStock memeriksa apakah stok buku masih tersedia
// func (r *GenericRepository) CheckStock(ctx context.Context, bookID uint) (bool, error) {
// 	var count int64
// 	err := r.db.WithContext(ctx).Model(&domains.Book{}).
// 		Where("id = ? AND stock > 0", bookID).
// 		Count(&count).Error
// 	if err != nil {
// 		return false, err
// 	}
// 	return count > 0, nil
// }

func (r *GenericRepository) CheckStock(ctx context.Context, bookID uint) (bool, error) {
    var book domains.Book
    // Mencari buku berdasarkan ID
    err := r.db.WithContext(ctx).Where("id = ?", bookID).First(&book).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return false, errors.New("book does not exist")
        }
        return false, err
    }
    // Kembalikan true jika stok lebih dari 0, sebaliknya false
    return book.Stock > 0, nil
}


// DecreaseStock mengurangi jumlah stok buku sebesar satu
// DecreaseStock mengurangi jumlah stok buku sebesar satu
func (r *GenericRepository) DecreaseStock(ctx context.Context, bookID uint) error {
    var book domains.Book
    // Ambil buku berdasarkan ID
    err := r.db.WithContext(ctx).Where("id = ?", bookID).First(&book).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return errors.New("book does not exist")
        }
        return err
    }
    
    // Periksa apakah stok masih tersedia
    if book.Stock <= 0 {
        return errors.New("stock not available")
    }

    // Mengurangi stok jika masih ada
    return r.db.WithContext(ctx).Model(&domains.Book{}).
        Where("id = ?", bookID).
        Update("stock", gorm.Expr("stock - ?", 1)).Error
}


func (r *GenericRepository) ApproveOrReject(ctx context.Context, id uint, approve bool, reason string) (time.Time, error) {
	var request domains.BookRequest
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&request).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return time.Time{}, errors.New("loan request not found")
		}
		return time.Time{}, err
	}

	if request.Status != "DIPROSES" {
		return time.Time{}, errors.New("request cannot be processed because it is not in processing status")
	}

	loc, err := time.LoadLocation("Asia/Jakarta") // Set to the desired time zone
	if err != nil {
		return time.Time{}, err
	}

	if approve {
		if err := r.DecreaseStock(ctx, request.BookID); err != nil {
			return time.Time{}, err // Kembalikan error jika stok tidak tersedia
		}
		request.Status = "DISETUJUI"
		request.RequestDate = time.Now().In(loc).Add(3 * 24 * time.Hour) // Set due date in local timezone
	} else {
		request.Status = "DITOLAK"
		request.Reason = reason
	}

	if err := r.db.WithContext(ctx).Save(&request).Error; err != nil {
		return time.Time{}, err
	}

	return request.RequestDate, nil // Return the due date
}


func (r *GenericRepository) GetBookByID(ctx context.Context, bookID uint) (*domains.Book, error) {
    var book domains.Book
    err := r.db.WithContext(ctx).First(&book, bookID).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("book with ID %d not found", bookID) // Ganti error dengan pesan lebih spesifik
        }
        return nil, err
    }
    return &book, nil
}

func (r *GenericRepository) CreateLoanRecord(ctx context.Context, record *domains.LoanRecord) error {
	return r.db.WithContext(ctx).Table("book_loans").Create(record).Error
}


func (r *GenericRepository) CreateLoan(ctx context.Context, request *domains.BookRequest) error {
    // Misalkan Anda memiliki tabel LoanRecord dan bukan BookRequest untuk penyimpanan
    loan := &domains.LoanRecord{
        BookRequestID: request.ID,
        BookID:        request.BookID,
        UserID:        request.UserID,
        DueDate:       time.Now().Add(1 * 24 * time.Hour), // Contoh: dua minggu
    }
    return r.db.Create(loan).Error
}

func (r *GenericRepository) IncreaseBookStock(ctx context.Context, bookID int) error {
    return r.db.WithContext(ctx).
        Model(&domains.Book{}).
        Where("id = ?", bookID).
        Update("stock", gorm.Expr("stock + ?", 1)).Error
}

// ReturnBook mengembalikan buku dan memperbarui stok di database
func (r *GenericRepository) Return(ctx context.Context, loanID uint, returnDate time.Time, lateFee int) error {
    var loan domains.BookLoans

    // Mengambil record peminjaman berdasarkan loanID
    err := r.db.WithContext(ctx).First(&loan, loanID).Error
    if err != nil {
        return err
    }

    // Pastikan status pengembalian belum di-set
    if loan.ReturnStatus {
        return errors.New("book already returned")
    }

    // Update status, tanggal pengembalian, dan denda
    loan.ReturnDate = &returnDate
    loan.ReturnStatus = true
    loan.LateFee = lateFee

    // Simpan perubahan ke database
    if err := r.db.WithContext(ctx).Save(&loan).Error; err != nil {
        return err
    }

    return nil
}


// Contoh fungsi untuk memperbarui peminjaman buku dengan denda
func (r *GenericRepository) UpdateLoanWithFine(ctx context.Context, loanID int, denda int) error {
    return r.db.WithContext(ctx).
        Model(&domains.BookLoans{}).
        Where("id = ?", loanID).
        Update("denda", denda).Error
}

func (r *GenericRepository) AddLoan(ctx context.Context, loan *domains.BookLoans) error {
	return r.db.WithContext(ctx).Create(loan).Error
}