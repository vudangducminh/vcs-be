curl -X POST "http://localhost:9200/server/_count?pretty" -H 'Content-Type: application/json' -d'{
    "query": {
        "bool": {
            "must": [
                {
                    "wildcard": {
                        "ipv4": {
                            "value": "*23*"
                        }
                    }
                },
                {
                    "term": {
                        "status": "active"
                    }
                }
            ]
        }
    }
}'