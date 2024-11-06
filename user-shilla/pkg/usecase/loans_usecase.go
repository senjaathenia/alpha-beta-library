package usecase

import (
	"context"
	"errors"
	"project-golang-crud/domains"
	"time"

	"fmt"

	"github.com/google/uuid"
)

type loanUsecase struct {
	repo domains.BookLoansRepository
}

func NewLoanUsecase(repo domains.BookLoansRepository) domains.LoanUsecase {
	return &loanUsecase{
		repo: repo,
	}
}

func (u *loanUsecase) Create(ctx context.Context, loan *domains.BookLoans) error {
	if loan.BookID == 0 {
		return errors.New("BookID is required")
	}

	if loan.UserID == uuid.Nil {
		return errors.New("UserID is required")
	}

	// Cek apakah buku dan pengguna ada
	if err := u.repo.CheckBookExists(ctx, loan.BookID); err != nil {
		return err
	}

	if err := u.repo.CheckUserExists(ctx, loan.UserID); err != nil {
		return err
	}

	// Validasi tanggal pinjam dan jatuh tempo
	if loan.LoanDate.IsZero() || loan.DueDate.IsZero() {
		return errors.New("LoanDate and DueDate are required")
	}

	err := u.repo.Create(ctx, loan)
	if err != nil {
		return err
	}

	return nil
}

func (u *loanUsecase) GetAll(ctx context.Context) ([]domains.BookLoans, error) {
	var tmplScan []domains.BookLoans
	err := u.repo.GetAll(ctx, &domains.BookLoans{}, &tmplScan, "")
	if err != nil {
		return nil, err
	}
	return tmplScan, nil
}

func (u *loanUsecase) GetByID(ctx context.Context, id uint) (*domains.BookLoans, error) {
	var tmplScan domains.BookLoans

	// Panggil repository dengan parameter lengkap
	_, err := u.repo.GetByID(ctx, &domains.BookLoans{}, &tmplScan, nil, "", "", id)
	if err != nil {
		return nil, err
	}

	return &tmplScan, nil
}

func (u *loanUsecase) Update(ctx context.Context, updatedLoan *domains.BookLoans) error {
	// Cek apakah pinjaman ada
	var existingLoan domains.BookLoans
	_, err := u.repo.GetByID(ctx, &existingLoan, &existingLoan, nil, "", "", uint(updatedLoan.ID))
	if err != nil {
		return err
	}

	// Validasi LoanDate dan DueDate jika perlu
	if updatedLoan.LoanDate.IsZero() || updatedLoan.DueDate.IsZero() {
		return errors.New("LoanDate and DueDate are required")
	}

	// Update logika sesuai kebutuhan
	return u.repo.Update(ctx, updatedLoan, "id = ?", existingLoan.ID)
}

func (u *loanUsecase) Delete(ctx context.Context, id uint) error {
	var existingLoan domains.BookLoans
	_, err := u.repo.GetByID(ctx, &existingLoan, &existingLoan, nil, "", "", id)
	if err != nil {
		return err
	}

	if existingLoan.DeletedAt != nil {
		return errors.New("Loan already deleted")
	}

	deletedAt := time.Now()
	existingLoan.DeletedAt = &deletedAt

	err = u.repo.Delete(ctx, existingLoan, "id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (u *loanUsecase) CheckBookExists(ctx context.Context, bookID int) error {
	err := u.repo.CheckBookExists(ctx, bookID)
	if err != nil {
		return err
	}
	return nil
}

func (u *loanUsecase) CheckUserExists(ctx context.Context, userID uuid.UUID) error {
	err := u.repo.CheckUserExists(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}

// CreateLoanRecord menyimpan catatan peminjaman
func (u *loanUsecase) CreateLoanRecord(ctx context.Context, loan *domains.LoanRecord) error {
	if err := u.repo.CreateLoanRecord(ctx, loan); err != nil {
		return err
	}
	return nil
}

func (u *loanUsecase) Return(ctx context.Context, loanID uint, returnDate time.Time) error {
    // Ambil data peminjaman
    loan, err := u.repo.GetByID(ctx, &domains.BookLoans{}, &domains.BookLoans{}, nil, "", "", loanID)
    if err != nil {
        return err
    }

    bookLoan, ok := loan.(*domains.BookLoans)
    if !ok {
        return errors.New("failed to assert loan to BookLoans")
    }

    if bookLoan.ReturnStatus {
        return errors.New("book already returned")
    }

    // Hitung denda jika tanggal pengembalian setelah tanggal jatuh tempo
    const lateFeeRate = 5000 // Denda per hari
    var lateFee int

    // Hitung denda
    if returnDate.After(bookLoan.DueDate) {
        lateDays := int(returnDate.Sub(bookLoan.DueDate).Hours() / 24)
        lateFee = lateDays * lateFeeRate
    }

    // Set tanggal pengembalian dan status pengembalian
    bookLoan.ReturnStatus = true
    bookLoan.ReturnDate = &returnDate // Set tanggal pengembalian
    bookLoan.LateFee = lateFee         // Simpan denda ke dalam data peminjaman

    // Panggil repository untuk memperbarui data peminjaman
    if err := u.repo.Return(ctx, uint(bookLoan.ID), returnDate, lateFee); err != nil {
        return err
    }

    // Panggil fungsi untuk menambah stok buku
    if err := u.repo.IncreaseBookStock(ctx, bookLoan.BookID); err != nil {
        return fmt.Errorf("failed to increase book stock: %v", err)
    }

    return nil
}

func (u *loanUsecase) GetByUsernameLoans(ctx context.Context, username string) (*domains.UserWithLoans, error) {
	userWithLoans, err := u.repo.GetByUsernameLoans(ctx, username)
	if err != nil {
		return nil, err
	}
	return userWithLoans, nil
}


func (u *loanUsecase) CalculateFine(ctx context.Context, loanID uint) (int, error) {
	loan, err := u.repo.GetByID(ctx, &domains.BookLoans{}, &domains.BookLoans{}, nil, "", "", loanID)
	if err != nil {
		return 0, err
	}

	bookLoan, ok := loan.(*domains.BookLoans)
	if !ok {
		return 0, errors.New("failed to assert loan to BookLoans")
	}

	if bookLoan.ReturnDate != nil {
		// Hitung denda berdasarkan keterlambatan
		dueDate := bookLoan.DueDate
		if time.Now().After(dueDate) {
			daysLate := time.Now().Sub(dueDate).Hours() / 24
			finePerDay := 1000 // Misalkan denda per hari adalah 1000
			return int(daysLate) * finePerDay, nil
		}
	}
	return 0, nil
}
