curl -X POST "localhost:9200/server/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "size": 10,
  "query": {
    "script": {
      "script": {
        "source": "doc[\"uptime\"].size() > 0 && doc[\"uptime\"][0] > 0",
        "lang": "painless"
      }
    }
  }
}'