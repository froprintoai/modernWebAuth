package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"

	"github.com/froprintoai/modernWeb/loglog"
	"github.com/julienschmidt/httprouter"
)

type configuration struct {
	Path              string
	Path_without_port string
	Gmail             string
	Password          string
}

//Tech33 Parse template once
var signup_template = template.Must(template.ParseFiles("templates/loginSignup.html"))
var conf configuration

func main() {
	//Tech3 configure app from the configuration file
	file, err := os.Open("conf.json")
	if err != nil {
		loglog.LogWTF("Cannot Open conf file", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	if err != nil {
		loglog.LogWTF("Cannot decode json conf file", err)
	}

	//Tech9 use httprouter for fast and flexible routing
	mux := httprouter.New()
	mux.GET("/", home)
	mux.POST("/signup", signup)
	mux.POST("/login", login)
	mux.GET("/confirm/:hashed_email/:activation_code", confirm)
	server := http.Server{
		Addr:    conf.Path,
		Handler: mux,
	}
	server.ListenAndServe()
}
