package usecase

import (
	"context"
	"errors"
	"log"
	"project-golang-crud/domains"
	"time"

	"github.com/google/uuid"
)

type bookRequestUsecase struct {
	repo domains.BookRequestRepository
}

func NewBookRequestUsecase(repo domains.BookRequestRepository) domains.BookRequestUsecase {
	return &bookRequestUsecase{repo: repo}
}

func (uc *bookRequestUsecase) GetBookByID(ctx context.Context, bookID uint) (*domains.Book, error) {
	return uc.repo.GetBookByID(ctx, bookID)
}

func (uc *bookRequestUsecase) Create(ctx context.Context, request *domains.BookRequest) error {
	// Cek apakah buku ada dan stok tersedia
	if err := uc.repo.CheckBookExists(ctx, int(request.BookID)); err != nil {
		return err
	}

	available, err := uc.repo.CheckStock(ctx, request.BookID)
	if err != nil {
		return err
	}
	if !available {
		return errors.New("Book out of stock")
	}

	// Cek apakah user ada
	if err := uc.repo.CheckUserExists(ctx, request.UserID); err != nil {
		return err
	}

	// Cek apakah UserID dan username sesuai
	if err := uc.repo.CheckUserIDByUsername(ctx, request.UserID, request.Username); err != nil {
		return err
	}

	// Set data dan buat permintaan peminjaman
	request.RequestDate = time.Now()
	request.Status = "DIPROSES"
	return uc.repo.Create(ctx, request)
}

func (uc *bookRequestUsecase) GetAll(ctx context.Context) ([]domains.BookRequest, error) {
	var requests []domains.BookRequest
	err := uc.repo.GetAll(ctx, &domains.BookRequest{}, &requests, "")
	return requests, err
}

func (uc *bookRequestUsecase) GetByUsernameRequest(ctx context.Context, username string) (*domains.UserWithRequest, error) {
	userWithLoans, err := uc.repo.GetByUsernameRequest(ctx, username)
	log.Printf(username)
	if err != nil {
		return nil, err
	}
	return userWithLoans, nil
}

// func (uc *bookReques) GetByUsernameRequest(ctx context.Context, username string) (*domains.UserWithRequest, error) {
// 	user, err := u.userRepo.GetByUsername(ctx, username)
// 	if err != nil {
// 		return nil, err
// 	}

// 	requests, err := u.bookRequestRepo.GetByUsernameRequest(ctx, username)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &domains.UserWithRequest{
// 		User:        *user,
// 		BookRequest: requests,
// 	}, nil
// }

func (uc *bookRequestUsecase) GetByID(ctx context.Context, id uint) (*domains.BookRequest, error) {
	var request domains.BookRequest
	data, err := uc.repo.GetByID(ctx, &domains.BookRequest{}, &request, []string{}, "", "", id)
	if err != nil {
		return nil, err
	}
	return data.(*domains.BookRequest), nil
}

func (uc *bookRequestUsecase) Update(ctx context.Context, request *domains.BookRequest) error {
	return uc.repo.Update(ctx, request, "id = ?", request.ID)
}

func (uc *bookRequestUsecase) Delete(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, &domains.BookRequest{}, "id = ?", id)
}

func (uc *bookRequestUsecase) ApproveOrReject(ctx context.Context, id uint, approve bool, reason string) (time.Time, error) {
	if approve {
		// Ambil permintaan buku untuk peminjaman
		request, err := uc.GetByID(ctx, id)
		if err != nil {
			return time.Time{}, err
		}

		// Periksa ketersediaan stok buku
		stockAvailable, err := uc.repo.CheckStock(ctx, request.BookID)
		if err != nil {
			return time.Time{}, err
		}

		// Pastikan stok lebih dari 0 sebelum menyetujui
		if !stockAvailable {
			return time.Time{}, errors.New("stok tidak tersedia untuk disetujui")
		}

		// Jika stok tersedia, lanjutkan dengan membuat catatan peminjaman
		loanRecord := &domains.LoanRecord{
			BookID:         request.BookID,
			UserID:         request.UserID,
			LoanDate:       time.Now(),
			DueDate:        request.RequestDate.AddDate(0, 0, 7), // Contoh: due date 7 hari setelah pinjam
			ReturnStatus:   false,
			BookRequestID:  request.ID, // Menyimpan ID permintaan buku
		}

		// Membuat catatan peminjaman
		if err := uc.CreateLoanRecord(ctx, loanRecord); err != nil {
			return time.Time{}, err
		}

		// Kurangi stok buku
		if err := uc.DecreaseStock(ctx, request.BookID); err != nil {
			return time.Time{}, err
		}

		// Ubah status permintaan menjadi disetujui
		if dueDate, err := uc.repo.ApproveOrReject(ctx, id, true, ""); err != nil {
			return time.Time{}, err
		} else {
			return dueDate, nil // Kembalikan due date
		}
	}

	// Jika ditolak, ubah status permintaan
	if _, err := uc.repo.ApproveOrReject(ctx, id, false, reason); err != nil {
		return time.Time{}, err
	}

	return time.Time{}, nil
}


func (uc *bookRequestUsecase) CheckStock(ctx context.Context, id uint) (bool, error) {
	return uc.repo.CheckStock(ctx, id)
}

func (uc *bookRequestUsecase) CheckBookExists(ctx context.Context, id int) error {
	return uc.repo.CheckBookExists(ctx, id)
}

func (uc *bookRequestUsecase) CheckUserExists(ctx context.Context, userID uuid.UUID) error {
	return uc.repo.CheckUserExists(ctx, userID)
}

func (uc *bookRequestUsecase) CreateLoanRecord(ctx context.Context, record *domains.LoanRecord) error {
	return uc.repo.CreateLoanRecord(ctx, record)
}

func (u *bookRequestUsecase) DecreaseStock(ctx context.Context, bookID uint) error {
	var book domains.Book
	// Mengambil data buku berdasarkan ID
	_, err := u.repo.GetByID(ctx, &book, &book, nil, "", "", bookID)
	if err != nil {
		return err
	}

	// Periksa apakah stok cukup
	if book.Stock <= 0 {
		return errors.New("stok buku tidak cukup")
	}

	// Kurangi stok
	book.Stock--
	return u.repo.Update(ctx, &book, "id = ?", book.ID)
}


func (uc *bookRequestUsecase) MoveToLoan(ctx context.Context, request *domains.BookRequest) error {
   // Logika untuk memindahkan permintaan buku ke dalam peminjaman
   loanRecord := &domains.LoanRecord{
	BookID:   request.BookID,
	UserID:   request.UserID,
	LoanDate: time.Now(),
	// Tambahkan field lain sesuai kebutuhan
}

if err := uc.CreateLoanRecord(ctx, loanRecord); err != nil {
	return err
}

// Update status permintaan jika perlu
return uc.Update(ctx, &domains.BookRequest{
	ID:     request.ID,
	Status: "DISETUJUI",
})
}