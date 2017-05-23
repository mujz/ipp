package authentication

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

type FacebookAuthenticator struct {
	model              FacebookUserCreator
	auth               *Authenticator
	authURL            string
	successRedirectURL *url.URL
	failureRedirectURL string
	*oauth2.Config
}

type fbUser struct {
	ID   string
	Name string
}

type FacebookUserCreator interface {
	CreateFacebookUser(string) (int, error)
}

type Key int

const (
	oauthStateString = "randoms123sawr32vyj"
	TokenKey         = Key(1)
	ErrorKey         = Key(2)
)

var (
	facebookProfileAPI      = "https://graph.facebook.com/me"
	InvalidFacebookStateErr = errors.New("authentication: invalid oauth state.")
	InvalidFacebookTokenErr = errors.New("authentication: failed to parse Facebook token.")
	InvalidUserInfoErr      = errors.New("authentication: failed to parse Facebook user info.")
	UserInfoFetchFailedErr  = errors.New("authentication: failed to fetch user info from Facebook's Graph API")
)

func NewFacebookAuthenticator(appID, appSecret, redirectURL, successRedirectURL, failureRedirectURL string, scopes []string, userCreator FacebookUserCreator, auth *Authenticator) *FacebookAuthenticator {

	successURL, err := url.Parse(successRedirectURL)
	if err != nil {
		panic(err)
	}

	f := &FacebookAuthenticator{
		model:              userCreator,
		auth:               auth,
		successRedirectURL: successURL,
		failureRedirectURL: failureRedirectURL,
		Config: &oauth2.Config{
			ClientID:     appID,
			ClientSecret: appSecret,
			RedirectURL:  redirectURL,
			Scopes:       scopes,
			Endpoint:     facebook.Endpoint,
		},
	}

	f.setAuthURL()

	return f
}

func (f *FacebookAuthenticator) setAuthURL() {
	Url, _ := url.Parse(f.Endpoint.AuthURL)
	params := url.Values{}

	params.Add("client_id", f.ClientID)
	params.Add("state", oauthStateString)
	if scopes := f.Scopes; scopes != nil && len(scopes) > 0 {
		params.Add("scope", strings.Join(scopes, ","))
	}
	params.Add("redirect_uri", f.RedirectURL)
	params.Add("response_type", "code")

	Url.RawQuery = params.Encode()
	f.authURL = Url.String()
}

// Redirects the user to the Facebook login page
func (f *FacebookAuthenticator) LoginHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, f.authURL, http.StatusTemporaryRedirect)
}

func (f *FacebookAuthenticator) LoginCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Verify state is the same as the one we sent
	state := r.FormValue("state")
	if state != oauthStateString {
		log.Println(InvalidFacebookStateErr)
		http.Redirect(w, r, f.failureRedirectURL, http.StatusTemporaryRedirect)
		return
	}

	// Get the token from the code
	code := r.FormValue("code")
	fbToken, err := f.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, f.failureRedirectURL, http.StatusTemporaryRedirect)
		return
	}

	// Get the user info (Name and Facebook ID)
	resp, err := http.Get(facebookProfileAPI + "?access_token=" +
		url.QueryEscape(fbToken.AccessToken))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, f.failureRedirectURL, http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	// Decode the user info
	user := new(fbUser)
	if err := json.NewDecoder(resp.Body).Decode(user); err != nil {
		log.Println(err)
		http.Redirect(w, r, f.failureRedirectURL, http.StatusTemporaryRedirect)
		return
	}

	// Create a new Facebook user
	userID, err := f.model.CreateFacebookUser(user.ID)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, f.failureRedirectURL, http.StatusTemporaryRedirect)
		return
	}

	// Generate a token for the user
	token, err := f.auth.generateToken(userID)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, f.failureRedirectURL, http.StatusTemporaryRedirect)
		return
	}

	params := url.Values{}
	params.Add("token", token.Subject)
	f.successRedirectURL.RawQuery = params.Encode()
	// Redirect to success redirect url with the token
	http.Redirect(w, r, f.successRedirectURL.String(), http.StatusTemporaryRedirect)
}
