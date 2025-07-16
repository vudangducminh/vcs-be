curl -X POST "http://localhost:9200/server/_search" -H "Content-Type: application/json" -d '{
    "query": {
        "match_all": {}
    }
}'