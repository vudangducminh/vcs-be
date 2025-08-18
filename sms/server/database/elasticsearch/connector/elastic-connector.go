package elastic

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

var Es *elasticsearch.Client
var err error
var isConnected = false

func IsConnected() bool {
	return isConnected
}
func ConnectToEs() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
		APIKey: "",
	}

	Es, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Println("Error creating the client:", err)
	}
	isConnected = true
	log.Println("Connected to Elasticsearch successfully")
}
