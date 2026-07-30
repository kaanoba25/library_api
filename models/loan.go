package models

import "time"

type LoanStatus string

const (
	StatusBorrowed LoanStatus = "borrowed"
	StatusReturned LoanStatus = "returned"
)

type Loan struct {
	ID         int        `json:"id"`
	UserID     int        `json:"user_id"`
	BookID     int        `json:"book_id"`
	BorrowedAt time.Time  `json:"borrowed_at"`
	DueDate    time.Time  `json:"due_date"`
	ReturnedAt *time.Time `json:"returned_at,omitempty"`
	Status     LoanStatus `json:"status"`
}

type BorrowRequest struct {
	BookID int `json:"book_id"`
}