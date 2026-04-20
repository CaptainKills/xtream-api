MAKEFLAGS += --silent

run:
	go run .

test:
	go test ./...

build:
	go build .

upload:
	go build .
	cp xtream-api ~/sshfs/docker-proxmox/xtream-api/
