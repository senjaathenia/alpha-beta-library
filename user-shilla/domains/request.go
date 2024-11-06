package domains

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// BookRequest merepresentasikan model permintaan peminjaman buku.
type BookRequest struct {
	ID          uint       `json:"id"`
	BookID      uint       `json:"book_id"`
	UserID      uuid.UUID  `json:"user_id"`
	Username     string    `json:"username"` // Menyertakan username di sini
	RequestDate time.Time  `gorm:"autoCreateTime" json:"request_date"`
	Status      string     `json:"status"` // Status permintaan (DIPROSES, DISETUJUI, DITOLAK)
	Reason      string     `json:"reason,omiempty"` // Alasan penolakan (jika ada)
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

// BookRequestUsecase interface untuk logika bisnis permintaan peminjaman buku.
type BookRequestUsecase interface {
	Create(ctx context.Context, request *BookRequest) error
	GetAll(ctx context.Context) ([]BookRequest, error)
	GetByID(ctx context.Context, id uint) (*BookRequest, error)
	Update(ctx context.Context, request *BookRequest) error
	GetBookByID(ctx context.Context, bookID uint) (*Book, error)
	Delete(ctx context.Context, id uint) error
	ApproveOrReject(ctx context.Context, id uint, approve bool, reason string) (time.Time, error) // Metode baru
	CheckStock(ctx context.Context, id uint) (bool, error)                           // Memeriksa ketersediaan stok buku
	CheckBookExists(ctx context.Context, id int) error
	CheckUserExists(ctx context.Context, userID uuid.UUID) error
	CreateLoanRecord(ctx context.Context, record *LoanRecord) error
	DecreaseStock(ctx context.Context, id uint) error // Mengurangi stok buku
	MoveToLoan(ctx context.Context, request *BookRequest) error // Method baru
	GetByUsernameRequest(ctx context.Context, username string) (*UserWithRequest, error)
}

// BookRequestRepository interface untuk operasi permintaan peminjaman buku.
// BookRequestRepository interface untuk operasi permintaan peminjaman buku.
type BookRequestRepository interface {
	Create(ctx context.Context, model interface{}) error
	GetAll(ctx context.Context, model interface{}, tmplScan interface{}, criteria string) error
	GetByID(ctx context.Context, model interface{}, tmplScan interface{}, selectedFields []string, order string, criteria string, id uint) (interface{}, error)
	Update(ctx context.Context, model interface{}, criteria string, args ...interface{}) error
	Delete(ctx context.Context, model interface{}, criteria string, id uint) error
	CheckStock(ctx context.Context, id uint) (bool, error) // Memeriksa ketersediaan stok buku
	CheckBookExists(ctx context.Context, id int) error
	CheckUserExists(ctx context.Context, userID uuid.UUID) error
	ApproveOrReject(ctx context.Context, id uint, approve bool, reason string) (time.Time, error) // Metode baru
	GetBookByID(ctx context.Context, bookID uint) (*Book, error)
	CreateLoanRecord(ctx context.Context, record *LoanRecord) error
	DecreaseStock(ctx context.Context, id uint) error // Mengurangi stok buku
	CreateLoan(ctx context.Context, request *BookRequest) error
	GetByUsernameRequest(ctx context.Context, username string) (*UserWithRequest, error)
	CheckUserIDByUsername(ctx context.Context, userID uuid.UUID, username string) error
}


type LoanUpdateRequest struct {
    Approved         bool   `json:"approved"`
    RejectionReason  string `json:"reason"`
    // Pastikan semua field yang digunakan di kode benar-benar ada
}

	type UserWithRequest struct {
		User      User        `json:"user"`
		BookRequest []BookRequest `json:"book_requests"`
	}

	type RequestBodyy struct {
		Username string `json:"username"` // Field untuk username
	}