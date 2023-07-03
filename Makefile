build:
	go test ./...
	export GOMAXPROCS=1 && go build ./cmd/i3rotonda/

install: build
	mv -f ./i3rotonda ~/.local/bin/

uninstall:
	rm ~/.local/bin/i3rotonda