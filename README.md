## About

This is a for fun project of mine that mimics a cloud storage. Used primarily
for learning how the all the systems interacts with each other: front end, back end,
network, security, authentication, and more.

It is built with Go, Python, Bash, React TS, Docker, and MariaDB.

## Requirements

The host machine *must be a Unix device*. WSL is supported.

Software requirements:
- Go >=1.25.0
- Node.js >= 22.22.1
- npm >= 10.9.4
- Git
- Docker
- Docker Compose

## Getting Started

Before starting the server, the configuration YAML must be set up.
This is required for the server to run.

```yml
# sample config

```

The frontend service must also have its own `.env` file with
the following:

**frontend/.env**:

```env
VITE_SERVER_BASE_URL=<server_host>
```