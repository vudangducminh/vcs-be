curl -X POST "localhost:9200/server/_search" -H 'Content-Type: application/json' -d'
{
  "size": 0,
  "query": {
    "term": {
      "status": "active"
    }
  },
  "aggs": {
    "total_created_time": {
      "sum": {
        "field": "uptime"
      }
    }
  }
}'