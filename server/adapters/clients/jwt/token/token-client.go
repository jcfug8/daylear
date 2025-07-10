package token

import (
	"context"
	"fmt"

	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/ports/token"
	"go.uber.org/fx"

	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog"
)

var _ token.Client = (*TokenClient)(nil)

type TokenClient struct {
	l      zerolog.Logger
	secret string
	issuer string
}

type NewTokenClientParams struct {
	fx.In

	L      zerolog.Logger
	Secret string `name:"jwt_secret"`
	Issuer string `name:"jwt_issuer"`
}

func NewTokenClient(params NewTokenClientParams) TokenClient {
	return TokenClient{
		l:      params.L,
		secret: params.Secret,
		issuer: params.Issuer,
	}
}

func (t TokenClient) Encode(ctx context.Context, user model.User) (string, error) {
	log := logutil.EnrichLoggerWithContext(t.l, ctx)
	var err error

	log.Debug().Msgf("saving user to token: %+v", user)

	claims := toJWTClaims(user)
	claims["iss"] = t.issuer

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))

	tokenString, err := token.SignedString([]byte(t.secret))
	if err != nil {
		log.Error().Err(err).Msg("unable to sign token")
		return "", fmt.Errorf("unable to sign token")
	}

	return tokenString, nil
}

func (t TokenClient) Decode(ctx context.Context, tn string) (model.User, error) {
	log := logutil.EnrichLoggerWithContext(t.l, ctx)
	token, err := jwt.Parse(tn, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Error().Msg("unexpected signing method")
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(t.secret), nil
	})

	if err != nil {
		log.Error().Err(err).Msg("error parsing token")
		return model.User{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Error().Msg("invalid token claims")
		return model.User{}, err
	}

	user := fromJWTClaims(claims)

	return user, nil
}

// toJWTClaims - Converts a user to a jwt.MapClaims
func toJWTClaims(user model.User) jwt.MapClaims {
	return jwt.MapClaims{
		"id":         user.Id.UserId,
		"email":      user.Email,
		"amazonId":   user.AmazonId,
		"facebookId": user.FacebookId,
		"googleId":   user.GoogleId,
	}
}

// fromJWTClaims - Converts a jwt.MapClaims to a user
func fromJWTClaims(m jwt.MapClaims) model.User {
	floatId, _ := m["id"].(float64)
	id := int64(floatId)
	return model.User{
		Id:         model.UserId{UserId: id},
		Email:      m["email"].(string),
		AmazonId:   m["amazonId"].(string),
		FacebookId: m["facebookId"].(string),
		GoogleId:   m["googleId"].(string),
	}
}
