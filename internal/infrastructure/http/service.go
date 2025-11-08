package http

import "net/http"

type HttpService interface {
	GetRouter() http.Handler
}
