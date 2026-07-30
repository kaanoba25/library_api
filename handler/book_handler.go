package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kaanoba25/library_api/models"
	"github.com/kaanoba25/library_api/service"

	"github.com/kaanoba25/library_api/utils"
)

type BookHandler struct {
	service *service.BookService
}

// Constructor
func NewBookHandler(service *service.BookService) *BookHandler {
	return &BookHandler{service: service}
}

// Methods
func (h *BookHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	books, err := h.service.GetAllBooks()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, books)
}

func (b *BookHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid id param")
		return
	}

	book, err := b.service.GetBookByID(id)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, book)
}

func (b *BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	var book models.Book

	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	createdBook, err := b.service.CreateBook(book)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, createdBook)
}