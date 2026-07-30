package main

import (
	_ "fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kaanoba25/library_api/handler"
	"github.com/kaanoba25/library_api/middleware"
	"github.com/kaanoba25/library_api/repository"
	"github.com/kaanoba25/library_api/services"
)

func main() {
	bookRepo := repository.NewBookRepository()
	bookService := services.NewBookService(bookRepo)
	bookHandler := handler.NewBookHandler(bookService)

	r := mux.NewRouter()

	// Global middleware
	r.Use(middleware.JSONContentTypeMiddleware)
	r.Use(middleware.LoggerMiddleware)

	// Routes
	r.HandleFunc("/api/books", bookHandler.GetAll).Methods("GET")
	r.HandleFunc("/api/books/{id}", bookHandler.GetByID).Methods("GET")
	r.HandleFunc("/api/books", bookHandler.Create).Methods("POST")

	log.Println("Library API working on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}