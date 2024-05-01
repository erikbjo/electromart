package db

import (
	"context"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
	"time"
)

var Client *sql.DB

// getDataSourceName returns a Data Source Name string for connecting to a MySQL database.
func getDataSourceName() string {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbNet := os.Getenv("DB_NET")
	dbIP := os.Getenv("DB_IP")
	dbName := os.Getenv("DB_NAME")

	cfg := mysql.Config{
		User:                 dbUser,
		Passwd:               dbPassword,
		Net:                  dbNet,
		Addr:                 dbIP,
		DBName:               dbName,
		AllowNativePasswords: true,
	}

	return cfg.FormatDSN()
}

// OpenDatabaseConnection opens a connection to the database and returns a pointer to the database.
func OpenDatabaseConnection() *sql.DB {
	db, err := sql.Open("mysql", getDataSourceName())
	if err != nil {
		log.Fatal("Error when opening the database: ", err)
	}

	// Call the ping function to ensure the database is up and running
	pingDatabase(db)

	return db
}

// pingDatabase pings the database to verify the connection.
func pingDatabase(db *sql.DB) {
	// Ping the database to verify the connection.
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal("Error when pinging database: ", err)
	}
}
