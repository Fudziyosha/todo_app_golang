# Go Todo App
![Static Badge](https://img.shields.io/badge/golang-00ADD8?&style=plastic&logo=go&logoColor=white)
[![Go Reference](https://pkg.go.dev/badge/golang.org/x/pkgsite.svg)](https://pkg.go.dev/golang.org/x/pkgsite)
![Go Version](https://img.shields.io/badge/go-1.25-blue.svg)
![Psql](https://img.shields.io/badge/PostgreSQL-316192?logo=postgresql&logoColor=white)


![Imgur](https://i.imgur.com/gvXRkfP.png)

## Description

A simple and intuitive task management web app that allows users to create lists, add tasks, organize them with tabs, and keep everything up to date.

## Features

- User authentication (login / registration)
- Create and manage task lists
- Add, update, and delete tasks
- Tab-based organization
- Bulk deletion of lists
- Responsive UI

## Tech stack

- Go (Golang) 1.25
- Fiber — fast HTTP web framework
- PostgreSQL
- Redis 
- Server-side rendering (HTML templates)
- HTML, CSS (Bootstrap)

## Getting Started

1. Clone the repository:

```bash
git clone https://github.com/Fudziyosha/todo_app_golang
```
2. Copy `.env.example` to `.env` and configure the database and Redis credentials:

```bash
cp .env.example .env
```

3. Update the config.yaml file: set the logger level to info and enable static file serving (server.statics: true)
4. Run the application locally:
```bash
cd cmd/app
go run main.go
```

## Configuration

The application uses:

- `.env` — for environment variables (database, Redis)
- `config.yaml` — for server and logger configuration

## Makefile Commands

Common commands:

- `make run` — start the application
- `make build` — build the binary
- `make dm-up` — run services with Docker
