package data

import (
	"fmt"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	db2 "github.com/upper/db/v4"
)

func TestNew(t *testing.T) {
	fakeDB, _, _ := sqlmock.New()
	defer fakeDB.Close()

	_ = os.Setenv("DATABASE_TYPE", "postgres")
	m := New(fakeDB)
	if fmt.Sprintf("%T", m) != "data.Models" {
		t.Error("Wrong type", fmt.Sprintf("%T", m))
	}

	_ = os.Setenv("DATABASE_TYPE", "mysql")
	m = New(fakeDB)
	if fmt.Sprintf("%T", m) != "data.Models" {
		t.Error("Wrong type", fmt.Sprintf("%T", m))
	}

	_ = os.Setenv("DATABASE_TYPE", "mssql")
	m = New(fakeDB)
	if fmt.Sprintf("%T", m) != "data.Models" {
		t.Error("Wrong type", fmt.Sprintf("%T", m))
	}
}

func TestGetInsertID(t *testing.T) {
	var id db2.ID
	id = int64(1)

	returnedID := getInsertId(id) //from postgres (int64)
	if fmt.Sprintf("%T", returnedID) != "int" {
		t.Error("wrong type returned")
	}

	id = 1

	returnedID = getInsertId(id) //For mysql and mssql just an int...
	if fmt.Sprintf("%T", returnedID) != "int" {
		t.Error("wrong type returned")
	}

}
