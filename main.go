package main

import (
	"log"
	"net/http"
    "go-cloud-ssl/db"
    "go-cloud-ssl/handlers"
	"os"
	"cloud.google.com/go/compute/metadata"
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

	config.Database.Driver = os.Getenv("driver")
	config.Database.InstanceConnectionName = os.Getenv("instance_connection_name")
	config.Database.User = os.Getenv("db_user")
	config.Database.Name = os.Getenv("db_name")
	config.Database.Private = os.Getenv("private")

		// Automatically fetch service account email from metadata server (Cloud Run)
	if metadata.OnGCE() {
		email, err := metadata.Email("default")
		if err != nil {
			log.Fatalf("Failed to get service account email: %v", err)
		}
		config.Database.User = email
		log.Printf("Using service account: %s", email)
	} else {
		// Fallback for local development
		config.Database.User = os.Getenv("db_user")
	}

}
