package main //Added a vendor folder so VSCode will not take forever to clear an error thats not there (Ie.. another func that is not immediately recognized from the other pkg(s)..
import (
	"myappDemo/handlers"

	"github.com/xsdrt/hiSpeed"
)

type application struct {
	App      *hiSpeed.HiSpeed
	Handlers *handlers.Handlers
}

func main() {
	h := initApplication() //After init the app store value in a variable h (for hiSpeed)...
	h.App.ListenAndServe() // This should start the web server...
}
