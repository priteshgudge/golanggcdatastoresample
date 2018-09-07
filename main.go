package main

import (
	"log"
	"os"
	"time"

	// Imports the Google Cloud Datastore client package.
	"cloud.google.com/go/datastore"
	"github.com/priteshgudge/golanggcdatastoresample/dao"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()

	// Set your Google Cloud Platform project ID.
	//projectID := "YOUR_PROJECT_ID"
	projectID := os.Getenv("GC_PROJECT_ID")

	// Creates a client.
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// // Creates a Customer instance.
	customer := dao.Customer{
		Name:         "Test Customer",
		PhoneNumbers: []string{"38827729388"},
		Created:      time.Now(),
	}
	// Saves the new entity.
	if key, err := customer.SaveToDB(ctx, client); err != nil {
		log.Fatalf("Failed To save Customer: %v", err)
	} else {
		log.Printf("Key saved: %v\n", key)
	}

	if customerObj, err := dao.GetCustomerFromDB(ctx, client, customer.UUID); err != nil {
		log.Fatalf("Failed To Get Customer: %v", err)
	} else {

		log.Printf("Got Customer Object: %v", customerObj)
	}

	if err := dao.ArchiveCustomer(ctx, client, customer.UUID); err != nil {
		log.Fatalf("Failed To Archive Customer: %v", err)
	} else {
		log.Printf("Archived Customer: %v", customer.UUID)
	}

}
