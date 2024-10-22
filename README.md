# Chitty-Chat: Distributed Chat System with Lamport Timestamps

## Project Overview
Chitty-Chat is a distributed chat service built using gRPC and Go. The system supports multiple participants, maintains logical time using Lamport timestamps, and logs all chat messages and events.

## Features
- Distributed chat system with real-time message broadcasting.
- Logical time tracking using Lamport timestamps.
- Server-client architecture with gRPC communication.
- Logs all join, leave, and chat events.

## How to Run

### Prerequisites
- Go 1.20 or higher
- Protobuf compiler (`protoc`)

### Running the Server
1. Navigate to the `server/` directory.
2. Build the server:
    ```bash
    go build -o server
    ```
3. Run the server:
    ```bash
    ./server.go
    ```

### Running the Client
1. Navigate to the `client/` directory.
2. Build the client:
    ```bash
    go build -o client
    ```
3. Run the client with a unique participant name:
    ```bash
    ./client.go <participant_name>
    ```

    Example:
    ```bash
    ./client Alice
    ```

## Contributing
Feel free to fork this repository and open pull requests for any improvements.
