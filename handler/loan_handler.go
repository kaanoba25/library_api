package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kaanoba25/library_api/middleware"
	"github.com/kaanoba25/library_api/models"
	"github.com/kaanoba25/library_api/service"
	"github.com/kaanoba25/library_api/utils"
)

type LoanHandler struct {
	service *service.LoanService
}

func NewLoanHandler(service *service.LoanService) *LoanHandler {
	return &LoanHandler{service: service}
}

func (h *LoanHandler) Borrow(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized user")
		return
	}

	var req models.BorrowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	loan, err := h.service.BorrowBook(userID, req.BookID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, loan)
}

func (h *LoanHandler) Return(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized user")
		return
	}

	vars := mux.Vars(r)
	loanID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid loan ID")
		return
	}

	loan, err := h.service.ReturnBook(loanID, userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, loan)
}