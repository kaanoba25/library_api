package service

import (
	"errors"

	"github.com/kaanoba25/library_api/models"
	"github.com/kaanoba25/library_api/repository"
)

type LoanService struct {
	repo *repository.LoanRepository
}

func NewLoanService(repo *repository.LoanRepository) *LoanService {
	return &LoanService{repo: repo}
}

func (s *LoanService) BorrowBook(userID, bookID int) (models.Loan, error) {
	if bookID <= 0 {
		return models.Loan{}, errors.New("invalid book id")
	}

	// Rule: Maximum 3 active loans per user
	count, err := s.repo.GetActiveLoansCount(userID)
	if err != nil {
		return models.Loan{}, err
	}
	if count >= 3 {
		return models.Loan{}, errors.New("user has reached maximum borrow limit (3 books)")
	}

	return s.repo.BorrowBook(userID, bookID)
}

func (s *LoanService) ReturnBook(loanID, userID int) (models.Loan, error) {
	if loanID <= 0 {
		return models.Loan{}, errors.New("invalid loan id")
	}
	return s.repo.ReturnBook(loanID, userID)
}