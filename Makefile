GO := go
PATH_GO_FILE := ./src/ # file you winat run it
RUN_TIME_RELOADER := air

all: run 

run: generate
	$(RUN_TIME_RELOADER)

build:
	$(GO) build -o app $(PATH_GO_FILE)

test:
	$(GO) test ./...

generate:
	$(GO) generate ./...

# TODO: build type