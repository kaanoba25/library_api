package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/kaanoba25/library_api/docs"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kaanoba25/library_api/config"
	"github.com/kaanoba25/library_api/handler"
	"github.com/kaanoba25/library_api/middleware"
	"github.com/kaanoba25/library_api/repository"
	"github.com/kaanoba25/library_api/service"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	_ = godotenv.Load()

	db := config.InitDB()
	defer db.Close()

	// Repositories
	bookRepo := repository.NewBookRepository(db)
	userRepo := repository.NewUserRepository(db)
	loanRepo := repository.NewLoanRepository(db)

	// Services
	bookService := service.NewBookService(bookRepo)
	userService := service.NewUserService(userRepo)
	loanService := service.NewLoanService(loanRepo)

	// Handlers
	bookHandler := handler.NewBookHandler(bookService)
	userHandler := handler.NewUserHandler(userService)
	loanHandler := handler.NewLoanHandler(loanService)

	r := mux.NewRouter()
	r.Use(middleware.JSONContentTypeMiddleware)
	r.Use(middleware.LoggerMiddleware)

	// SWAGGER UI ENDPOINT
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Public Routes
	r.HandleFunc("/api/auth/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/api/auth/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/api/books", bookHandler.GetAll).Methods("GET")
	r.HandleFunc("/api/books/{id}", bookHandler.GetByID).Methods("GET")

	// Protected Routes (All Logged-in Users)
	userProtected := r.PathPrefix("/api").Subrouter()
	userProtected.Use(middleware.AuthMiddleware)
	userProtected.HandleFunc("/loans/borrow", loanHandler.Borrow).Methods("POST")
	userProtected.HandleFunc("/loans/return/{id}", loanHandler.Return).Methods("POST")

	// Protected Admin Routes
	adminProtected := r.PathPrefix("/api").Subrouter()
	adminProtected.Use(middleware.AuthMiddleware)
	adminProtected.Use(middleware.RequireRole("admin"))
	adminProtected.HandleFunc("/books", bookHandler.Create).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}