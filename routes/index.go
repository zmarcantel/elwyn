package routes

import (
	"net/http"
	"strconv"

	"github.com/go-martini/martini"

	"github.com/zmarcantel/elwyn/logging"
	"github.com/zmarcantel/elwyn/routes/web"
)

var errorChan chan error

func Initialize(lock chan error, logger *logging.Router, port int) (err error) {
	errorChan = lock

	m := martini.Classic()
	m.Map(logger.Web())

	RegisterDefaults(m)
	err = http.ListenAndServe(":"+strconv.Itoa(port), m)
	if err != nil {
		return
	}

	err = web.Initialize(errorChan)
	return
}

func RegisterDefaults(m *martini.ClassicMartini) {
	m.Get("/", web.Home)
}
