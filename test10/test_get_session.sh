

#req=$(cat ./test10/session_req.json) ;

TaskId=task_1695622594262bavn

# time_asc 为 true表示 按时间升序
req='{
  "time_asc": true,
  "page_id": 1,
  "page_size": 1
}'
# echo "$req" ; curl -X GET "http://192.168.31.6:18081/msg/api/v1r/log/session/task_1695614017280avqz" -d "$req"; echo ;
#echo "$req" ; curl -X GET "http://192.168.31.7/msg/api/v1/log/session/$TaskId" -d "$req"; echo ;
echo "$req" ; curl -X GET "http://message7/msg/api/v1/log/session/$TaskId" -d "$req"; echo ;

req='{
  "time_asc": false,
  "page_id": 1,
  "page_size": 2
}'
# echo "$req" ; curl -X GET "http://192.168.31.6:18081/msg/api/v1r/log/session/task_1695622594262bavn"  -d "$req"; echo ;
echo "$req" ; curl -X GET "http://message7/msg/api/v1/log/session/$TaskId"  -d "$req"; echo ;

