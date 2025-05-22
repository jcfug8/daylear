package oauth2

import (
	"encoding/json"

	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/ports/config"
	"github.com/jcfug8/daylear/server/ports/domain"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const googleProvider = "google"

type GoogleService struct {
	*service
}

func NewGoogleService(l zerolog.Logger, configClient config.Client, domain domain.Domain) (*GoogleService, error) {
	config := configClient.GetConfig()
	googleConfig := config["google"].(map[string]interface{})
	loginPath, ok := googleConfig["loginpath"].(string)
	if !ok || loginPath == "" {
		loginPath = "/auth/google"
	}
	callbackPath, ok := googleConfig["callbackpath"].(string)
	if !ok || callbackPath == "" {
		callbackPath = "/auth/google/callback"
	}
	clientId := googleConfig["clientid"].(string)
	clientSecret := googleConfig["clientsecret"].(string)

	service := newService(NewServiceParams{
		L:               l,
		Name:            googleProvider,
		Domain:          domain,
		LoginPath:       loginPath,
		CallbackPath:    callbackPath,
		GetUserInfoURL:  "https://www.googleapis.com/oauth2/v2/userinfo",
		UserInfoFactory: fromGoogleUserInfo,
		OAuth2Config: &oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
				"openid",
			},
			Endpoint: google.Endpoint,
		},
		SystemConfig: config,
	})
	return &GoogleService{service}, nil
}

type googleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// fromGoogleUserInfo -
func fromGoogleUserInfo(userInfoBytes []byte) (model.User, error) {
	var userInfo googleUserInfo

	// unmarshal the response into the user struct
	err := json.Unmarshal(userInfoBytes, &userInfo)
	if err != nil {
		return model.User{}, err
	}

	return model.User{
		GoogleId: userInfo.ID,
		Email:    userInfo.Email,
	}, nil
}
