package transhttp

import "net/http"

type Route struct {
	Name string
	Path string
	Handler http.Handler
}

type Routes []Route

func AddRoutes(routes Routes, server *http.Server)  {

}