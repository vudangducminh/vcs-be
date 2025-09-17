curl -X PUT http://localhost:9200/server -H 'Content-Type: application/json' -d'
{
  "mappings": {
    "properties": {
      "server_id": { "type": "keyword" },
      "server_name": { "type": "text" },
      "ipv4": { "type": "text" },
      "status": { "type": "keyword" },
      "uptime": { "type": "integer" },
      "created_time": { "type": "integer" },
      "last_updated_time": { "type": "integer" }
    }
  }
}
'