# Golang Reddit Clone

A simple Reddit clone built for Go web development training, inspired by [course.gowebexamples.com](https://course.gowebexamples.com).

## Features

- User authentication
- Post creation and listing
- Commenting system

## Development Tools

- **Podman Compose** or **Docker Compose** for container orchestration:
  ```sh
  podman compose up
  # or
  docker compose up
  ```
- **Makefile** for database migrations:
  ```sh
  make migrate
  ```

## Live Reloading

- **Air** for live reloading during development:
  ```sh
  air
  ```
  Install Air with:
  ```sh
  go install github.com/cosmtrek/air@latest
  ```

## Source Material

This project is based on the curriculum from [course.gowebexamples.com](https://course.gowebexamples.com).
