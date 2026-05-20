package driver

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// DB holds the database connection pool
type DB struct {
	Client *mongo.Client
	SQL    *sql.DB
}

var dbConn = &DB{}

// ConnectMongo creates database pool for MongoDB
func ConnectMongo(dsn string) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(dsn)
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	dbConn.Client = client
	return dbConn, nil
}

// ConnectPostgres creates a database pool for PostgreSQL
func ConnectPostgres(dsn string) (*DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("postgres ping failed: %w", err)
	}

	dbConn.SQL = db
	return dbConn, nil
}

// Close releases database connections
func (db *DB) Close() error {
	if db.SQL != nil {
		return db.SQL.Close()
	}
	if db.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return db.Client.Disconnect(ctx)
	}
	return nil
}
