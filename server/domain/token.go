package domain

import (
	"context"
	"fmt"

	model "github.com/jcfug8/daylear/server/core/model"
	uuid "github.com/satori/go.uuid"
)

func (d *Domain) CreateToken(ctx context.Context, user model.User) (string, error) {
	token, err := d.tokenClient.Encode(user)
	if err != nil {
		return "", err
	}

	key := uuid.NewV4().String()
	d.tokenStore.Store(key, token)

	return key, nil
}

func (d *Domain) RetrieveToken(ctx context.Context, key string) (string, error) {
	token, ok := d.tokenStore.Load(key)
	if !ok {
		return "", fmt.Errorf("token not found")
	}

	d.tokenStore.Delete(key)

	return token.(string), nil
}

func (d *Domain) ParseToken(ctx context.Context, token string) (model.User, error) {
	return d.tokenClient.Decode(token)
}
