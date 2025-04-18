http http://localhost:9200/book/_search <<EOF
{
  "query": {
    "match_all": {}
  },
  "sort": [
    {
      "title": {
        "order": "asc"
      }
    }
  ]
}