package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kaanoba25/library_api/config"
	"github.com/kaanoba25/library_api/handler"
	"github.com/kaanoba25/library_api/middleware"
	"github.com/kaanoba25/library_api/repository"
	"github.com/kaanoba25/library_api/service"
)

func main() {
	// 1. Veritabanı Bağlantısı
	db := config.InitDB()
	defer db.Close()

	// 2. Katmanlar (Dependency Injection)
	bookRepo := repository.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo)
	bookHandler := handler.NewBookHandler(bookService)

	r := mux.NewRouter()

	r.Use(middleware.JSONContentTypeMiddleware)
	r.Use(middleware.LoggerMiddleware)

	r.HandleFunc("/api/books", bookHandler.GetAll).Methods("GET")
	r.HandleFunc("/api/books/{id}", bookHandler.GetByID).Methods("GET")
	r.HandleFunc("/api/books", bookHandler.Create).Methods("POST")

	log.Println("Library API (PostgreSQL) 8080 portunda çalışıyor...")
	log.Fatal(http.ListenAndServe(":8080", r))
}