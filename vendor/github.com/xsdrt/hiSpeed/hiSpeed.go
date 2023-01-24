package hiSpeed

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/xsdrt/hiSpeed/render"
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
	JetViews *jet.Set
	config   config //This will not be exported as there is no reason any app that imports HiSpeed should have access to the config...
}

// This type will not be exported but will hold all the config values for this package...
type config struct {
	port     string
	renderer string
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
	h.InfoLog = infoLog
	h.ErrorLog = errorLog
	h.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	h.Version = version
	h.RootPath = rootPath
	h.Routes = h.routes().(*chi.Mux) //Cast to a * of chi.Mux from a htttphandler

	h.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
	}

	var views = jet.NewSet(
		jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rootPath)),
		jet.InDevelopmentMode(),
	)
	h.JetViews = views

	h.CreateRenderer()

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
	}
	h.Render = &myRenderer
}
