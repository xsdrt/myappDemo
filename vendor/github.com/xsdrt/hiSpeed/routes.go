package hiSpeed

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (h *HiSpeed) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	if h.Debug { //If in Debug mode log to console...
		mux.Use(middleware.Logger)
	}
	mux.Use(middleware.Recoverer)

	// mux.Get("/", func(w http.ResponseWriter, r *http.Request) { //Test route
	// 	fmt.Fprint(w, "Welcome to HiSpeed") //Should print to the web page the message just to make sure everything is working
	// })

	return mux
}
