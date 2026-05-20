package models

// Reservation holds reservation data
type Reservation struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate string
	EndDate   string
	RoomID    string
	RoomName  string
}
