package config

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // Driver'ın init fonksiyonunu çalıştırmak için anonim import
)

func InitDB() *sql.DB {
	// Bağlantı cümlesi (Connection String)
	connStr := "host=localhost port=5432 user=postgres password=covboy dbname=library_api sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Veritabanı sürücüsü başlatılamadı: %v", err)
	}

	// Bağlantıyı test edelim
	if err := db.Ping(); err != nil {
		log.Fatalf("Veritabanına bağlanılamadı (PostgreSQL açık mı?): %v", err)
	}

	log.Println("PostgreSQL veritabanına başarıyla bağlanıldı!")

	// Tablomuzu otomatik oluşturalım
	createTables(db)

	return db
}

func createTables(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS books (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		author VARCHAR(255) NOT NULL,
		isbn VARCHAR(50) UNIQUE NOT NULL,
		total_copies INT NOT NULL DEFAULT 1,
		available_copies INT NOT NULL DEFAULT 1
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Tablo oluşturulurken hata: %v", err)
	}
}