# stage build
FROM golang:1.22@sha256:a66eda637829ce891e9cf61ff1ee0edf544e1f6c5b0e666c7310dce231a66f28 as builder 

WORKDIR /app/

COPY . .

RUN apt-get update -y && \
    apt-get install -y curl make gcc protobuf-compiler && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

ENV export CGO_ENABLED=0
ENV export GOARCH=amd64
ENV export GOOS=linux

RUN make build_all

# auth api service
FROM ubuntu:24.04 as auth_api

WORKDIR /root/

COPY --from=builder /app/auth_api .

CMD [ "./auth_api" ]


# room api serive
FROM ubuntu:24.04 as room_api

WORKDIR /root/

COPY --from=builder /app/room_api .

CMD [ "./room_api" ]


# message api serive
FROM ubuntu:24.04 as message_api

WORKDIR /root/

COPY --from=builder /app/message_api .

CMD [ "./message_api" ]


# notifications api serive
FROM ubuntu:24.04 as notifications_api

WORKDIR /root/

COPY --from=builder /app/notifications_api .

CMD [ "./notifications_api" ]


# auth grpc serive
FROM ubuntu:24.04 as auth_api_grpc

WORKDIR /root/

COPY --from=builder /app/auth_api_grpc .

CMD [ "./auth_api_grpc" ]



# rooms task serive
FROM ubuntu:24.04 as rooms_task_server

WORKDIR /root/

COPY --from=builder /app/rooms_task_server .

CMD [ "./rooms_task_server" ]


# message task server
FROM ubuntu:24.04 as message_task_server

WORKDIR /root/

COPY --from=builder /app/message_task_server .

CMD [ "./message_task_server" ]


# notifications task server
FROM ubuntu:24.04 as  notifications_task_server

WORKDIR /root/

COPY --from=builder /app/notifications_task_server .

CMD [ "./notifications_task_server" ]

