package hiSpeed

import "net/http"

func (h *HiSpeed) SessionLoad(next http.Handler) http.Handler { //This func handles saving and loading the session on every request...
	h.InfoLog.Println("SessionLoad called") //This is just to test and make sure the middleware is working...
	return h.Session.LoadAndSave(next)
}
