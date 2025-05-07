package http

import (
	"fmt"
	"net"

	listener "github.com/jcfug8/daylear/server/adapters/servers/netutils/listener"

	"go.uber.org/fx"
)

// Fx tags
const (
	AddressTag         = `name:"httpServerAddress"`
	ListenerTag        = `name:"httpServerListener"`
	MappedListenersTag = `group:"httpMappedListeners"`
	MappedServicesTag  = `group:"httpMappedServices"`
	ServicesTag        = `group:"httpServices"`
)

// Module - the fx module for the http server.
var Module = fx.Module(
	"httpServer",
	fx.Provide(
		fx.Annotate(
			NewMappedServices,
			fx.ParamTags(ServicesTag),
			fx.ResultTags(`group:"httpMappedServices,flatten"`),
		),

		fx.Annotate(
			NewMultiServer,
			fx.OnStart(func(server *MultiServer) error {
				return server.Start()
			}),
			fx.OnStop(func(server *MultiServer) error {
				return server.Stop()
			}),
		),

		fx.Annotate(
			GetAddress,
			fx.ResultTags(fmt.Sprintf(`%s optional:"true"`, AddressTag)),
		),

		fx.Annotate(
			func(lis net.Listener) *MappedListener {
				return NewMappedListener(DefaultServerName, lis)
			},
			fx.ParamTags(ListenerTag),
			fx.ResultTags(MappedListenersTag),
		),
	),

	listener.ProvideListener("httpServerAddress", "httpServerListener"),

	fx.Invoke(func(*MultiServer) {}),
)

// ProvideAsService provides a type T as a Service for the named servers.
func ProvideAsService[T Service](serverNames ...string) fx.Option {
	if len(serverNames) == 0 {
		return fx.Provide(
			fx.Annotate(
				func(service T) *MappedService {
					return NewMappedService(DefaultServerName, service)
				},
				fx.ResultTags(MappedServicesTag),
			),
		)
	}

	opts := make([]any, 0, len(serverNames))

	for _, serverName := range serverNames {
		opts = append(opts, fx.Annotate(
			func(service T) *MappedService {
				return NewMappedService(serverName, service)
			},
			fx.ResultTags(MappedServicesTag),
		))
	}

	return fx.Provide(opts...)
}

// ProvideListener provides a listener for the named server.
func ProvideListener(lisName, serverName string) fx.Option {
	return fx.Provide(
		fx.Annotate(
			func(lis net.Listener) *MappedListener {
				return NewMappedListener(serverName, lis)
			},
			fx.ParamTags(lisName),
			fx.ResultTags(MappedListenersTag),
		),
	)
}

// ProvideDevcloneListener -
func ProvideDevcloneListener() fx.Option {
	return fx.Replace(
		fx.Annotate(
			&listener.Address{
				Port:      8080,
				IsDynamic: true,
			},
			fx.ResultTags(AddressTag),
		),
	)
}
