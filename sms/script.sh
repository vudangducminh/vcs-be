
curl -X PUT "http://localhost:9200/server" -H 'Content-Type: application/json' -d'{
    "mappings": {
        "properties": {
            "ipv4": {
                "type": "text"
            },
            "uptime": {
                "type": "integer"
            },
            "last_updated_time": {
                "type": "long"
            },
            "status": {
                "type": "keyword"
            },
            "server_name": {
                "type": "text"
            },
            "created_time": {
                "type": "long"
            }
        }
    }
}'
