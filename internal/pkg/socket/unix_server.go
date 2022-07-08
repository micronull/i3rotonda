package socket

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/user"
)

const (
	protocol = "unix"
	sockAddr = "/run/user/%s/i3rotonda.sock"
)

func Run(handler func(net.Conn)) (addr net.Addr) {
	if err := cleanup(); err != nil {
		log.Fatal(err)
	}

	listener := listen()

	go func() {
		defer func() {
			if err := listener.Close(); err != nil {
				log.Printf("WARN: %s", err.Error())
			}
		}()

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatal(err)
			}

			go handler(conn)
		}
	}()

	return listener.Addr()
}

func listen() net.Listener {
	listener, err := net.Listen(protocol, getSockAddr())
	if err != nil {
		log.Fatal(err)
	}

	return listener
}

func getSockAddr() string {
	info, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf(sockAddr, info.Uid)
}

func cleanup() (err error) {
	addr := getSockAddr()

	if _, err = os.Stat(addr); err == nil {
		if err = os.RemoveAll(addr); err != nil {
			return err
		}
	}

	return nil
}
