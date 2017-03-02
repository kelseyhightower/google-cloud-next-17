package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
)

func main() {
	ctx := context.Background()
	projectID := os.Getenv("PROJECT_ID")

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	topicName := "events"

	sub, err := client.CreateSubscription(context.Background(), "events", topic, 0, nil)
}
