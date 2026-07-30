package repository

import (
	"database/sql"
	"errors"

	"github.com/kaanoba25/library_api/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user models.User) (models.User, error) {
	query := `
		INSERT INTO users (full_name, email, password, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	err := r.db.QueryRow(query, user.FullName, user.Email, user.Password, user.Role).
		Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (models.User, error) {
	query := `SELECT id, full_name, email, password, role, created_at FROM users WHERE email = $1`
	row := r.db.QueryRow(query, email)

	var u models.User
	err := row.Scan(&u.ID, &u.FullName, &u.Email, &u.Password, &u.Role, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}

	return u, nil
}