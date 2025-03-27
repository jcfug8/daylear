package auth

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/jcfug8/daylear/server/adapters/grpc/users/user/v1alpha1/namer"
	"github.com/jcfug8/daylear/server/ports/config"
	"github.com/jcfug8/daylear/server/ports/domain"
)

type Service struct {
	domain      domain.Domain
	userNamer   namer.UserNamer
	uiDomainURL url.URL
}

func NewService(domain domain.Domain, configClient config.Client) *Service {
	userNamer, err := namer.NewUserNamer()
	if err != nil {
		panic(err)
	}

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
	return &Service{
		domain:      domain,
		userNamer:   userNamer,
		uiDomainURL: uiU,
	}
}

func (s *Service) Register(m *http.ServeMux) error {
	m.HandleFunc("/auth/check", s.AuthCheck)
	m.HandleFunc("/auth/token/", s.GetToken)
	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) Name() string {
	return "http-auth"
}
