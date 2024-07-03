GO := go

all: build_all 

build_all: generate build_task_server build_grpc_sevices build_api_service
	@echo "build done!"

run:  #TODO your option

generate:
	$(GO) generate ./...

run_docker_compose: create_vendor
	docker compose -f "docker-compose.yaml" down 
	docker compose -f "docker-compose.yaml" up -d --build 

create_vendor:
	$(GO) mod vendor

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

tests:
	@echo "run docker compose test database"
	@docker compose -f "docker/docker-compose.yaml" down
	@docker compose -f "docker/docker-compose.yaml" up -d --build 
	
	@echo "run test database & store"
	$(GO) test -timeout 30s -tags app_debug_mod -coverprofile=/tmp/vscode-gobr19vy/go-code-cover github.com/NoobforAl/real_time_chat_application/src/database

	@echo "restart database"
	@docker compose -f "docker/docker-compose.yaml" down
	@docker compose -f "docker/docker-compose.yaml" up -d --build 

	@echo "run grpc auth test"
	$(GO) test -timeout 30s -tags app_debug_mod -coverprofile=/tmp/vscode-gobr19vy/go-code-cover github.com/NoobforAl/real_time_chat_application/src/grpc/auth/tests

	@echo "restart database"
	@docker compose -f "docker/docker-compose.yaml" down
	@docker compose -f "docker/docker-compose.yaml" up -d --build 

	@echo "run test message task"
	$(GO) test -timeout 30s -tags app_debug_mod -coverprofile=/tmp/vscode-gobr19vy/go-code-cover github.com/NoobforAl/real_time_chat_application/src/tasks/messages/tasks_message/tests

	@echo "restart database"
	@docker compose -f "docker/docker-compose.yaml" down
	@docker compose -f "docker/docker-compose.yaml" up -d --build 

	@echo "run auth service test"
	$(GO) test -timeout 30s -tags app_debug_mod -coverprofile=/tmp/vscode-gobr19vy/go-code-cover github.com/NoobforAl/real_time_chat_application/src/services/auth/http/v1/router

	@echo "restart database"
	@docker compose -f "docker/docker-compose.yaml" down
	@docker compose -f "docker/docker-compose.yaml" up -d --build 

	@echo "run room service test"
	$(GO) test -timeout 30s -tags app_debug_mod -coverprofile=/tmp/vscode-gobr19vy/go-code-cover github.com/NoobforAl/real_time_chat_application/src/services/rooms/http/v1/router

	@echo "down databases"
	@docker compose -f "docker/docker-compose.yaml" down
	@echo "see TODOs"
