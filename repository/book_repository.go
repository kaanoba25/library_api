package repository

import (
	"database/sql"
	"errors"

	"github.com/kaanoba25/library_api/models"
)

type BookRepository struct {
	db *sql.DB
}

// Artık *sql.DB kabul ediyor
func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) GetAll() ([]models.Book, error) {
	query := `SELECT id, title, author, isbn, total_copies, available_copies FROM books`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.ISBN, &b.TotalCopies, &b.AvailableCopies); err != nil {
			return nil, err
		}
		books = append(books, b)
	}

	return books, nil
}

func (r *BookRepository) GetByID(id int) (models.Book, error) {
	query := `SELECT id, title, author, isbn, total_copies, available_copies FROM books WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var b models.Book
	err := row.Scan(&b.ID, &b.Title, &b.Author, &b.ISBN, &b.TotalCopies, &b.AvailableCopies)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Book{}, errors.New("book not found")
		}
		return models.Book{}, err
	}

	return b, nil
}

func (r *BookRepository) Create(book models.Book) (models.Book, error) {
	query := `
		INSERT INTO books (title, author, isbn, total_copies, available_copies)
		VALUES ($1, $2, $3, $4, $4)
		RETURNING id, available_copies`

	// RETURNING id sayesinde PostgreSQL'in atadığı otomatik ID'yi alıyoruz
	err := r.db.QueryRow(query, book.Title, book.Author, book.ISBN, book.TotalCopies).Scan(&book.ID, &book.AvailableCopies)
	if err != nil {
		return models.Book{}, err
	}

	return book, nil
}