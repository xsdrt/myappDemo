package handlers

import "net/http"

func (h *Handlers) UserLogin(w http.ResponseWriter, r *http.Request) { //User login page...
	err := h.App.Render.Page(w, r, "login", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}
