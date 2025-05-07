package listener

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Listen used for testing
var Listen = net.Listen

// Address -
type Address struct {
	Host      string
	Port      int
	IsDynamic bool
}

// BuildListener -
func BuildListener(addr *Address) (net.Listener, error) {
	if addr == nil {
		return nil, nil
	}

	return NewListener(addr.Host, addr.Port, addr.IsDynamic)
}

// NewListener creates a net.Listener on the given host and port. When isDynamic is true
// the first available port in the range (port - 65535] is used. When isDynamic is false
// it errors if the port is not available.
func NewListener(host string, port int, isDynamic bool) (lis net.Listener, err error) {
	maxPort := 65535
	if port > maxPort {
		return nil, fmt.Errorf("provided port exceeds max range")
	}

	for {
		addr := net.JoinHostPort(host, strconv.Itoa(port))
		lis, err = Listen("tcp", addr)

		if err == nil {
			break
		}

		if isDynamic && strings.Contains(err.Error(), "bind: address already in use") {
			if port > maxPort {
				err = fmt.Errorf("could not find open port")
			} else {
				port++
				continue
			}
		}

		return nil, fmt.Errorf("could not create tcp listener: %v", err)
	}

	return lis, nil
}
