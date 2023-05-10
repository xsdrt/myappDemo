package hiSpeed

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/xsdrt/hiSpeed/render"
	"github.com/xsdrt/hiSpeed/session"
)

const version = "1.0.0"

// HiSpeed is the overall type for the HiSpeed Package.  Member that are exported in this type
// are available to any application that uses it...
type HiSpeed struct {
	AppName  string
	Debug    bool
	Version  string
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	RootPath string
	Routes   *chi.Mux
	Render   *render.Render
	Session  *scs.SessionManager
	DB       Database
	JetViews *jet.Set
	config   config //This will not be exported as there is no reason any app that imports HiSpeed should have access to the config...
}

// This type will not be exported but will hold all the config values for this package...
type config struct { //Application config
	port        string
	renderer    string
	cookie      cookieConfig //setup in types.go
	sessionType string       //setup in types.go
	database    databaseConfig
}

// New reads the .env file, creates our application config, populates the HiSpeed type with settings
// based on .env values, and creates necessary folders and files if they don't exist...
func (h *HiSpeed) New(rootPath string) error {
	pathConfig := initPaths{
		rootPath:    rootPath,
		folderNames: []string{"handlers", "migrations", "views", "data", "public", "tmp", "logs", "middleware"},
	}

	err := h.Init(pathConfig)
	if err != nil {
		return err
	}

	err = h.checkDotEnv(rootPath)
	if err != nil {
		return err
	}

	//read .env file
	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	//Create loggers
	infoLog, errorLog := h.startLoggers()

	// Connect to the database after loggers created and before populating the variable...
	if os.Getenv("DATABASE_TYPE") != "" { // If not equal to empty string connect to the database...
		db, err := h.OpenDB(os.Getenv("DATABASE_TYPE"), h.BuildDSN()) //h.OpenDB is in the driver.go file (just a reference to remember where coming from...)
		if err != nil {
			errorLog.Println(err)
			os.Exit(1) //If cannot connect to database , serious issue so exit...
		}
		h.DB = Database{
			DataType: os.Getenv("DATABASE_TYPE"),
			Pool:     db,
		}
	}

	h.InfoLog = infoLog
	h.ErrorLog = errorLog
	h.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	h.Version = version
	h.RootPath = rootPath
	h.Routes = h.routes().(*chi.Mux) //Cast to a * of chi.Mux from a htttphandler

	h.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
		cookie: cookieConfig{
			name:     os.Getenv("COOKIE_NAME"),
			lifetime: os.Getenv("COOKIE_LIFETIME"),
			persist:  os.Getenv("COOKIE_PERSISTS"),
			secure:   os.Getenv("COOKIE_SECURE"),
			domain:   os.Getenv("COOKIE_DOMAIN"),
		},
		sessionType: os.Getenv("SESSION_TYPE"),
		database: databaseConfig{ //Populate the database config and values in case needed...
			database: os.Getenv("DATABASE_TYPE"),
			dsn:      h.BuildDSN(),
		},
	}

	//Need to create a session...  just like render is in its own package , putting session in its own pkg also...

	sess := session.Session{
		CookieLifetime: h.config.cookie.lifetime,
		CookiePersist:  h.config.cookie.persist,
		CookieName:     h.config.cookie.name,
		SessionType:    h.config.sessionType,
		CookieDomain:   h.config.cookie.domain,
		DBPool:         h.DB.Pool,
	}

	h.Session = sess.InitSession()

	var views = jet.NewSet(
		jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rootPath)),
		jet.InDevelopmentMode(),
	)
	h.JetViews = views

	h.createRenderer()

	return nil
}

func (h *HiSpeed) Init(p initPaths) error {
	root := p.rootPath
	for _, path := range p.folderNames {
		//Create folder if it doesn't exist..
		err := h.CreateDirIfNotExist(root + "/" + path)
		if err != nil {
			return err
		}
	}
	return nil
}

// This (ListenAndServe) will start the web server and keep it running in the back ground...
func (h *HiSpeed) ListenAndServe() {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")), //Get the variable from the env file
		ErrorLog:     h.ErrorLog,
		Handler:      h.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second, //Really long for development purposes, can/should change for prod...
	}
	//After the server shuts down close the DB
	defer h.DB.Pool.Close()

	h.InfoLog.Printf("Listening on port %s", os.Getenv("PORT"))
	err := srv.ListenAndServe()
	h.ErrorLog.Fatal(err)
}

func (h *HiSpeed) checkDotEnv(path string) error {
	err := h.CreateFileIfNotExists(fmt.Sprintf("%s/.env", path))
	if err != nil {
		return err
	}
	return nil
}

func (h *HiSpeed) startLoggers() (*log.Logger, *log.Logger) {
	var infoLog *log.Logger
	var errorLog *log.Logger

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}

func (h *HiSpeed) createRenderer() {
	myRenderer := render.Render{
		Renderer: h.config.renderer,
		RootPath: h.RootPath,
		Port:     h.config.port,
		JetViews: h.JetViews,
		Session:  h.Session,
	}
	h.Render = &myRenderer
}

func (h *HiSpeed) BuildDSN() string {
	var dsn string //Made dsn a variable because different db's may/will be different then postgres...

	switch os.Getenv("DATABASE_TYPE") { //Of course we would get this below from the users env file in their app...
	case "postgres", "postgresql":
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=5",
			os.Getenv("DATABASE_HOST"),
			os.Getenv("DATABASE_PORT"),
			os.Getenv("DATABASE_USER"),
			os.Getenv("DATABASE_NAME"),
			os.Getenv("DATABASE_SSL_MODE"))

		//Builds a connection string for postgres and because some versions do not build with a password; need to support this... so use an if

		if os.Getenv("DATABASE_PASS") != "" { //If not blank then the user has supplied a password so get it and substitute
			dsn = fmt.Sprintf("%s password=%s", dsn, os.Getenv("DATABASE_PASS"))
		}

	default:

	}

	return dsn
}
