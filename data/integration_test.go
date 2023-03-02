//go:build integration

// run test with the following code: go test . --tags integration --count=1
package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
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
		log.Println("Error:", err)
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
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	resource, err = pool.RunWithOptions(&opts)
	if err != nil { //if a fatal error; purge the resource...
		log.Println("Error:", err)
		_ = pool.Purge(resource)
		log.Fatalf("could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		var err error
		testDB, err = sql.Open("pgx", fmt.Sprintf(dsn, host, port, user, password, dbName))
		if err != nil {
			log.Println("Error:", err)
			return err
		}
		return testDB.Ping() //Should keep pinging db until works properly...
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not connect to docker: %s", err)
	}

	err = createTables(testDB)
	if err != nil { //So exit ,because no tables...
		log.Fatalf("error creating tables: %s", err)
	}

	models = New(testDB)

	code := m.Run()

	os.Exit(code)

}

func createTables(db *sql.DB) error {
	stmt := `
	CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

drop table if exists users cascade;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    user_active integer NOT NULL DEFAULT 0,
    email character varying(255) NOT NULL UNIQUE,
    password character varying(60) NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE TRIGGER set_timestamp
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();

drop table if exists remember_tokens;

CREATE TABLE remember_tokens (
    id SERIAL PRIMARY KEY,
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    remember_token character varying(100) NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE TRIGGER set_timestamp
    BEFORE UPDATE ON remember_tokens
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();

drop table if exists tokens;

CREATE TABLE tokens (
    id SERIAL PRIMARY KEY,
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    first_name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    token character varying(255) NOT NULL,
    token_hash bytea NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now(),
    expiry timestamp without time zone NOT NULL
);

CREATE TRIGGER set_timestamp
    BEFORE UPDATE ON tokens
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();
	`
	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

func TestUser_Table(t *testing.T) {
	s := models.Users.Table()
	if s != "users" {
		t.Error("wrong table name returned: ", s)
	}
}
