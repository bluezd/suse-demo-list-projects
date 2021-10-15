package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	//"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

type JSONTime struct {
	time.Time
}

func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", t.Format("Mon Jan _2"))
	return []byte(stamp), nil
}

type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Repository  string `json:"repository"`
	Twitter     string `json:"twitter"`
	Website     string `json:"website"`
	Description string `json:"description"`
}

var projects []Project

func GetProject(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range projects {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Project{})
}
func GetProjectsList(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(projects)
}

func CreateProject(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var project Project
	_ = json.NewDecoder(r.Body).Decode(&project)
	project.ID = params["id"]
	for _, item := range projects {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(map[string]string{"error": "duplicated id"})
			return
		}
	}
	projects = append(projects, project)
	json.NewEncoder(w).Encode(projects)
}

func DeleteProject(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range projects {
		if item.ID == params["id"] {
			projects = append(projects[:index], projects[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(projects)
}

func main() {
	var log = logrus.New()
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	}
	log.Out = os.Stdout

	router := mux.NewRouter()

	file, _ := ioutil.ReadFile("projects-contents.json")
	_ = json.Unmarshal([]byte(file), &projects)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})
	router.HandleFunc("/projects", GetProjectsList).Methods("GET")
	router.HandleFunc("/projects/{id}", GetProject).Methods("GET")
	router.HandleFunc("/projects/{id}", CreateProject).Methods("POST")
	router.HandleFunc("/projects/{id}", DeleteProject).Methods("DELETE")

	/*
		c := cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
		})
	*/

	var handler http.Handler = router
	handler = &logHandler{log: log, next: handler} // add logging

	log.Infof("starting backend on 0.0.0.0:8001")
	//log.Fatal(http.ListenAndServe(":8001", c.Handler(router)))
	log.Fatal(http.ListenAndServe(":8001", handler))
}
