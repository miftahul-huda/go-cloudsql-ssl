package main

import (
	"log"
	"net/http"
    "go-cloud-ssl/db"
    "go-cloud-ssl/handlers"
	"os"
	//"cloud.google.com/go/compute/metadata"
    "context"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
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

func loadConfigOld() {

	config.Database.Driver = os.Getenv("driver")
	config.Database.InstanceConnectionName = os.Getenv("instance_connection_name")
	config.Database.User = os.Getenv("db_user")
	config.Database.Name = os.Getenv("db_name")
	config.Database.Private = os.Getenv("private")

		// Automatically fetch service account email from metadata server (Cloud Run)
		/*
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
		*/

}


func loadConfig() {
	ctx := context.Background()

	driver, err := accessSecret(ctx, "db_driver")
	if err != nil {
		log.Fatalf("Failed to load driver: %v", err)
	}
	config.Database.Driver = driver

	instance, err := accessSecret(ctx, "db_instance")
	if err != nil {
		log.Fatalf("Failed to load db_instance: %v", err)
	}
	config.Database.InstanceConnectionName = instance

	user, err := accessSecret(ctx, "db_user")
	if err != nil {
		log.Fatalf("Failed to load db_user: %v", err)
	}
	config.Database.User = user

	dbName, err := accessSecret(ctx, "db_name")
	if err != nil {
		log.Fatalf("Failed to load db_name: %v", err)
	}
	config.Database.Name = dbName

	/*
	private, err := accessSecret(ctx, "private")
	if err != nil {
		log.Fatalf("Failed to load private: %v", err)
	}
	config.Database.Private = private
	*/
}

func accessSecret(ctx context.Context, name string) (string, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	// Replace PROJECT_ID with your actual GCP project ID
	PROJECT_ID := os.Getenv("PROJECT_ID")
	secretName := "projects/" + PROJECT_ID + "/secrets/" + name + "/versions/latest"

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretName,
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", err
	}

	return string(result.Payload.Data), nil
}