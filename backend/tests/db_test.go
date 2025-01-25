package db_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestMockedDatabaseConnection(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))

	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}

	defer mockDB.Close()

	// Wrap sql.DB with sqlx
	db := sqlx.NewDb(mockDB, "sqlmock")

	// Mock a ping
	mock.ExpectPing()

	// Test the connection
	if err := db.Ping(); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}

}
