build:
	go test ./...
	go build ./cmd/i3rotonda/

install: build
	mv -f ./i3rotonda ~/.local/bin/

uninstall:
	rm ~/.local/bin/i3rotonda