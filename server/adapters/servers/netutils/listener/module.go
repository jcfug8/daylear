package listener

import (
	"fmt"

	"go.uber.org/fx"
)

// ProvideListener -
func ProvideListener(addrName, lisName string) fx.Option {
	return fx.Provide(
		fx.Annotate(
			BuildListener,
			fx.ParamTags(fmt.Sprintf(`name:"%s"`, addrName)),
			fx.ResultTags(fmt.Sprintf(`name:"%s" optional:"true"`, lisName)),
		),
	)
}
