package server

import (
	"net/http"
)

// Service constants
const (
	DefaultServerName = "default"
)

// Service - the interface for the http service.
type Service interface {
	Register(*http.ServeMux) error
	Close() error
	Name() string
}

// NewMappedServices -
func NewMappedServices(services ...Service) []*MappedService {
	result := make([]*MappedService, len(services))
	for i, service := range services {
		result[i] = NewMappedService(DefaultServerName, service)
	}
	return result
}

// NewMappedService -
func NewMappedService(name string, service Service) *MappedService {
	return &MappedService{
		Service:    service,
		ServerName: name,
	}
}

// MappedService -
type MappedService struct {
	Service
	ServerName string
}
