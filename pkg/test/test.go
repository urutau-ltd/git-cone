package test

import (
	"fmt"
	"net"
	"sync"
)

var (
	used = map[int]struct{}{}
	lock sync.Mutex
)

// RandomPort returns a random port number.
// This is mainly used for testing.
func RandomPort() int {
	addr, err := net.Listen("tcp", "127.0.0.1:0") //nolint:gosec,noctx
	if err != nil {
		panic(fmt.Sprintf("reserve random tcp port: %v", err))
	}
	_ = addr.Close()
	port := addr.Addr().(*net.TCPAddr).Port
	lock.Lock()

	if _, ok := used[port]; ok {
		lock.Unlock()
		return RandomPort()
	}

	used[port] = struct{}{}
	lock.Unlock()
	return port
}
