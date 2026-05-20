package dbrepo

import (
	"github.com/sabrodigan/bookings-app/internal/config"
	"github.com/sabrodigan/bookings-app/internal/repository"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type mongoDBRepo struct {
	App *config.AppConfig
	DB  *mongo.Client
}

// NewMongoRepo creates a new MongoDB repository
func NewMongoRepo(conn *mongo.Client, a *config.AppConfig) repository.DatabaseRepo {
	return &mongoDBRepo{
		App: a,
		DB:  conn,
	}
}
