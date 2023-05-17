package db

import (
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testMock sqlmock.Sqlmock

func TestMain(m *testing.M) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	testMock = mock
	testQueries = New(db)

	exitCode := m.Run()

	if err := testMock.ExpectationsWereMet(); err != nil {
		panic(err)
	}

	os.Exit(exitCode)
}
