package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kaanoba25/library_api/config"
	"github.com/kaanoba25/library_api/handler"
	"github.com/kaanoba25/library_api/middleware"
	"github.com/kaanoba25/library_api/repository"
	"github.com/kaanoba25/library_api/service"
)

func main() {
	_ = godotenv.Load()

	db := config.InitDB()
	defer db.Close()

	// Repositories
	bookRepo := repository.NewBookRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Services
	bookService := service.NewBookService(bookRepo)
	userService := service.NewUserService(userRepo)

	// Handlers
	bookHandler := handler.NewBookHandler(bookService)
	userHandler := handler.NewUserHandler(userService)

	r := mux.NewRouter()
	r.Use(middleware.JSONContentTypeMiddleware)
	r.Use(middleware.LoggerMiddleware)

	// Public Auth Routes
	r.HandleFunc("/api/auth/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/api/auth/login", userHandler.Login).Methods("POST")

	// Public Book Routes
	r.HandleFunc("/api/books", bookHandler.GetAll).Methods("GET")
	r.HandleFunc("/api/books/{id}", bookHandler.GetByID).Methods("GET")

	// Protected Admin Routes
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.Use(middleware.RequireRole("admin"))

	protected.HandleFunc("/books", bookHandler.Create).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}