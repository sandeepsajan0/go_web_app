package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/go-redis/redis"
	"github.com/gorilla/sessions"
	"html/template"
)

var client *redis.Client
var templates *template.Template
var store = sessions.NewCookieStore([]byte("top-s3cr3t"))

func handler(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Hello world")
}

func byehandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Bye Bye world")
}

func getIndexHandler(w http.ResponseWriter, r *http.Request){
	comments, err := client.LRange("comments", 0, 10).Result()
	if err != nil{
		return
	}
	templates.ExecuteTemplate(w, "index.html", comments)
}

func postIndexHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	comment := r.PostForm.Get("comment")
	client.LPush("comments", comment)
	http.Redirect(w, r, "/get_comments", 302)
}

func getLoginHandler(w http.ResponseWriter, r *http.Request){
	templates.ExecuteTemplate(w, "login.html", nil)
}

func postLoginHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	username := r.PostForm.Get("username")
	session, _ := store.Get(r, "session")
	session.Values["username"] = username
	session.Save(r, w)
}

func testLogin(w http.ResponseWriter, r *http.Request){
	session, _ := store.Get(r, "session")
	untyped, ok := session.Values["username"]
	if !ok{
		return
	}
	username, ok := untyped.(string)
	if !ok{
		return
	}
	w.Write([]byte(username))
}

func main() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	templates = template.Must(template.ParseGlob("templates/*.html"))
	r := mux.NewRouter()
	r.HandleFunc("/hello", handler)
	r.HandleFunc("/goodbye", byehandler)
	r.HandleFunc("/get_comments", getIndexHandler).Methods("GET")
	r.HandleFunc("/get_comments", postIndexHandler).Methods("POST")
	r.HandleFunc("/login", getLoginHandler).Methods("GET")
	r.HandleFunc("/login", postLoginHandler).Methods("POST")
	r.HandleFunc("/login/session", testLogin).Methods("GET")
	f := http.FileServer(http.Dir("static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", f))
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
