package web

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/gzip"
	"github.com/martini-contrib/render"

	"github.com/zmarcantel/elwyn/logging"
	"github.com/zmarcantel/elwyn/web/routes"
)

var errorChan chan error

func Initialize(lock chan error, logger *logging.Router, port int, directory string) {
	errorChan = lock
	logger.Banner("Starting Server")

	m := martini.Classic()
	m.Map(logger.Web())

	RegisterMiddleware(m, directory)
	RegisterDefaults(m)

	go func() {
		lock <- http.ListenAndServe(":"+strconv.Itoa(port), m)
	}()
}

func RegisterMiddleware(m *martini.ClassicMartini, directory string) {
	var webDir = filepath.Join(directory, "web")

	m.Use(gzip.All())

	m.Use(render.Renderer(render.Options{
		Directory:       filepath.Join(webDir, "views"),
		Extensions:      []string{".tmpl", ".html"},
		IndentJSON:      true,
		HTMLContentType: "text/html",
	}))

	m.Use(martini.Static(filepath.Join(webDir, "static/css"), martini.StaticOptions{
		Prefix: "/css",
	}))

	m.Use(martini.Static(filepath.Join(webDir, "static/js"), martini.StaticOptions{
		Prefix: "/js",
	}))
}

func RegisterDefaults(m *martini.ClassicMartini) {
	m.Get("/", routes.Home)
}
