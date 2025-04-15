package main

import (
	"context"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

func main() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
		APIKey: "",
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// API Key should have cluster monitoring rights
	infores, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	fmt.Println(infores)

	searchResp, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("books"),
		es.Search.WithQuery(""),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)

	fmt.Println(searchResp, err)
}
