curl -X DELETE "http://localhost9200/server/_delete_by_query?pretty=true" -H "Content-Type: application/json" -d '{
   "query": {
	 "match_all": {}
   }
 }'