
package domains

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type BookLoans struct {
	ID           int        `json:"id"`
	BookID       int        `json:"book_id"`
	UserID       uuid.UUID  `json:"user_id"`
	Username     string    `json:"username"` // Menyertakan username di sini
	LoanDate     time.Time  `json:"loan_date"`
	DueDate      time.Time  `json:"due_date"`
	ReturnStatus bool       `json:"return_status"`
	ReturnDate   *time.Time `json:"return_date,omitempty"` // Optional field for return date
    LateFee       int       `json:"late_fee,omitempty"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at"`
}

// LoanUsecase interface defines the methods for loan business logic.
type LoanUsecase interface {
	Create(ctx context.Context, loans *BookLoans) error
	GetAll(ctx context.Context) ([]BookLoans, error)
	GetByID(ctx context.Context, id uint) (*BookLoans, error)
	Update(ctx context.Context, loans *BookLoans) error
	Delete(ctx context.Context, id uint) error
	CheckBookExists(ctx context.Context, id int) error
	CheckUserExists(ctx context.Context, userID uuid.UUID) error
	CreateLoanRecord(ctx context.Context, record *LoanRecord) error
	Return(ctx context.Context, loanID uint, returnDate time.Time)error
	GetByUsernameLoans(ctx context.Context, username string) (*UserWithLoans, error)
}

type LoanRecord struct {
	ID            uint      `gorm:"primaryKey"`
	BookRequestID uint      `gorm:"not null"`
	BookID        uint      `gorm:"not null"`
	UserID        uuid.UUID `gorm:"not null"`
	DueDate       time.Time `gorm:"not null"`
	LoanDate      time.Time `gorm:"not null"`
	ReturnStatus  bool      // True jika buku sudah dikembalikan, False jika belum
}


type BookLoansRepository interface {
	Create(ctx context.Context, model interface{}) error
	GetAll(ctx context.Context, model interface{}, tmplScan interface{}, criteria string) error
	GetByID(ctx context.Context, model interface{}, tmplScan interface{}, selectedFields []string, order string, criteria string, id uint) (interface{}, error)
	Update(ctx context.Context, model interface{}, criteria string, args ...interface{}) error
	Delete(ctx context.Context, model interface{}, criteria string, id uint) error
	CheckBookExists(ctx context.Context, id int) error
	CheckUserExists(ctx context.Context, userID uuid.UUID) error
	CreateLoanRecord(ctx context.Context, record *LoanRecord) error
	IncreaseBookStock(ctx context.Context, bookID int) error // Method baru
	Return(ctx context.Context, loanID uint, returnDate time.Time, lateFee int)error
	GetByUsernameLoans(ctx context.Context, username string) (*UserWithLoans, error)
}

type UserWithLoans struct {
	User      User        `json:"user"`
	BookLoans []BookLoans `json:"book_loans"`
}

type RequestBody struct {
	Username string `json:"username"` // Field untuk username
}