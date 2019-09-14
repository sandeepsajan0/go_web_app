package middleware

import (
	"net/http"
	"../sessions"
)

func AuthMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		session, _ := session.Store.Get(r, "session")
		_, ok := session.Values["username"]
		if !ok{
			http.Redirect(w, r, "/login", 302)
		}
		handler.ServeHTTP(w, r)
	}
}