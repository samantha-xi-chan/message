




#req=$(cat ./test10/session_req.json) ;


req='{
  "time_asc": true,
  "page_id": 1,
  "page_size": 2
}'

echo "$req" ; curl -X GET "http://192.168.31.117:18081/api/v1r/session/task_1695206325647zrfr"  -H "request-id: $RANDOM" -d "$req"

