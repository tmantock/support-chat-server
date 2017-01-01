package main

import (
	"fmt"
	"net/http"

	r "github.com/dancannon/gorethink"
	"github.com/gorilla/websocket"
)

//Handler defines a new Handler type
type Handler func(*Client, interface{})

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

//Router implemenst a new Roueter type
type Router struct {
	rules   map[string]Handler
	session *r.Session
}

//FindHandler function finds a handler
func (rt *Router) FindHandler(msgName string) (Handler, bool) {
	handler, found := rt.rules[msgName]
	return handler, found
}

//ServeHTTP implements websockets
func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, error := upgrader.Upgrade(w, r, nil)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, error.Error())
		return
	}

	client := NewClient(socket, rt.FindHandler, rt.session)
	defer client.Close()
	go client.Write()
	client.Read()
}

//Handle implements a handle method
func (rt *Router) Handle(msgName string, handler Handler) {
	rt.rules[msgName] = handler
}

//NewRouter allows for new routers to be constructed
func NewRouter(session *r.Session) *Router {
	return &Router{
		rules:   make(map[string]Handler),
		session: session,
	}
}
