OutDir=./
protoc ./api/proto/message.proto --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=$OutDir -I .
protoc ./api/proto/message.proto --go_out=$OutDir

make build