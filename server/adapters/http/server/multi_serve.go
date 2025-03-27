package server

import (
	"fmt"
	"sync"

	"github.com/jcfug8/daylear/server/ports/config"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

// MultiServerParams -
type MultiServerParams struct {
	fx.In

	Logger       zerolog.Logger
	Services     []*MappedService  `group:"httpMappedServices"`
	Listeners    []*MappedListener `group:"httpMappedListeners"`
	ConfigClient config.Client
}

// NewMultiServer -
func NewMultiServer(p MultiServerParams) (*MultiServer, error) {
	l := p.Logger.With().Str("module", "http_server").Logger()

	servers := make([]Server, 0, len(p.Listeners))
	names := make([]string, 0, len(p.Listeners))

	for _, listener := range p.Listeners {
		services := make([]Service, 0, len(p.Services))
		serviceNames := make([]string, 0, len(p.Services))
		for _, service := range p.Services {
			if service.ServerName == listener.ServerName {
				services = append(services, service.Service)
				serviceNames = append(serviceNames, service.Name())
			}
		}

		if len(services) == 0 {
			l.Warn().Msgf("no services found for listener %s, skipping", listener.ServerName)
			continue
		}

		l.Info().Strs("services", serviceNames).Msgf("creating server for listener %s with %d services on %s", listener.ServerName, len(services), listener.Addr())
		server, err := NewServer(l, listener.Listener, services, p.ConfigClient)
		if err != nil {
			return nil, fmt.Errorf("unable to create server: %w", err)
		}

		servers = append(servers, server)
		names = append(names, listener.ServerName)
	}

	l.Info().Strs("servers", names).Msgf("created %d http servers", len(servers))
	return &MultiServer{
		log:     l,
		servers: servers,
	}, nil
}

// MultiServer -
type MultiServer struct {
	log     zerolog.Logger
	servers []Server
}

// Start -
func (s *MultiServer) Start() (err error) {
	wg := sync.WaitGroup{}
	wg.Add(len(s.servers))

	errored := false
	for _, server := range s.servers {
		go func(server Server) {
			defer wg.Done()

			if err = server.Start(); err != nil {
				s.log.Error().Err(err).Msgf("unable to start server")
				errored = true
			}
		}(server)
	}

	wg.Wait()

	if errored {
		err = fmt.Errorf("encountered errors when starting servers")
	}

	return err
}

// Stop -
func (s *MultiServer) Stop() (err error) {
	wg := sync.WaitGroup{}
	wg.Add(len(s.servers))

	errored := false
	for _, server := range s.servers {
		go func(server Server) {
			defer wg.Done()

			if err = server.Stop(); err != nil {
				s.log.Error().Err(err).Msgf("unable to stop server")
				errored = true
			}
		}(server)
	}

	wg.Wait()

	if errored {
		err = fmt.Errorf("encountered errors when stopping servers")
	}

	return err
}
