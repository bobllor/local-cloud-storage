package api

import "net/http"

type HandlerMap map[string]func(http.ResponseWriter, *http.Request)

type Handler interface{}
