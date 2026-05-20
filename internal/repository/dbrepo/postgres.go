package dbrepo

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/sabrodigan/bookings-app/internal/config"
	"github.com/sabrodigan/bookings-app/internal/models"
	"github.com/sabrodigan/bookings-app/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// NewPostgresRepo creates a new PostgreSQL repository
func NewPostgresRepo(db *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  db,
	}
}

func (m *postgresDBRepo) InsertReservation(res models.Reservation) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if res.StartDate == "" {
		res.StartDate = time.Now().Format("2006-01-02")
	}
	if res.EndDate == "" {
		res.EndDate = time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	}

	roomID, err := strconv.Atoi(res.RoomID)
	if err != nil || roomID < 1 {
		roomID = 1
	}

	var id int
	query := `
		INSERT INTO reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5::date, $6::date, $7, NOW(), NOW())
		RETURNING id`

	err = m.DB.QueryRowContext(ctx, query,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		roomID,
	).Scan(&id)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(id), nil
}

func (m *postgresDBRepo) GetReservationByID(id string) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res models.Reservation
	var dbID int
	var roomID int

	query := `
		SELECT r.id, r.first_name, r.last_name, r.email, r.phone,
		       r.start_date::text, r.end_date::text, r.room_id, COALESCE(rm.room_name, '')
		FROM reservations r
		LEFT JOIN rooms rm ON r.room_id = rm.id
		WHERE r.id = $1`

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&dbID,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Phone,
		&res.StartDate,
		&res.EndDate,
		&roomID,
		&res.RoomName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, errors.New("reservation not found")
		}
		return res, err
	}

	res.ID = strconv.Itoa(dbID)
	res.RoomID = strconv.Itoa(roomID)
	return res, nil
}
