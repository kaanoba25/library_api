package repository

import (
	"errors"

	"github.com/kaanoba25/library_api/models"
)

type BookRepository struct {
	books []models.Book
	nextID int
}


// Constructor 
func NewBookRepository() *BookRepository {
	return &BookRepository{
		books: []models.Book{
			{ID: 1, Title: "Nutuk", Author: "Mustafa Kemal Atatürk", ISBN: "9789751600000", TotalCopies: 5, AvailableCopies: 5},
			{ID: 2, Title: "Suç ve Ceza", Author: "Dostoyevski", ISBN: "9789750711111", TotalCopies: 3, AvailableCopies: 3},
		},
		nextID: 3,
		
	}
}

// Methods
func (b *BookRepository) GetAllBooks() []models.Book {
	return b.books
}

func (b *BookRepository) GetByID(id int) (models.Book, error) {
	for _ , book := range b.books {
		if book.ID == id {
			return book, nil
		}
	}

	return models.Book{}, errors.New("Book not found")
}

func (b *BookRepository) Create(book models.Book) models.Book {
	book.ID = b.nextID	
	book.AvailableCopies = book.TotalCopies
	b.nextID++
	b.books = append(b.books, book)
	return book
}