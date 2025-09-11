default:
	@just --list

update:
	go get -u ./...
	go mod tidy
