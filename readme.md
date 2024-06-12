# REDIS-X

REDIS-X is a lightweight, custom implementation of a Redis-like key-value store written in Go. This project includes basic functionalities such as setting and getting string values, as well as working with hash maps. Additionally, it supports the Redis Serialization Protocol (RESP) for communication.

## Features

- RESP (REdis Serialization Protocol) parser and writer.
- Basic Redis commands: `PING`, `SET`, `GET`, `HSET`, `HGET`, `HGETALL`.
- Append-only file (AOF) for persistence.
- Concurrent safe operations using mutexes.

## Getting Started

### Prerequisites

- Go (version 1.15+)
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/utsavvirani/redis-x.git
   cd redis-x
   ```

2. Build the project:
   ```bash
   go build -o redis-x
   ```

3. Run the server:
   ```bash
   ./redis-x
   ```

The server will start listening on port `6379`.

## Usage

### Supported Commands

- **PING**: Check if the server is running.
  ```shell
  $ redis-cli PING
  PONG
  ```

- **SET**: Set a string key-value pair.
  ```shell
  $ redis-cli SET key value
  OK
  ```

- **GET**: Get the value of a string key.
  ```shell
  $ redis-cli GET key
  value
  ```

- **HSET**: Set a hash field to a value.
  ```shell
  $ redis-cli HSET myhash field value
  OK
  ```

- **HGET**: Get the value of a hash field.
  ```shell
  $ redis-cli HGET myhash field
  value
  ```

- **HGETALL**: Get all fields and values of a hash.
  ```shell
  $ redis-cli HGETALL myhash
  field
  value
  ```

### Persistence

The server supports persistence through an Append-Only File (AOF). All `SET` and `HSET` commands are written to the `database.aof` file for durability. The file is synced to disk every second.

## Code Structure

- `main.go`: The entry point of the server.
- `resp.go`: Contains the RESP protocol implementation.
- `handlers.go`: Implements the command handlers.
- `aof.go`: Manages the Append-Only File (AOF) for persistence.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

- Utsav Virani

---

Feel free to reach out if you have any questions or suggestions. Happy coding!