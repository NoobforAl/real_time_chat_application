# Document

Task Description: [link](./Task_dis.md)

Diagram MicroService in doc folder.  
Postman file in doc folder (http request only~).

For request to grpc and webSockets use this [link]("https://www.postman.com/science-astronomer-71562693/workspace/my-workspace/collection/668413c45a8d9a7d9fbe4817?action=share&creator=30256855"), this like is public postman file for test app ;)

## Requirements

- postman
- make
- protoc-compiler
- golang protoc pkg
- docker
- redis
- air (restart code run again if one go file changed, this tools for dev)

## How run ?

First setup `.env` file like `.env.example`.

### Run with Docker

Run:
> docker compose up -d

and run postman, use postman file in doc folder for request to api!

### Run local

First run docker-compose.yaml in Docker folder ( this docker-compose file for run database, you can install data base locally )

and run each service you want run it.

Note: command for build in Makefile.

## Run tests

> make tests
