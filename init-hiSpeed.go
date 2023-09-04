package main

import (
	"log"
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
	his.Debug = true

	app := &application{
		App: his,
	}

	return app

}
