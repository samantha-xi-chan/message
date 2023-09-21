




#req=$(cat ./test10/session_req.json) ;

# time_asc 为 true表示 按时间升序
req='{
  "time_asc": true,
  "page_id": 1,
  "page_size": 2
}'
echo "$req" ; curl -X GET "http://192.168.31.7:18081/msg/api/v1r/log/session/task_1695268492571hglf" -d "$req"; echo ;

req='{
  "time_asc": false,
  "page_id": 1,
  "page_size": 2
}'
echo "$req" ; curl -X GET "http://192.168.31.7:18081/msg/api/v1r/log/session/task_1695268492571hglf"  -d "$req"; echo ;

