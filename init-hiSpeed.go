package main

import (
	"log"
	"myappDemo/data"
	"myappDemo/handlers"
	"myappDemo/middleware"
	"os"

	"github.com/xsdrt/hiSpeed"
)

func initApplication() *application {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	//init hiSpeed
	his := &hiSpeed.HiSpeed{}
	err = his.New(path)
	if err != nil {
		log.Fatal(err)
	}

	his.AppName = "myappDemo"

	myMiddleware := &middleware.Middleware{
		App: his,
	}

	// his.InfoLog.Println("Debug is set to", his.Debug) commented out do not need , but saved for a reference...

	myHandlers := &handlers.Handlers{
		App: his,
	}

	app := &application{
		App:        his,
		Handlers:   myHandlers,
		Middleware: myMiddleware,
	}

	app.App.Routes = app.routes()

	app.Models = data.New(app.App.DB.Pool) //This initializes the models from the call from main.go...
	myHandlers.Models = app.Models
	app.Middleware.Models = app.Models

	return app
}
