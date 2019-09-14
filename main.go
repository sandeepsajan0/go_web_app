package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/go-redis/redis"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
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

func getComments(w http.ResponseWriter, r *http.Request){
	comments, err := client.LRange("comments", 0, 10).Result()
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Database Error"))
		return
	}
	templates.ExecuteTemplate(w, "index.html", comments)
}

func postComments(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	comment := r.PostForm.Get("comment")
	err := client.LPush("comments", comment).Err()
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Database Error"))
		return
	}
	http.Redirect(w, r, "/get_comments", 302)
}

func getLogin(w http.ResponseWriter, r *http.Request){
	templates.ExecuteTemplate(w, "login.html", nil)
}

func postLogin(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	hash, err := client.Get("user:" + username).Bytes()
	if err == redis.Nil{
		templates.ExecuteTemplate(w, "login.html", "Incorrect username")
		return
	} else if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil{
		templates.ExecuteTemplate(w, "login.html", "Incorrect password")
		return
	}
	session, _ := store.Get(r, "session")
	session.Values["username"] = username
	session.Save(r, w)
	http.Redirect(w, r, "/get_comments", 302)
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

func getRegister(w http.ResponseWriter, r *http.Request){
	templates.ExecuteTemplate(w, "register.html", nil)
}

func postRegister(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	e := client.Set("user:" + username, hash, 0).Err()
	if e != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Database Error"))
		return
	}
	http.Redirect(w, r, "/login", 302)
}

func authMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		session, _ := store.Get(r, "session")
		_, ok := session.Values["username"]
		if !ok{
			http.Redirect(w, r, "/login", 302)
		}
		handler.ServeHTTP(w, r)
	}
}

func main() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	store.Options = &sessions.Options{
		MaxAge:   60 * 1, 
		HttpOnly: true,}
	templates = template.Must(template.ParseGlob("templates/*.html"))
	r := mux.NewRouter()
	r.HandleFunc("/hello", handler)
	r.HandleFunc("/goodbye", byehandler)
	r.HandleFunc("/get_comments", authMiddleware(getComments)).Methods("GET")
	r.HandleFunc("/get_comments", authMiddleware(postComments)).Methods("POST")
	r.HandleFunc("/login", getLogin).Methods("GET")
	r.HandleFunc("/login", postLogin).Methods("POST")
	r.HandleFunc("/login/session", testLogin).Methods("GET")
	r.HandleFunc("/register", getRegister).Methods("GET")
	r.HandleFunc("/register", postRegister).Methods("POST")
	f := http.FileServer(http.Dir("static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", f))
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
