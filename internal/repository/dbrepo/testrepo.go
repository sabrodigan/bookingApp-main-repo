package dbrepo

import (
	"github.com/sabrodigan/bookings-app/internal/config"
	"github.com/sabrodigan/bookings-app/internal/models"
	"github.com/sabrodigan/bookings-app/internal/repository"
)

type testDBRepo struct {
	App *config.AppConfig
}

// NewTestRepo creates a new test repository
func NewTestRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}

func (m *testDBRepo) InsertReservation(res models.Reservation) (string, error) {
	return "test_id", nil
}

func (m *testDBRepo) GetReservationByID(id string) (models.Reservation, error) {
	return models.Reservation{}, nil
}
