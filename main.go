package main

import (
	"net/http"
	"./models"
	"./sessions"
	"github.com/gorilla/sessions"
	"./routes"
)

func main() {
	models.Init()
	session.Store.Options = &sessions.Options{
		MaxAge:   60 * 1, 
		HttpOnly: true,}
	r := routes.GetRoutes()
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
