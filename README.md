# Task Manager API
API responsible for creating and managing tasks.

## Introduction

The API itself is for study purposes, and aims to be a simple API for crud and task management, using Docker, Redis and working in concurrency with goroutines.

## 🛠 Main libraries used

- **pgxpool** - A package from the pgx library, pgxpool provides a connection pool for PostgreSQL databases. It allows efficient management of database connections and helps in handling multiple concurrent requests to the database with reduced latency and resource consumption.

- **github.com/google/uuid** - This library provides utilities for generating and working with UUIDs (Universally Unique Identifiers). It supports the creation of various types of UUIDs, including UUIDv1, UUIDv4, and UUIDv5, and is used to ensure unique identifiers in distributed systems and databases.

- **github.com/redis/go-redis/v9** - This is a Go client library for interacting with Redis, a popular in-memory data structure store. The go-redis library provides support for various Redis commands and features, including caching, pub/sub messaging, and data persistence, allowing Go applications to efficiently communicate with Redis servers.

For a complete list of dependencies see the file [go.mod](https://github.com/felipeversiane/task-api/blob/main/go.mod).


## 🚀 Development build

To be able to run the project as a developer, follow the steps:

 - Clone the repository to your machine.
 - Remember to checkout the branch you are developing.

From there you can run the following command to run the project and test the application

```bash
  docker-compose up --build
```

That's it, the API is running, be happy!

## Suport

For support, please email me [felipeversiane09@gmail.com](mailto:felipeversiane09@gmail.com)

