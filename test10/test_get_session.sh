


#req=$(cat ./test10/session_req.json) ;

TaskId=task_1695701155504scsh

# echo "$req" ; curl -X GET "http://192.168.34.8/msg/api/v1/log/session/$TaskId" -d "$req"; echo ;


# time_asc 为 true表示 按时间升序
echo "$req" ; curl -X GET "http://message/msg/api/v2/log/session/$TaskId?time_asc=true&page_id=1&page_size=1" -d "$req"; echo ;
echo "$req" ; curl -X GET "http://message/msg/api/v2/log/session/$TaskId?time_asc=false&page_id=1&page_size=5" -d "$req"; echo ;

