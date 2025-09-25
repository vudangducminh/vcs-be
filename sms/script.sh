curl -X GET "http://localhost:9200/server/_search?pretty" -H 'Content-Type: application/json' -d'{
    "query": {
        "term": {
            "status": "inactive"
        }
    }
}'