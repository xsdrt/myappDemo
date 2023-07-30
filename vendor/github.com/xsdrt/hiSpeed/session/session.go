package session

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/scs/mssqlstore"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
)

type Session struct { //This will be exportable..
	CookieLifetime string
	CookiePersist  string
	CookieName     string
	CookieDomain   string
	SessionType    string
	CookieSecure   string
	DBPool         *sql.DB //Fill/populate the DBPool in the hiSpeed.go file sess...
}

func (h *Session) InitSession() *scs.SessionManager {
	var persist, secure bool

	//How long should sessions last?
	minutes, err := strconv.Atoi(h.CookieLifetime) //changed all c. to h. 7/29/23
	if err != nil {
		minutes = 60
	}

	//Should cookies persist?
	if strings.ToLower(h.CookiePersist) == "true" {
		persist = true
	}

	//Must cookies be secure?
	if strings.ToLower(h.CookieSecure) == "true" {
		secure = true
	}

	//Create the session
	session := scs.New()
	session.Lifetime = time.Duration(minutes) * time.Minute
	session.Cookie.Persist = persist
	session.Cookie.Name = h.CookieName
	session.Cookie.Secure = secure
	session.Cookie.Domain = h.CookieDomain
	session.Cookie.SameSite = http.SameSiteLaxMode

	//Which session store to use...
	switch strings.ToLower(h.SessionType) {
	case "redis":

	case "mysql", "mariadb":
		session.Store = mysqlstore.New(h.DBPool)

	case "postgres", "postgresql":
		session.Store = postgresstore.New(h.DBPool) //Set the sessions store

	case "mssql":
		session.Store = mssqlstore.New(h.DBPool)

	default:
		// cookie

	}

	return session
}
