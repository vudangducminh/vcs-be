curl -X POST "localhost:9200/server/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "query": {
    "range": {
      "uptime": {
        "gt": 0
      }
    }
  }
}'
