package data

import (
	"database/sql"
	"os"

	db2 "github.com/upper/db/v4" // Added upper/db to use as an ORM
	"github.com/upper/db/v4/adapter/mssql"
	"github.com/upper/db/v4/adapter/mysql"
	"github.com/upper/db/v4/adapter/postgresql"
)

//A file to store various things from the database, such as a structure for models that correspond to data in tables stored in the db and allow it to be easy for user to use...

var db *sql.DB        // Package level variable assigned a value in the New func; now accessible to the entire package...
var upper db2.Session //Pakage level variable available to the entire data package

type Models struct {
	// Any models inserted here (and in the New function)
	// are easily available thru the whole application...
}

func New(databasePool *sql.DB) Models {
	db = databasePool // This variable is package wide accessible...

	//Check what db using first...
	if os.Getenv("DATABASE_TYPE") == "mysql" || os.Getenv("DATABASE_TYPE") == "mariadb" {
		// TODO, going to just concentrate on postgres for now...
		upper, _ = mysql.New(databasePool)
	} else if os.Getenv("DATABASE_TYPE") == "mssql" {
		// TODO also for MSSQL...
		upper, _ = mssql.New(databasePool)
	} else {
		upper, _ = postgresql.New(databasePool)
	}

	return Models{}
}
