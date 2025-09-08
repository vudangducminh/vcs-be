curl -X POST "localhost:9200/server/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "size": 10000,
  "_source": ["server_id"],
  "query": {
    "match_all": { }
  }
}'
