package dao

import (
	"context"
	"log"
	"time"

	// Imports the Google Cloud Datastore client package.
	"cloud.google.com/go/datastore"
	uuid "github.com/satori/go.uuid"
)

var gcDBCustomerKind = "Customer"

type Customer struct {
	UUID         string    `datastore:"uuid"` // The UUID used in the datastore.
	Name         string    `datastore:"name"`
	PhoneNumbers []string  `datastore:"phoneNumbers"`
	Created      time.Time `datastore:"created"`
	Archived     bool      `datastore:"archived"`
	Remarks      string    `datastore:"remarks,noindex"`
}

//SaveToDB saves the customer object initialized to DB
func (c *Customer) SaveToDB(ctx context.Context, client *datastore.Client) (*datastore.Key, error) {

	//For Customer Object created First Time
	if c.UUID == "" {
		uuid := uuid.Must(uuid.NewV4())
		c.UUID = uuid.String()
	}

	//key := datastore.IncompleteKey("Customer", nil)
	// // Creates a Key instance.
	key := datastore.NameKey(gcDBCustomerKind, c.UUID, nil)

	responseKey, err := client.Put(ctx, key, c)
	if err != nil {
		log.Fatalf("Failed to save Customer: %v", err)
	}

	return responseKey, err
}

//ArchiveCustomer archives customer in atomic fashion
func ArchiveCustomer(ctx context.Context, client *datastore.Client, customerUUID string) error {

	// // Creates a Key instance.
	key := datastore.NameKey(gcDBCustomerKind, customerUUID, nil)
	// In a transaction load each Customer, set archived to true and store.
	_, err := client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var customerObj Customer
		if err := tx.Get(key, &customerObj); err != nil {
			return err
		}
		customerObj.Archived = true
		_, err := tx.Put(key, customerObj)
		return err
	})
	return err
}

//GetCustomerFromDB Gets an existing customer object from the DB
//Returns nil and error if not found
func GetCustomerFromDB(ctx context.Context, client *datastore.Client, UUID string) (*Customer, error) {
	key := datastore.NameKey(gcDBCustomerKind, UUID, nil)
	customerObj := Customer{}
	err := client.Get(ctx, key, &customerObj)

	if err != nil {
		log.Fatalf("Failed Get Customer: %v", err)
		return nil, err
	}

	return &customerObj, nil
}
