package models

type Order struct {
	ID     int
	UserID int
	Total  float64
	Status string
}
