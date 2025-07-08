package token

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/users/user/v1alpha1"
	"github.com/jcfug8/daylear/server/ports/config"
	"github.com/jcfug8/daylear/server/ports/domain"
)

type Service struct {
	domain      domain.Domain
	userNamer   namer.ReflectNamer
	uiDomainURL url.URL
}

func NewService(domain domain.Domain, configClient config.Client) *Service {
	userNamer, err := namer.NewReflectNamer[*pb.User]()
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
	m.Handle("/auth/token/", http.HandlerFunc(s.GetToken))
	m.Handle("/auth/check/token", headers.NewAuthTokenMiddleware(s.domain)(http.HandlerFunc(s.CheckToken)))
	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) Name() string {
	return "http-auth"
}
