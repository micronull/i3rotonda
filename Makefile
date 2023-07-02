build:
	go build ./cmd/i3rotonda/

install: build
	mv -f ./i3rotonda ~/go/bin/

uninstall:
	rm ~/go/bin/i3rotonda