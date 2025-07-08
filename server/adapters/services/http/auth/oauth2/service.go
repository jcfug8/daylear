package oauth2

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/ports/domain"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
)

type service struct {
	l    zerolog.Logger
	name string

	domain domain.Domain

	loginPath    string
	callbackPath string

	getUserInfoURL  string
	userInfoFactory func(userInfoBytes []byte) (model.User, error)

	config       *oauth2.Config
	uiDomainURL  *url.URL
	apiDomainURL *url.URL
}

type NewServiceParams struct {
	L               zerolog.Logger
	Name            string
	Domain          domain.Domain
	LoginPath       string
	CallbackPath    string
	GetUserInfoURL  string
	UserInfoFactory func(userInfoBytes []byte) (model.User, error)
	OAuth2Config    *oauth2.Config
	SystemConfig    map[string]interface{}
}

func newService(params NewServiceParams) *service {
	params.L.Printf("NewOAuthServer: loginPath=%s, callbackPath=%s", params.LoginPath, params.CallbackPath)

	uiDomainConfig := params.SystemConfig["uidomain"].(map[string]interface{})
	uiScheme := uiDomainConfig["scheme"].(string)
	uiHost := uiDomainConfig["host"].(string)
	uiPort, ok := uiDomainConfig["port"].(string)
	if ok {
		uiHost = fmt.Sprintf("%s:%s", uiHost, uiPort)
	}

	uiU := &url.URL{
		Scheme: uiScheme,
		Host:   uiHost,
	}

	apiDomainConfig := params.SystemConfig["apidomain"].(map[string]interface{})
	apiScheme := apiDomainConfig["scheme"].(string)
	apiHost := apiDomainConfig["host"].(string)
	apiPort, ok := apiDomainConfig["port"].(string)
	if ok {
		apiHost = fmt.Sprintf("%s:%s", apiHost, apiPort)
	}
	apiPath, _ := apiDomainConfig["path"].(string)
	apiPath = path.Join(apiPath, params.CallbackPath)

	apiU := &url.URL{
		Scheme: apiScheme,
		Host:   apiHost,
		Path:   apiPath,
	}

	params.OAuth2Config.RedirectURL = apiU.String()

	return &service{
		l:               params.L,
		name:            params.Name,
		domain:          params.Domain,
		loginPath:       params.LoginPath,
		callbackPath:    params.CallbackPath,
		getUserInfoURL:  params.GetUserInfoURL,
		userInfoFactory: params.UserInfoFactory,
		config:          params.OAuth2Config,
		uiDomainURL:     uiU,
		apiDomainURL:    apiU,
	}
}

func (s *service) Name() string {
	return s.name
}

func (s *service) Register(mux *http.ServeMux) error {
	mux.HandleFunc(s.loginPath, s.GetConsentURL)
	mux.HandleFunc(s.callbackPath, s.Callback)

	return nil
}

func (s *service) Close() error {
	return nil
}

func (s *service) GetConsentURL(w http.ResponseWriter, r *http.Request) {
	// Create oauthState cookie
	var expiration = time.Now().Add(20 * time.Minute)

	b := make([]byte, 16)
	rand.Read(b)
	oauthState := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: oauthState, Expires: expiration}
	http.SetCookie(w, &cookie)

	redirectURI := r.URL.Query().Get("redirect_uri")
	redirectCookie := http.Cookie{Name: "oauthclientredirecturi", Value: redirectURI, Expires: expiration}
	http.SetCookie(w, &redirectCookie)

	u := s.config.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func (s *service) Callback(w http.ResponseWriter, r *http.Request) {
	oauthState, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthState.Value {
		s.l.Printf("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	accessToken, err := s.config.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		s.l.Printf("code exchange wrong: %s", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	user, err := s.getUserByAccessToken(accessToken.AccessToken)
	if err != nil {
		s.l.Print(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	storedUser, err := s.domain.IdentifyUser(r.Context(), user)
	if err != nil && errors.Is(err, domain.ErrNotFound{}) {
		s.l.Print(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if storedUser.Id.UserId == 0 {
		storedUser, err = s.domain.CreateUser(r.Context(), user)
		if err != nil {
			s.l.Print(err.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
	}

	// login the user to get a login token
	tokenKey, err := s.domain.CreateToken(r.Context(), storedUser)
	if err != nil {
		s.l.Print(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Create a copy of the uiDomainURL
	redirectURL := s.uiDomainURL

	oauthclientredirecturi, _ := r.Cookie("oauthclientredirecturi")
	if oauthclientredirecturi != nil && oauthclientredirecturi.Value != "" {
		redirectURL, err = url.Parse(oauthclientredirecturi.Value)
		if err != nil {
			s.l.Print(err.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
	}

	http.SetCookie(w, &http.Cookie{Name: "oauthclientredirecturi", Value: "", Expires: time.Now().Add(-1 * time.Hour)})
	http.SetCookie(w, &http.Cookie{Name: "oauthstate", Value: "", Expires: time.Now().Add(-1 * time.Hour)})

	// Add a query parameter to the copied URL
	query := redirectURL.Query()
	query.Set("token_key", tokenKey)
	redirectURL.RawQuery = query.Encode()

	http.Redirect(w, r, redirectURL.String(), http.StatusTemporaryRedirect)
}

func (s service) getUserByAccessToken(accessToken string) (model.User, error) {
	// get value of user info url and set the access_token query param
	getUserInfoURL, err := url.Parse(s.getUserInfoURL)
	if err != nil {
		return model.User{}, err
	}
	q := getUserInfoURL.Query()
	q.Set("access_token", accessToken)
	getUserInfoURL.RawQuery = q.Encode()

	// make the request getting the user info
	resp, err := http.Get(getUserInfoURL.String())
	if err != nil {
		return model.User{}, err
	} else if resp.StatusCode != 200 {
		return model.User{}, fmt.Errorf("error getting user info from oauth2 api: %s", resp.Status)
	}
	defer resp.Body.Close()

	// read the response body
	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.User{}, err
	}

	// return the user struct generated by the user info factory
	user, err := s.userInfoFactory(contents)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
