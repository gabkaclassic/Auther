build_app:
	go build -o auther cmd/main.go

test:
	go test ./...

clean:
	rm auther

