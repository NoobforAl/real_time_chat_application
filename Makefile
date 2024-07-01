GO := go

all: build_all 

build_all: generate build_task_server build_grpc_sevices build_api_service
	@echo "build done!"

run:  #TODO your option

tests:
	$(GO) test ./...

generate:
	$(GO) generate ./...


build_api_service: auth_api room_api message_api notifications_api
	@echo "done build server api app"

auth_api:
	@echo "start build auth_api"
	$(GO) build -mod=vendor -o auth_api  src/services/auth/http/v1/main.go
	@echo "start finish auth_api"

room_api:
	@echo "start build room_api"
	$(GO) build -mod=vendor -o room_api  src/services/rooms/http/v1/main.go 
	@echo "start finish room_api"

message_api:
	@echo "start build message_api"
	$(GO) build -mod=vendor -o message_api  src/services/messages/http/v1/main.go
	@echo "start finish message_api"

notifications_api:
	@echo "start build notifications_api"
	$(GO) build -mod=vendor -o notifications_api  src/services/notifications/http/v1/main.go
	@echo "start finish notifications_api"


build_grpc_sevices: auth_api_grpc
	@echo "done build grpc serivce"

auth_api_grpc:
	@echo "start build auth_api_grpc"
	$(GO) build -mod=vendor -o auth_api_grpc  src/services/auth/grpc/main.go 
	@echo "start finish auth_api_grpc"


build_task_server: rooms_task_server message_task_server notifications_task_server
	@echo "done build server tasks app"

rooms_task_server:
	@echo "start build rooms_task_server"
	$(GO) build -mod=vendor -o rooms_task_server  src/tasks/rooms/rooms.go 
	@echo "start finish rooms_task_server"

message_task_server:
	@echo "start build message_task_server"
	$(GO) build -mod=vendor -o message_task_server src/tasks/messages/message.go 
	@echo "start finish message_task_server"

notifications_task_server:
	@echo "start build notifications_task_server"
	$(GO) build -mod=vendor -o notifications_task_server src/tasks/notifications/notifications.go 
	@echo "start finish notifications_task_server"
