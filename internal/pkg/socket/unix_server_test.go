package socket_test

import (
	"net"
	"testing"
	"time"

	"github.com/micronull/i3rotonda/internal/pkg/socket"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	ch := make(chan bool)
	defer close(ch)

	addr := socket.Run(func(conn *net.Conn) {
		ch <- true
	})

	go func() {
		_, err := net.Dial("unix", addr.String())
		assert.NoError(t, err)
	}()

	assert.Eventually(t, func() bool { return <-ch }, time.Second*5, time.Millisecond, "no connected")
}
