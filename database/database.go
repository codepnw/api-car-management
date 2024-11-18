package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func ConnectPostgres(connStr string) *sql.DB {
	fmt.Println("waiting the database start up...")
	time.Sleep(2 * time.Second)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}

	if err := conn.Ping(); err != nil {
		log.Fatalf("error conneting to the database: %v", err)
	}

	db = conn

	fmt.Println("database conneted...")
	return db
}

func GetDB() *sql.DB {
	return db
}

func ExecuteSQLSchema(db *sql.DB, fileName string) error {
	sqlFile, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("error reading schema file: %v", err)
	}

	_, err = db.Exec(string(sqlFile))
	if err != nil {
		return fmt.Errorf("error executing schema file: %v", err)
	}

	return nil
}
