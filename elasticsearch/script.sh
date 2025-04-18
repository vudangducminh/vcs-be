http http://localhost:9200/book/_search <<EOF
{
  "query": {
    "range": {
      "rating": {
        "gte": 0.0
      }
    }
  }
}