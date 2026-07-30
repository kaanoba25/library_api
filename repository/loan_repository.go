package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/kaanoba25/library_api/models"
)

type LoanRepository struct {
	db *sql.DB
}

func NewLoanRepository(db *sql.DB) *LoanRepository {
	return &LoanRepository{db: db}
}

func (r *LoanRepository) GetActiveLoansCount(userID int) (int, error) {
	query := `SELECT COUNT(*) FROM loans WHERE user_id = $1 AND status = 'borrowed'`
	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	return count, err
}

func (r *LoanRepository) BorrowBook(userID, bookID int) (models.Loan, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return models.Loan{}, err
	}
	defer tx.Rollback()

	// 1. Check stock & decrement available copies
	var available int
	err = tx.QueryRow(`SELECT available_copies FROM books WHERE id = $1 FOR UPDATE`, bookID).Scan(&available)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Loan{}, errors.New("book not found")
		}
		return models.Loan{}, err
	}

	if available <= 0 {
		return models.Loan{}, errors.New("no copies available for this book")
	}

	_, err = tx.Exec(`UPDATE books SET available_copies = available_copies - 1 WHERE id = $1`, bookID)
	if err != nil {
		return models.Loan{}, err
	}

	// 2. Create loan record (Default 14 days due date)
	dueDate := time.Now().Add(14 * 24 * time.Hour)
	query := `
		INSERT INTO loans (user_id, book_id, due_date, status)
		VALUES ($1, $2, $3, 'borrowed')
		RETURNING id, borrowed_at, status`

	var loan models.Loan
	loan.UserID = userID
	loan.BookID = bookID
	loan.DueDate = dueDate

	err = tx.QueryRow(query, userID, bookID, dueDate).Scan(&loan.ID, &loan.BorrowedAt, &loan.Status)
	if err != nil {
		return models.Loan{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.Loan{}, err
	}

	return loan, nil
}

func (r *LoanRepository) ReturnBook(loanID, userID int) (models.Loan, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return models.Loan{}, err
	}
	defer tx.Rollback()

	// 1. Get loan record
	var loan models.Loan
	query := `SELECT id, user_id, book_id, borrowed_at, due_date, status FROM loans WHERE id = $1 AND user_id = $2 AND status = 'borrowed'`
	err = tx.QueryRow(query, loanID, userID).Scan(&loan.ID, &loan.UserID, &loan.BookID, &loan.BorrowedAt, &loan.DueDate, &loan.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Loan{}, errors.New("active loan not found")
		}
		return models.Loan{}, err
	}

	// 2. Update loan status
	now := time.Now()
	_, err = tx.Exec(`UPDATE loans SET status = 'returned', returned_at = $1 WHERE id = $2`, now, loanID)
	if err != nil {
		return models.Loan{}, err
	}

	// 3. Increment book stock
	_, err = tx.Exec(`UPDATE books SET available_copies = available_copies + 1 WHERE id = $1`, loan.BookID)
	if err != nil {
		return models.Loan{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.Loan{}, err
	}

	loan.Status = models.StatusReturned
	loan.ReturnedAt = &now
	return loan, nil
}