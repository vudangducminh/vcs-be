# http http://localhost:9200/book/_analyze field=book.title text="Elasticsearch: The Definitive Guide"
http http://localhost:9200/book/_doc <<EOF
{
  "title": "hehehehaw",
  "author": "CR",
  "published_date": "2015-02-01",
  "rating": 7.27
}