package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/jcfug8/daylear/server/adapters/netutils/listener"
	"github.com/jcfug8/daylear/server/ports/config"

	"github.com/rs/zerolog"
)

var _ Server = (*DefaultServer)(nil)

// Server -
type Server interface {
	Start() error
	Stop() error
}

// DefaultServer -
type DefaultServer struct {
	l          zerolog.Logger
	httpServer *http.Server
	services   []Service
	mux        *http.ServeMux
	listener   net.Listener
}

// NewServer - creates a new http server.
func NewServer(log zerolog.Logger, lis net.Listener, services []Service, configClient config.Client) (*DefaultServer, error) {
	log.Info().Msgf("creating new http server with %d services", len(services))

	mux := http.NewServeMux()

	uiDomainConfig := configClient.GetConfig()["uidomain"].(map[string]interface{})
	uiScheme := uiDomainConfig["scheme"].(string)
	uiHost := uiDomainConfig["host"].(string)
	uiPort, ok := uiDomainConfig["port"].(string)
	if ok {
		uiHost = fmt.Sprintf("%s:%s", uiHost, uiPort)
	}

	uiU := url.URL{
		Scheme: uiScheme,
		Host:   uiHost,
	}

	apiDomainConfig := configClient.GetConfig()["apidomain"].(map[string]interface{})
	apiScheme := apiDomainConfig["scheme"].(string)
	apiHost := apiDomainConfig["host"].(string)
	apiPort, ok := apiDomainConfig["port"].(string)
	if ok {
		apiHost = fmt.Sprintf("%s:%s", apiHost, apiPort)
	}

	apiU := url.URL{
		Scheme: apiScheme,
		Host:   apiHost,
	}

	corsConfig := configClient.GetConfig()["cors"].(map[string]interface{})
	originStrings := corsConfig["extraorigins"].(string)
	origins := append(strings.Split(originStrings, ","), uiU.String(), apiU.String())

	log.Info().Msgf("cors origins: %v", origins)

	return &DefaultServer{
		httpServer: &http.Server{
			Addr:    lis.Addr().String(),
			Handler: NewMiddlewareMux(log, origins, mux),
		},
		services: services,
		mux:      mux,
		listener: lis,
	}, nil
}

// start - starts http server and registers the services.
func (s *DefaultServer) Start() (err error) {
	s.l.Info().Msgf("starting http server on %s", s.listener.Addr().String())

	lis := s.listener
	if lis == nil {
		lis, err = net.Listen("tcp", s.httpServer.Addr)
	}

	if err != nil {
		return fmt.Errorf("unable to listen on %s: %w", s.httpServer.Addr, err)
	}

	// register services
	for _, service := range s.services {
		s.l.Info().Msgf("registering service %s", service.Name())
		if err = service.Register(s.mux); err != nil {
			return fmt.Errorf("unable to register service %s: %w", service.Name(), err)
		}
	}

	go func() {
		if err := s.httpServer.Serve(lis); err != nil {
			s.l.Error().Err(err).Msg("failed to serve http server")
		}
	}()

	return nil
}

// stop - stops the http server.
func (s *DefaultServer) Stop() error {
	s.l.Info().Msg("stopping http server")

	for _, service := range s.services {
		s.l.Info().Msgf("closing service %s", service.Name())
		if err := service.Close(); err != nil {
			s.l.Error().Err(err).Msgf("unable to close service %s", service.Name())
		}
	}

	s.httpServer.Shutdown(context.Background())
	return nil
}

// GetAddress -
func GetAddress(configClient config.Client) *listener.Address {
	config := configClient.GetConfig()["http"].(map[string]interface{})

	host, ok := config["host"].(string)
	if !ok || host == "" {
		host = "0.0.0.0"
	}
	port, ok := config["port"].(string)
	if !ok || port == "" {
		port = "8080"
	}
	iPort, err := strconv.Atoi(port)
	if err != nil {
		panic(fmt.Errorf("unable to convert port to int: %w", err))
	}

	return &listener.Address{
		Host: host,
		Port: iPort,
	}
}
