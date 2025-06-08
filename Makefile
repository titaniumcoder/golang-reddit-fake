.PHONY: postgres adminer migrate down

postgres:
	podman compose up

up:
	migrate -source file://migrations -database postgres://postgres:secret@localhost/postgres?sslmode=disable up

down:
	migrate -source file://migrations -database postgres://postgres:secret@localhost/postgres?sslmode=disable down
