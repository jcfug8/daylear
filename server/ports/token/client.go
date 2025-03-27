package token

import "github.com/jcfug8/daylear/server/core/model"

type Client interface {
	Encode(model.User) (string, error)
	Decode(string) (model.User, error)
}
