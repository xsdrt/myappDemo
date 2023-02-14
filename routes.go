package main

import (
	"fmt"
	"myappDemo/data"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (a *application) routes() *chi.Mux {
	// middleware must come before any routes..

	// add routes here...
	a.App.Routes.Get("/", a.Handlers.Home)
	a.App.Routes.Get("/go-page", a.Handlers.GoPage) //Still need to develop these handlers yet...
	a.App.Routes.Get("/jet-page", a.Handlers.JetPage)
	a.App.Routes.Get("/sessions", a.Handlers.SessionTest)

	a.App.Routes.Get("/users/login", a.Handlers.UserLogin)
	//a.App.Routes.Post("/users/login", a.Handlers.PostUserLogin) //commented out for now as it does not exist yet(TODO)...

	//Test routes for testing the user models with postgresql using inline funcs...
	a.App.Routes.Get("/create-user", func(w http.ResponseWriter, r *http.Request) { //Test route to add a user...
		u := data.User{
			FirstName: "Michael",
			LastName:  "Redinger",
			Email:     "somebody@there.com",
			Active:    1,
			Password:  "password",
		}

		id, err := a.Models.Users.Insert(u) //This should insert
		if err != nil {
			a.App.ErrorLog.Println(err) //if an error will just print a blank "page" in th web browser...
			return
		}

		fmt.Fprintf(w, "%d: %s", id, u.FirstName) //If this works should print out id and first name from the database after the insert
	})

	a.App.Routes.Get("/get-all-users", func(w http.ResponseWriter, r *http.Request) {
		users, err := a.Models.Users.GetAll()
		if err != nil {
			a.App.ErrorLog.Println(err) //if an error will just print a blank "page" in the web browser...
			return
		}
		for _, x := range users {
			fmt.Fprint(w, x.LastName)
		}
	})

	a.App.Routes.Get("/get-user/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		u, err := a.Models.Users.Get(id)
		if err != nil {
			a.App.ErrorLog.Println(err) //if an error will just print a blank "page" in th web browser...
			return
		}

		fmt.Fprintf(w, "%s %s %s", u.FirstName, u.LastName, u.Email)
	})

	a.App.Routes.Get("/update-user/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		u, err := a.Models.Users.Get(id)
		if err != nil {
			a.App.ErrorLog.Println(err) //if an error will just print a blank "page" in th web browser...
			return
		}

		u.LastName = a.App.RandomString(10) //Change the name to a random generated string and see if the update works...
		err = u.Update(*u)                  //Hand it a pointer to u...
		if err != nil {
			a.App.ErrorLog.Println(err) //if an error will just print a blank "page" in th web browser...
			return
		}

		fmt.Fprintf(w, "updated last name to  %s", u.LastName)

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
