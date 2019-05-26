all:
	go build -o calc

tutorial1:
	go build gotutorial/tutorial1
test:
	go test -v ./...
benchmark:
go test -bench . ./...