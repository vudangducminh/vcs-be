curl -X GET "http://localhost:9200/_search" -H "Content-Type: application/json" -d '{
  "query": {
    "match_all": {}
  }
}'