package socket

import "net"

func Connect() (net.Conn, error) {
	return net.Dial(protocol, getSockAddr())
}
