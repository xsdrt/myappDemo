package handlers

import (
	"net/http"

	"github.com/xsdrt/hiSpeed"
)

type Handlers struct {
	App *hiSpeed.HiSpeed
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
	myData := "Fun"

	h.App.Session.Put(r.Context(), "Times", myData)

	err := h.App.Render.JetPage(w, r, "sessions", nil, nil) //Should call the sessions.jet from the views folder...
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}
