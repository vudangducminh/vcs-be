
curl -X POST "http://localhost:9200/server/_search?pretty" -H 'Content-Type: application/json' -d'{
    "query": {
        "match": {
            "_id": "cedc893d43608ab14d8a73bce3f8804d1a389735726f2f986cf2ee06adbd026a"
        }
    }
}'
