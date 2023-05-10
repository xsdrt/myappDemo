package data

import (
	"database/sql"
	"fmt"
	"os"

	db2 "github.com/upper/db/v4" // Added upper/db to use as an ORM... upper/db still allows the use of developer written important/complicated raw sql as needed by the app...
	"github.com/upper/db/v4/adapter/mssql"
	"github.com/upper/db/v4/adapter/mysql"
	"github.com/upper/db/v4/adapter/postgresql"
)

// TODO: Need to add validation to all the models...

//A file to store various things from the database, such as a structure for models that correspond to data in tables stored in the db and allow it to be easy for user to use...

var db *sql.DB        // Package level variable assigned a value in the New func; now accessible to the entire package...
var upper db2.Session //Package level variable available to the entire data package

// Models is the wrapper for all database models...
type Models struct {
	// Any models inserted here (and in the New function)
	// are easily available/accessible in the whole application...
	Users  User
	Tokens Token
}

func New(databasePool *sql.DB) Models {
	db = databasePool // This variable is package wide accessible...

	//Check what db using first...
	if os.Getenv("DATABASE_TYPE") == "mysql" || os.Getenv("DATABASE_TYPE") == "mariadb" {
		// TODO, going to just concentrate on postgres for now...Ok added mySql and MSSQL
		upper, _ = mysql.New(databasePool)
	} else if os.Getenv("DATABASE_TYPE") == "mssql" {
		// TODO also for MSSQL...
		upper, _ = mssql.New(databasePool)
	} else {
		upper, _ = postgresql.New(databasePool)
	}

	return Models{
		Users:  User{},
		Tokens: Token{},
	}
}

func getInsertId(i db2.ID) int { //Get from the db2 (see import above) this used in the user.go
	idType := fmt.Sprintf("%T", i)
	if idType == "int64" { //Postgresql returns this type...
		return int(i.(int64))
	}

	return i.(int)
}
