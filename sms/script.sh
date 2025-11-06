
curl -X GET "http://localhost:9200/server/_count?pretty=true" -H 'Content-Type: application/json' -d'{
    "query": {
        "bool": {
            "must": [
                {
                    "wildcard": {
                        "server_name": {
                            "value": "*123*"
                        }
                    }
                },
                {
                    "term": {
                        "status": "inactive"
                    }
                }
            ]
        }
    }
}'
