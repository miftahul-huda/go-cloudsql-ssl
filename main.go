package main

import (
	"log"
	"net/http"
    "go-cloud-ssl/db"
    "go-cloud-ssl/handlers"

	"gopkg.in/yaml.v2"
	"os"
)



var config db.Config

func main() {
	loadConfig()
	db.InitDB(config)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.Index(w, r, config)
	})

	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		handlers.Create(w, r, config)
	})
	http.HandleFunc("/edit", func(w http.ResponseWriter, r *http.Request) {
		handlers.Edit(w, r, config)
	})
	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		handlers.Update(w, r, config)
	})
	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		handlers.Delete(w, r, config)
	})

	log.Println("Server running at http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}

func loadConfig() {
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("Cannot read config.yaml: ", err)
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatal("Cannot parse config.yaml: ", err)
	}
}
