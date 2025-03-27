package server

import "net"

// NewMappedListener -
func NewMappedListener(name string, listener net.Listener) *MappedListener {
	return &MappedListener{
		Listener:   listener,
		ServerName: name,
	}
}

// MappedListener -
type MappedListener struct {
	net.Listener
	ServerName string
}
