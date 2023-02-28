package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

//Remember for Integration Tests will require a spun up new docker container/image and an empty postgres DB (will populate with known table structures)
//Then run tests and when finished (all pass ) get rid of the docker image and DB....

var (
	host     = "localhost"
	user     = "postgres"
	password = "secret"
	dbName   = "hispeed_test"
	port     = "5435"
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5" //data source name...
)

var dummyUser = User{
	FirstName: "No",
	LastName:  "one",
	Email:     "not@here.com",
	Active:    1,
	Password:  "password",
}

var models Models
var testDB *sql.DB
var resource *dockertest.Resource
var pool *dockertest.Pool

func TestMain(m *testing.M) {
	os.Setenv("DATABASE_TYPE", "postgres")

	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	pool = p

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13.4",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"}, //This is the port inside the docker image...
		PortBindings: map[docker.Port][]docker.PortBinding{ //Bind the port to the local machine...
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		_ = pool.Purge(resource) //if an err get rid of the docker image if didn't work; then log and exit from the program...
		log.Fatalf("could not start resource: %s", err)
	}

	if err := pool.Retry(func() error { //this func is part of the docker test pkg.  Allows docker time to spin up the container and DB...
		var err error
		testDB, err = sql.Open("pgx", fmt.Sprintf(dsn, host, port, user, password, dbName)) //Open up the connection to the database
		if err != nil {
			return err
		}
		return testDB.Ping() //Should keep pinging until the connection works properly (The ddb is built and waiting connections in the docker image)
	}); err != nil {
		_ = pool.Purge(resource) // Again, if doesn't work, clean up the resources and log and exit...
		log.Fatalf("could not connect to docker: %s", err)
	}

}
