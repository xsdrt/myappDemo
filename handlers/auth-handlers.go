package handlers

import "net/http"

func (h *Handlers) UserLogin(w http.ResponseWriter, r *http.Request) { //User login page...
	err := h.App.Render.Page(w, r, "login", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}

func (h *Handlers) PostUserLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm() //Grab the info from the request to make sure enough info to log a user in...
	if err != nil {
		w.Write([]byte(err.Error()))
		return // With a web application would do some more here... but this not a web app per se  but an app to make building web apps easier in Go...
	}

	email := r.Form.Get("email")       // ok , now how do we get a user and compare their password in the database?
	password := r.Form.Get("password") // References to make this work (a reason for comments, for those who don't think comments are necessary)added code
	// to handlers.go (Models data.Models to the  type struct); to init.hiSpeed.go (myHandlers.Models = app.Models); Hey... now have access to
	// the necessary functions... Know what this does but needed a ref to remember how/why/where to connect the bits :) ...
	user, err := h.Models.Users.GetByEmail(email)
	if err != nil {
		w.Write([]byte(err.Error()))
		return // With a web application would do some more here... but this not a web app per se  but an app to make building web apps easier in Go...
	}

	matches, err := user.PasswordMatches(password)
	if err != nil {
		w.Write([]byte("Error validating users password"))
		return
	}

	if !matches {
		w.Write([]byte("Invalid password!"))
		return
	}

	h.App.Session.Put(r.Context(), "userID", user.ID) // Ok , now a valid user so log in as password matches...

	http.Redirect(w, r, "/", http.StatusSeeOther) //Then take the valid user back to the Home Page...
}
