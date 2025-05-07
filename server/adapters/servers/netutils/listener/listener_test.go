package listener_test

import (
	"context"
	"net"
	"testing"

	"github.com/jcfug8/daylear/server/adapters/servers/netutils/listener"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestListener(t *testing.T) {
	NewListenerSuite(t).Test_BuildListener()
}

func (suite *ListenerSuite) Test_BuildListener() {
	type T = ListenerSuite

	suite.It("should do nothing when address is nil", func(t *T) {
		called := false

		app := t.fxApp(
			listener.ProvideListener("addr", "lis"),

			fx.Supply(
				fx.Annotate(
					(*listener.Address)(nil),
					fx.ResultTags(`name:"addr"`),
				),
			),

			fx.Invoke(func(net.Listener) {
				called = true
			}),
		)

		err := app.Start(ctx)
		require.Regexp(t, "^missing dependencies.*", err)

		require.False(t, called)
	})
}

// ----------------------------------------------------------------------------

var (
	ctx = context.Background()
)

func NewListenerSuite(t *testing.T) *ListenerSuite {
	return &ListenerSuite{T: t}
}

type ListenerSuite struct {
	*testing.T
}

func (suite *ListenerSuite) It(name string, yield func(*ListenerSuite)) {
	suite.Helper()
	suite.Run(name, yield)
}

func (suite *ListenerSuite) Run(name string, yield func(*ListenerSuite)) {
	suite.Helper()
	suite.T.Run(name, func(t *testing.T) {
		suite := NewListenerSuite(t)
		suite.Helper()

		yield(suite)
	})
}

func (suite *ListenerSuite) fxApp(opts ...fx.Option) *fx.App {
	return fx.New(opts...)
}
