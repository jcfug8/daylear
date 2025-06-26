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
				requester, ok := jwtConfig["requester"].(string)
				if !ok || requester == "" {
					return "", fmt.Errorf("missing jwt requester")
				}
				return requester, nil
			},
			fx.ResultTags(`name:"jwt_requester"`),
		),
	),
)
