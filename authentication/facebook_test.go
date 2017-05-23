package authentication

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

const (
	appID     = "12345"
	appSecret = "1as2sd34sd5"
)

var (
	auth                      = NewAuthenticator([]byte("my,secret,key"), 5*time.Second)
	extractPathAndQueryRegexp = regexp.MustCompile(`^((http[s]?):\/)?\/?([^:\/\s]+)((\/\w+)*(:\d+)?\/?)`)
)

type model struct{}

func (m model) CreateFacebookUser(fbID string) (int, error) {
	return 123456, nil
}

func (f FacebookAuthenticator) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken: "123456",
		TokenType:   "Bearer",
	}, nil
}

func fbHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			t.Logf("Facebook API called with invalid url %s", r.URL)
			t.Fail()
		}

		if id, ok := q["client_id"]; !ok || id[0] != appID {
			t.Logf("Facebook API: Expected client_id = %s; Got = %s", appID, id)
		}

		redirectURL, ok := q["redirect_uri"]
		if !ok || len(redirectURL) < 1 {
			t.Log("Facebook API: No redirect_uri")
		}

		state, ok := q["state"]
		if !ok || len(state) < 1 {
			t.Log("Facebook API: No state")
		}

		Url, err := url.Parse(redirectURL[0])
		if err != nil {
			t.Logf("Facebook API was passed an invalid redirect URL %s", redirectURL[0])
			t.Fail()
		}

		params := url.Values{}
		params.Add("code", "AQDrdwQcIPsec6NxsjFbf95du5A30npoHwlrRHlkWIU4p_KnGqh0MJV1zLhJ1GHXR3m1nvCJy-gzXub3hJzm_Rf7k0ecABOHd__Wq_U-C7WKykOQ6FtpNJWwe0W-94_IzvNOgq5rr-mxiQG-4AO_iOrSDxn7c_Q6XCG1mKnR0PNlzuThoFI965Qv5QNcZLBVqsk38xu3NHOLCEhxtA8CG5c6xbEG6-MsdFAFC3YHSd-BPdWpaKwiwnrzX0amDZygBBlb-5umSLcCif6Y8Xw-VnatO9_AJQjSRIRJ0pDJc8Lfv8V5XJfZ6Wbt422ROTduugA")
		params.Add("state", state[0])

		Url.RawQuery = params.Encode()

		http.Redirect(w, r, Url.String(), http.StatusTemporaryRedirect)
	}
}

func successHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			t.Logf("Facebook authenticator redirected to an invalid url %s", r.URL)
			t.Fail()
		}

		if len(q["token"]) < 1 {
			t.Logf("No token in callback URL %s", r.URL)
			t.Fail()
		}

		fmt.Println("Done! Done! Done!")
		w.WriteHeader(200)

	}
}

func errorHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			t.Logf("Facebook authenticator redirected to an invalid url %s", r.URL)
			t.Fail()
		}

		if len(q["error"]) < 1 {
			t.Logf("No error in callback URL %s", r.URL)
			t.Fail()
		}
	}
}

func profileHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if token := r.FormValue("access_token"); len(token) < 1 {
			t.Logf("No access token was sent to facebook user profile API %s", r.URL)
			t.Fail()
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"name": "Mujtaba Zuhair Al-Tameemi",
			"id": "1326314725"
		}`))
	}
}

func TestFacebookLogin(t *testing.T) {
	sMux := http.NewServeMux()
	sMux.HandleFunc("/", fbHandler(t))
	sMux.HandleFunc("/success", successHandler(t))
	sMux.HandleFunc("/error", errorHandler(t))
	sMux.HandleFunc("/me", profileHandler(t))

	s := httptest.NewServer(sMux)
	defer s.Close()

	fb := NewFacebookAuthenticator(
		appID,
		appSecret,
		s.URL+"/callback",
		s.URL+"/success",
		s.URL+"/error",
		nil,
		model{},
		auth,
	)

	sMux.HandleFunc("/login", fb.LoginHandler)
	sMux.HandleFunc("/callback", fb.LoginCallbackHandler)

	facebookProfileAPI = s.URL + "/me"
	fb.Endpoint.AuthURL = s.URL
	fb.setAuthURL()

	r, _ := http.NewRequest("GET", "/login", nil)
	w := httptest.NewRecorder()

	sMux.ServeHTTP(w, r)

	redirectURL, ok := w.Header()["Location"]
	if !ok {
		t.Fatalf("Expected redirect, but no Location header was found.\nHeaders: %v", w.Header())
	}

	u, err := url.Parse(redirectURL[0])
	if err != nil || u.Scheme+"://"+u.Host != s.URL {
		t.Fatalf("Expected redirect to %s; Got: %v", s.URL, redirectURL[0])
	}

	url := extractPathAndQueryRegexp.ReplaceAllString(redirectURL[0], "")
	r, _ = http.NewRequest("GET", "/"+url, nil)
	sMux.ServeHTTP(w, r)

	url = extractPathAndQueryRegexp.ReplaceAllString(w.Header()["Location"][0], "")
	r, _ = http.NewRequest("GET", "/"+url, nil)
	sMux.ServeHTTP(w, r)

	url = extractPathAndQueryRegexp.ReplaceAllString(w.Header()["Location"][0], "")
	r, _ = http.NewRequest("GET", "/"+url, nil)
	sMux.ServeHTTP(w, r)
}
