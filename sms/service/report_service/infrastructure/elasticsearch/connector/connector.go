package connector

import (
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

var Es *elasticsearch.Client
var err error
var isConnected = false

func IsConnected() bool {
	return isConnected
}

func Reconnect() error {
	Connect()
	if !isConnected {
		return fmt.Errorf("failed to reconnect to Elasticsearch")
	}
	return nil
}
func Connect() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://elasticsearch:9200",
		},
		APIKey: "",
	}

	Es, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Println("Error creating the client:", err)
		isConnected = false
		return
	}

	// Test the connection with a simple ping
	res, err := Es.Ping()
	if err != nil {
		log.Println("Error pinging Elasticsearch:", err)
		isConnected = false
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Elasticsearch ping failed: %s", res.Status())
		isConnected = false
		return
	}

	isConnected = true
	log.Println("Connected to Elasticsearch successfully")
}
