package services

import (
	"errors"

	"github.com/kaanoba25/library_api/models"
	"github.com/kaanoba25/library_api/repository"
)

type BookService struct {
	repo *repository.BookRepository
}

// Constructor
func NewBookService(repo *repository.BookRepository) *BookService {
	return &BookService{repo: repo}
} 	

// Methods
func (b *BookService) GetAllBooks() []models.Book {
	return b.repo.GetAllBooks()
}

func (b *BookService) GetBookByID(id int) (models.Book, error) {
	if id <= 0 {
		return models.Book{}, errors.New("invalid book id")
	}

	return b.repo.GetByID(id)
}

func (b *BookService) CreateBook(book models.Book) (models.Book, error) {
	if book.Title == "" || book.Author == "" {
		return models.Book{}, errors.New("title and author cannot be empty")
	}

	if book.TotalCopies <= 0 {
		return models.Book{}, errors.New("total copies must be greater than zero")
	}

	return b.repo.Create(book), nil
}