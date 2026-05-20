package repository

import "github.com/sabrodigan/bookings-app/internal/models"

type DatabaseRepo interface {
	InsertReservation(res models.Reservation) (string, error)
	GetReservationByID(id string) (models.Reservation, error)
}
