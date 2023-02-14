package handlers

import (
	"myappDemo/data"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/xsdrt/hiSpeed"
)

type Handlers struct {
	App    *hiSpeed.HiSpeed
	Models data.Models //When create handlers should have already have models; ie. type data.Models...
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.Page(w, r, "home", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

func (h *Handlers) GoPage(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.GoPage(w, r, "home", nil) //Just use the home.page.tmpl already in the views folder for this one
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

func (h *Handlers) JetPage(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.JetPage(w, r, "jet-template", nil, nil) //create a jet-template in the views folder
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

func (h *Handlers) SessionTest(w http.ResponseWriter, r *http.Request) {
	myData := "Times" //Following is just a test to make sure this works...

	h.App.Session.Put(r.Context(), "Fun", myData) //Context is "Fun" and the myData "Times" is; well the data of course...

	myValue := h.App.Session.GetString(r.Context(), "Fun")

	vars := make(jet.VarMap)
	vars.Set("Fun", myValue)

	err := h.App.Render.JetPage(w, r, "sessions", vars, nil) //Should call the sessions.jet from the views folder...
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}
