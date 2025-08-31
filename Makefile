udpsen:
	@go run ./cmd/udpsender

tcplis:
	@go run ./cmd/tcplistener

test:
	@go test -v ./...
