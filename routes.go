package main

import (
	"fmt"
	"myappDemo/data"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *application) routes() *chi.Mux {
	// middleware must come before any routes..

	// add routes here...
	a.App.Routes.Get("/", a.Handlers.Home)
	a.App.Routes.Get("/go-page", a.Handlers.GoPage) //Still need to develop these handlers yet...
	a.App.Routes.Get("/jet-page", a.Handlers.JetPage)
	a.App.Routes.Get("/sessions", a.Handlers.SessionTest)

	a.App.Routes.Get("/create-user", func(w http.ResponseWriter, r *http.Request) {
		u := data.User{
			FirstName: "Michael",
			LastName:  "Redinger",
			Email:     "somebody@there.com",
			Active:    1,
			Password:  "password",
		}

		id, err := a.Models.Users.Insert(u) //This should insert
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}

		fmt.Fprintf(w, "%d: %s", id, u.FirstName) //If this works should print out id and first name from the database after the insert
	})

	/* a.App.Routes.Get("/test-database", func(w http.ResponseWriter, r *http.Request) { //Temp route (inline func)...  //commented out inline test; saved for reference
		query := "select id, first_name from users where id = 1" //Query a temp table (need to create)...
		row := a.App.DB.Pool.QueryRowContext(r.Context(), query) //Just query at most one row to test for now...

		var id int
		var name string
		err := row.Scan(&id, &name)
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}

		fmt.Fprintf(w, "%d %s", id, name) //Write to the browser window to show it works..
	}) */

	// a.App.Routes.Get("/jet", func(w http.ResponseWriter, r *http.Request) {  	//commented out was a inline test
	// 	a.App.Render.JetPage(w, r, "testjet", nil, nil)
	// })

	// static routes
	fileServer := http.FileServer(http.Dir("./public"))
	a.App.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return a.App.Routes

}
