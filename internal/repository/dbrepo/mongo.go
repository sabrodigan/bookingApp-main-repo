package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/sabrodigan/bookings-app/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type mongoReservation struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	FirstName string        `bson:"first_name"`
	LastName  string        `bson:"last_name"`
	Email     string        `bson:"email"`
	Phone     string        `bson:"phone"`
	StartDate string        `bson:"start_date,omitempty"`
	EndDate   string        `bson:"end_date,omitempty"`
	RoomID    string        `bson:"room_id,omitempty"`
	RoomName  string        `bson:"room_name,omitempty"`
}

func toModel(doc mongoReservation) models.Reservation {
	res := models.Reservation{
		FirstName: doc.FirstName,
		LastName:  doc.LastName,
		Email:     doc.Email,
		Phone:     doc.Phone,
		StartDate: doc.StartDate,
		EndDate:   doc.EndDate,
		RoomID:    doc.RoomID,
		RoomName:  doc.RoomName,
	}
	if !doc.ID.IsZero() {
		res.ID = doc.ID.Hex()
	}
	return res
}

func fromModel(res models.Reservation) mongoReservation {
	return mongoReservation{
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Email:     res.Email,
		Phone:     res.Phone,
		StartDate: res.StartDate,
		EndDate:   res.EndDate,
		RoomID:    res.RoomID,
		RoomName:  res.RoomName,
	}
}

// InsertReservation inserts a reservation into the database
func (m *mongoDBRepo) InsertReservation(res models.Reservation) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	collection := m.DB.Database("bookings_db").Collection("reservations")

	result, err := collection.InsertOne(ctx, fromModel(res))
	if err != nil {
		return "", err
	}

	oid, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return "", errors.New("cannot convert to object id")
	}

	return oid.Hex(), nil
}

// GetReservationByID gets a reservation by ID
func (m *mongoDBRepo) GetReservationByID(id string) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var doc mongoReservation
	collection := m.DB.Database("bookings_db").Collection("reservations")

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.Reservation{}, err
	}

	err = collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc)
	if err != nil {
		return models.Reservation{}, err
	}

	return toModel(doc), nil
}
