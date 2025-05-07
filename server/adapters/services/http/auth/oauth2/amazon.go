package oauth2

import (
	"encoding/json"

	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/ports/config"
	"github.com/jcfug8/daylear/server/ports/domain"

	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/amazon"
)

const amazonProvider = "amazon"

type AmazonService struct {
	*service
}

func NewAmazonService(l zerolog.Logger, configClient config.Client, domain domain.Domain) (*AmazonService, error) {
	config := configClient.GetConfig()
	amazonConfig := config["amazon"].(map[string]interface{})
	loginPath, ok := amazonConfig["loginpath"].(string)
	if !ok || loginPath == "" {
		loginPath = "/auth/amazon"
	}
	callbackPath, ok := amazonConfig["callbackpath"].(string)
	if !ok || loginPath == "" {
		callbackPath = "/auth/amazon/callback"
	}
	clientId := amazonConfig["clientid"].(string)
	clientSecret := amazonConfig["clientsecret"].(string)

	service := newService(NewServiceParams{
		L:               l,
		Name:            amazonProvider,
		Domain:          domain,
		LoginPath:       loginPath,
		CallbackPath:    callbackPath,
		GetUserInfoURL:  "https://api.amazon.com/user/profile",
		UserInfoFactory: fromAmazonUserInfo,
		OAuth2Config: &oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			Scopes: []string{
				"profile",
			},
			Endpoint: amazon.Endpoint,
		},
		SystemConfig: config,
	})
	return &AmazonService{service}, nil
}

type amazonUserInfo struct {
	ID    string `json:"user_id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// fromAmazonUserInfo -
func fromAmazonUserInfo(userInfoBytes []byte) (model.User, error) {
	var userInfo amazonUserInfo

	// unmarshal the response into the user struct
	err := json.Unmarshal(userInfoBytes, &userInfo)
	if err != nil {
		return model.User{}, err
	}

	return model.User{
		AmazonId: userInfo.ID,
		Email:    userInfo.Email,
	}, nil
}
