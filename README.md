# Vigilant

A self-hosted, lightweight, and simple event monitoring tool written in Go and gRPC with Kafka as the message broker.


## Setup

1. Clone the repository.
2. Make sure you have Go installed.
    - If you don't have Go installed, you can download it from [here](https://golang.org/dl/).
    - You can check if you have Go installed by running `go version`.
    - You can check if Go is in your PATH by running `go env`.
3. Install gRPC.
    - You can install gRPC by running `go get -u google.golang.org/grpc`.
    - You can install the gRPC tools by running `go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc`.
4. Install the dependencies.
    - You can install the dependencies by running `go mod download`.
5. Make sure you have Docker installed.
    - If you don't have Docker installed, you can download it from [here](https://docs.docker.com/get-docker/).
    - You can check if you have Docker installed by running `docker --version`.
    - Newer versions of Docker come with Docker Compose.
6. Run the docker-compose file.
    - You can run the docker-compose file by running `docker compose up`.
    - This file will get Kafka up and running.
    - You can check if Kafka is running by running `docker ps`. You should see a container with the name `kafka` running.
7. Run the Go server.
    - You can run the Go server by running `go run main.go`.
    - The server will start on port `50051`.


🎉 You're all set! 🎉


### (Optional) Setting up a client

1. Install the Wails CLI.
    - You can install the wails CLI by running `go install github.com/wailsapp/wails/v2/cmd/wails@latest`.
    - You can check if you have the wails CLI installed by running `wails --version`.
    - If you encounter any issues, you can check the [Wails documentation](https://wails.io/docs/gettingstarted/installation).
2. Run the Wails app.
    - You can run the Wails app by running `wails dev` in the `cmd/desktop` directory.



## Running

Repeat step 6 and 7 from the setup section to run the server.

To run the client, repeat step 2 from the setup section.

### Deploying to production

Steps to take to deal with databases, Kafka, and security.

To do


## How it works

1. Events can be logged through gRPC. There is a `.proto` file that defines the service and the message that can be sent over gRPC.
   - Check the [Logging](#logging) section for more information on how to log events.
2. When events are sent over gRPC, they are sent to the server. The server takes the events and has Kafka produce them to a topic.
3. There is a concurrent goroutine that consumes the events every n seconds. (default is 30 seconds)
4. The events are then processed and inserted into a database. (default is a local SQLite database)
   - On each server start, the database is essentially wiped clean. See the [Configuration](#configuration) section for more information on how to use a different database.
5. The events can be queried through gRPC. There is a `.proto` file that defines the service and the message that can be sent over gRPC.
  - Check the [Querying](#querying) section for more information on how to query events.
6. The client is built with Wails which binds Go code to a frontend. `Sveltekit` is used for the frontend. The client can be used to log events and query events.
   - Check the [Running](#running) section for more information on how to run the client.
   - Allows you to use the service with a GUI; either through the desktop app or the web app.

Should you want to modify the server to your needs, you can edit the `.proto` files and regenerate the Go code. Check the [Editing the proto file](#editing-the-proto-file) section for more information.


### Logging

Logs can be sent over gRPC. The following is an example of a log message:

```json
{
  "id": 1,
  "message": "This is a log message",
  "level": "INFO",
  "severity": 0,
  "timestamp": 1612345678,
  "origin": "localhost",
  "source": "source",
  "type": "log",
  "group": "default",
  "tags": "tag1,tag2",
  "data": {
      "key1": "value1",
      "key2": "value2"
  }
}
```


- `id`: The unique identifier of the log message.
- `message`: The message of the log.
- `level`: The level of the log message; one of `DEBUG`, `INFO`, `WARN`, `ERROR`.
- `severity`: The severity of the message. Up to the user to interpret this value.
- `timestamp`: The timestamp of the log message.
- `origin`: The origin of the log message.
- `source`: The source of the log message. Source is different from origin in that source is the actual source of the log message, while origin is the origin of the log message.
- `type`: The type of the log message. You can use this to separate different types of log messages (e.g. `log`, `metric`, `event`).
- `group`: The group of the log messages. You can use this to separate different groups of log messages (e.g. `default`, `auth`, `db`).
- `tags`: The tags of the log message. You can use this to tag log messages (e.g. `tag1`, `tag2`).
- `data`: The data of the log message. You can use this to attach additional data to the log message.
- Optional fields: `origin`, `source`, `type`, `group`, `tags`, `data`.



### Querying

To do



## Testing

To do


## Configuration

### Using a different database

### Editing the proto file

If you want to edit the proto file, you will need to regenerate the Go code.
You can do this by running the following command:

```bash
protoc --go_out=. --go-grpc_out=. ./internal/logger/log.proto
```

### Changing the Kafka topic

To do


## Contributing

To do
