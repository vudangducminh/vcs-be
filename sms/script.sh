
curl -X POST "http://localhost:9200/server/_search?pretty" -H 'Content-Type: application/json' -d'{
    "query": {
        "match": {
            "server_name": "8256"
        }
    }
}'
