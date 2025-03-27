package oauth2

import (
	"encoding/json"

	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/ports/config"
	"github.com/jcfug8/daylear/server/ports/domain"

	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

const facebookProvider = "facebook"

type FacebookService struct {
	*service
}

func NewFacebookService(l zerolog.Logger, configClient config.Client, domain domain.Domain) (*FacebookService, error) {
	config := configClient.GetConfig()
	facebookConfig := config["facebook"].(map[string]interface{})
	loginPath, ok := facebookConfig["loginpath"].(string)
	if !ok || loginPath == "" {
		loginPath = "/auth/facebook"
	}
	callbackPath, ok := facebookConfig["callbackpath"].(string)
	if !ok || loginPath == "" {
		callbackPath = "/auth/facebook/callback"
	}
	clientId := facebookConfig["clientid"].(string)
	clientSecret := facebookConfig["clientsecret"].(string)

	service := newService(NewServiceParams{
		L:               l,
		Name:            facebookProvider,
		Domain:          domain,
		LoginPath:       loginPath,
		CallbackPath:    callbackPath,
		GetUserInfoURL:  "https://graph.facebook.com/me?fields=email,name,first_name,last_name,middle_name,picture",
		UserInfoFactory: fromFacebookUserInfo,
		OAuth2Config: &oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			Scopes: []string{
				"email",
				"public_profile",
			},
			Endpoint: facebook.Endpoint,
		},
		SystemConfig: config,
	})
	return &FacebookService{service}, nil
}

type facebookUserInfo struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	// Picture    string `json:"picture"`
}

// fromFacebookUserInfo -
func fromFacebookUserInfo(userInfoBytes []byte) (model.User, error) {
	var userInfo facebookUserInfo

	// unmarshal the response into the user struct
	err := json.Unmarshal(userInfoBytes, &userInfo)
	if err != nil {
		return model.User{}, err
	}

	return model.User{
		FacebookId: userInfo.ID,
		Email:      userInfo.Email,
	}, nil
}
