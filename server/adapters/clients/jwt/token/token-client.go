package token

import (
	"fmt"

	"github.com/jcfug8/daylear/server/core/model"
	"go.uber.org/fx"

	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog"
)

type TokenClient struct {
	l         zerolog.Logger
	secret    string
	requester string
}

type NewTokenClientParams struct {
	fx.In

	L         zerolog.Logger
	Secret    string `name:"jwt_secret"`
	requester string `name:"jwt_requester"`
}

func NewTokenClient(params NewTokenClientParams) TokenClient {
	return TokenClient{
		l:         params.L,
		secret:    params.Secret,
		requester: params.requester,
	}
}

func (t TokenClient) Encode(user model.User) (string, error) {
	var err error

	t.l.Printf("saving user to token: %+v", user)

	claims := toJWTClaims(user)
	claims["iss"] = t.requester

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))

	tokenString, err := token.SignedString([]byte(t.secret))
	if err != nil {
		t.l.Printf("unable to sign token: %s", err)
		return "", fmt.Errorf("unable to sign token")
	}

	return tokenString, nil
}

func (t TokenClient) Decode(tn string) (model.User, error) {
	token, err := jwt.Parse(tn, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			t.l.Printf("incorrect signing method")
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(t.secret), nil
	})

	if err != nil {
		t.l.Printf("error parsing token: ", err)
		return model.User{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
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
