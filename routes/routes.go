package routes

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"../models"
	"html/template"
	"../sessions"
	"../middleware"
)

var templates *template.Template

func handler(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Hello world")
}

func byehandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Bye Bye world")
}

func getComments(w http.ResponseWriter, r *http.Request){
	comments, err := models.GetComments()
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
	err := models.PostComments(comment)
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
	err := models.LoginUser(username, password)
	if err != nil{
		switch err{
		case models.Errusername:
			templates.ExecuteTemplate(w, "login.html", "Incorrect username")
		case models.Errpassword:
			templates.ExecuteTemplate(w, "login.html", "Incorrect password")
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		}
		return
	}
	session, _ := session.Store.Get(r, "session")
	session.Values["username"] = username
	session.Save(r, w)
	http.Redirect(w, r, "/get_comments", 302)
}

func testLogin(w http.ResponseWriter, r *http.Request){
	session, _ := session.Store.Get(r, "session")
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
	err := models.RegisterUser(username, password)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Database Error"))
		return
	}
	http.Redirect(w, r, "/login", 302)
}


func GetRoutes() *mux.Router{
	templates = template.Must(template.ParseGlob("templates/*.html"))
	r := mux.NewRouter()
	r.HandleFunc("/hello", handler)
	r.HandleFunc("/goodbye", byehandler)
	r.HandleFunc("/get_comments", middleware.AuthMiddleware(getComments)).Methods("GET")
	r.HandleFunc("/get_comments", middleware.AuthMiddleware(postComments)).Methods("POST")
	r.HandleFunc("/login", getLogin).Methods("GET")
	r.HandleFunc("/login", postLogin).Methods("POST")
	r.HandleFunc("/login/session", testLogin).Methods("GET")
	r.HandleFunc("/register", getRegister).Methods("GET")
	r.HandleFunc("/register", postRegister).Methods("POST")
	f := http.FileServer(http.Dir("static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", f))
	return r
}