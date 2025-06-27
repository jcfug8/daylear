package token

import (
	"fmt"

	"github.com/jcfug8/daylear/server/ports/config"
	"github.com/jcfug8/daylear/server/ports/token"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"jwt",
	fx.Provide(
		fx.Annotate(
			NewTokenClient,
			fx.As(new(token.Client)),
		),
		fx.Annotate(
			func(configClient config.Client) (string, error) {
				jwtConfig := configClient.GetConfig()["jwt"].(map[string]interface{})
				secret, ok := jwtConfig["secret"].(string)
				if !ok || secret == "" {
					return "", fmt.Errorf("missing jwt secret")
				}
				return secret, nil
			},
			fx.ResultTags(`name:"jwt_secret"`),
		),
		fx.Annotate(
			func(configClient config.Client) (string, error) {
				jwtConfig := configClient.GetConfig()["jwt"].(map[string]interface{})
				issuer, ok := jwtConfig["issuer"].(string)
				if !ok || issuer == "" {
					return "", fmt.Errorf("missing jwt issuer")
				}
				return issuer, nil
			},
			fx.ResultTags(`name:"jwt_issuer"`),
		),
	),
)
