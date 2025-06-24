package elastic

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

var Es *elasticsearch.Client
var err error

func ConnectToEs() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://elasticsearch:9200",
		},
		APIKey: "",
	}

	Es, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
}
