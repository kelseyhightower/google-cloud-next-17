package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"google.golang.org/api/option"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
)

type Event struct {
	ID        string
	Message   string
	Region    string
	Timestamp time.Time
}

func main() {
	ctx := context.Background()

	databaseID := os.Getenv("DATABASE_ID")
	if databaseID == "" {
		log.Fatal("DATABASE_ID must be non-empty")
	}

	client, err := spanner.NewClient(ctx, databaseID+"/databases/example",
		option.WithServiceAccountFile("/var/run/secret/cloud.google.com/service-account.json"))
	if err != nil {
		log.Fatalf("Failed to create client %v", err)
	}
	defer client.Close()

	data := struct {
		PodName string
		Region  string
	}{
		os.Getenv("POD_NAME"),
		os.Getenv("REGION"),
	}

	t, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		event := Event{
			ID:        uuid.New().String(),
			Message:   "Website access",
			Region:    os.Getenv("REGION"),
			Timestamp: time.Now(),
		}

		m, err := spanner.InsertStruct("Event", event)
		if err != nil {
			log.Println(err)
		}
		_, err = client.Apply(ctx, []*spanner.Mutation{m})
		if err != nil {
			log.Println(err)
		}

		err = t.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
