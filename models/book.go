package models

type Book struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
	ISBN string `json:"isbn"`
	TotalCopies int `json:"total_copies"`
	AvailableCopies int `json:"available_copies"`
}