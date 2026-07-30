package service

import (
	"errors"

	"github.com/kaanoba25/library_api/models"
	"github.com/kaanoba25/library_api/repository"
)

type BookService struct {
	repo *repository.BookRepository
}

func NewBookService(repo *repository.BookRepository) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) GetAllBooks() ([]models.Book, error) {
	return s.repo.GetAll()
}

func (s *BookService) GetBookByID(id int) (models.Book, error) {
	if id <= 0 {
		return models.Book{}, errors.New("invalid book id")
	}
	return s.repo.GetByID(id)
}

func (s *BookService) CreateBook(book models.Book) (models.Book, error) {
	if book.Title == "" || book.Author == "" || book.ISBN == "" {
		return models.Book{}, errors.New("title, author and ISBN cannot be empty")
	}
	if book.TotalCopies <= 0 {
		return models.Book{}, errors.New("total copies must be greater than zero")
	}
	return s.repo.Create(book)
}