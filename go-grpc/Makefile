
.PHONY: generate-chat-proto
generate-chat-proto:
	protoc -I=. \
	    --go_out . \
	    --go_opt paths=source_relative \
	    --go-grpc_out . --go-grpc_opt paths=source_relative \
	    helloworld/chat/chat.proto

.PHONY: generate-route-proto
generate-route-proto:
	protoc -I=. \
	    --go_out . \
	    --go_opt paths=source_relative \
	    --go-grpc_out . --go-grpc_opt paths=source_relative \
	    route_guide/routeguide/route_guide.proto