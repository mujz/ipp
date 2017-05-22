/* Integration Tests */
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/jsonapi"
	"github.com/mujz/ipp/authentication"
	"github.com/mujz/ipp/config"
	"github.com/mujz/ipp/util/testutil"
)

var secret = config.AuthSecretKey

const (
	e = "jd%v@m.ca" // email template
	p = "mpa$$!23"  // password
)

/* --- Test Signup --- */

func TestSuccessfulSignup(t *testing.T) {
	testAuth(t, "/signup", fmt.Sprintf(e, time.Now().UnixNano()), p, http.StatusOK)
}

func TestAnotherSuccessfulSignup(t *testing.T) {
	testAuth(t, "/signup", fmt.Sprintf("j!#$&'*+-/=?^_`{|}~.01UP%v@m.ca", time.Now().UnixNano()), p, http.StatusOK)
}

func TestFailedSignupNoPassword(t *testing.T) {
	testAuth(t, "/signup", fmt.Sprintf(e, time.Now().UnixNano()), "", http.StatusBadRequest)
}

func TestFailedSignupEmailWithoutDomain(t *testing.T) {
	testAuth(t, "/signup", fmt.Sprintf("jd%v", time.Now().UnixNano()), p, http.StatusBadRequest)
}

func TestFailedSignupEmailHasTwoConsecutivePeriods(t *testing.T) {
	testAuth(t, "/signup", fmt.Sprintf("jd..%v@m.ca", time.Now().UnixNano()), "", http.StatusBadRequest)
}

func TestFailedSignupMalformattedDomainName(t *testing.T) {
	testAuth(t, "/signup", fmt.Sprintf("jd%v@bad.^_^", time.Now().UnixNano()), p, http.StatusBadRequest)
}

func TestFailedSignupEmptyEmail(t *testing.T) {
	testAuth(t, "/signup", "", p, http.StatusBadRequest)
}

func TestFailedSignupEmailAlreadyExists(t *testing.T) {
	email := fmt.Sprintf(e, time.Now().UnixNano())
	testAuth(t, "/signup", email, p, http.StatusOK)
	testAuth(t, "/signup", email, p, http.StatusBadRequest)
}

/* --- Test Login --- */

func TestSuccessfulLogin(t *testing.T) {
	email := fmt.Sprintf(e, time.Now().UnixNano())
	testAuth(t, "/signup", email, p, http.StatusOK)
	testAuth(t, "/login", email, p, http.StatusOK)
}

func TestFailedLoginWrongPassword(t *testing.T) {
	email := fmt.Sprintf(e, time.Now().UnixNano())
	testAuth(t, "/signup", email, p, http.StatusOK)
	testAuth(t, "/login", email, "wrongPassword", http.StatusUnauthorized)
}

func TestFailedLoginWrongEmail(t *testing.T) {
	email := fmt.Sprintf(e, time.Now().UnixNano())
	testAuth(t, "/signup", email, p, http.StatusOK)
	testAuth(t, "/login", "wrong@email.ca", p, http.StatusUnauthorized)
}

/* --- Test Current Get --- */

func TestSuccessfulGet(t *testing.T) {
	token := testAuth(t, "/signup", fmt.Sprintf(e, time.Now().UnixNano()), p, http.StatusOK)
	testCurrentGet(t, token, 1, http.StatusOK)
}

func TestUnauthorizedGet(t *testing.T) {
	testCurrentGet(t, "", 0, http.StatusUnauthorized)
}

/* --- Test Current Update --- */

func TestSuccessfulUpdate(t *testing.T) {
	token := testAuth(t, "/signup", fmt.Sprintf(e, time.Now().UnixNano()), p, http.StatusOK)
	testCurrentUpdate(t, token, 1, http.StatusOK)
}

func TestUnauthorizedUpdate(t *testing.T) {
	testCurrentUpdate(t, "", 0, http.StatusUnauthorized)
}

func TestNegativeUpdate(t *testing.T) {
	token := testAuth(t, "/signup", fmt.Sprintf(e, time.Now().UnixNano()), p, http.StatusOK)
	testCurrentUpdate(t, token, -3, http.StatusBadRequest)
}

func TestOutOfRangeUpdate(t *testing.T) {
	token := testAuth(t, "/signup", fmt.Sprintf(e, time.Now().UnixNano()), p, http.StatusOK)
	testCurrentUpdate(t, token, 2147483649, http.StatusBadRequest)
}

func TestNotNumberUpdate(t *testing.T) {
	token := testAuth(t, "/signup", fmt.Sprintf(e, time.Now().UnixNano()), p, http.StatusOK)
	testCurrentUpdate(t, token, "a", http.StatusBadRequest)
}

/* --- Test Increment --- */

func TestSuccessfulIncrement(t *testing.T) {
	token := testAuth(t, "/signup", fmt.Sprintf(e, time.Now().UnixNano()), p, http.StatusOK)
	testNext(t, token, 2, http.StatusOK)
	testNext(t, token, 3, http.StatusOK)
	testNext(t, token, 4, http.StatusOK)
}

func TestUnauthorizedIncrement(t *testing.T) {
	testNext(t, "", 0, http.StatusUnauthorized)
}

/* --- Test 404 --- */

func Test404(t *testing.T) {
	s := NewServer()
	r, _ := http.NewRequest("GET", "/not-found", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)

	// Assert Headers
	checkHeaders(w, http.StatusNotFound, t)
}

/* --- Helper Methods --- */

func testAuth(t *testing.T, path, email, password string, status int) string {
	/* Seting up test */
	s := NewServer()

	user := &User{
		Email:    email,
		Password: password,
	}
	body := testutil.Body{}
	if err := jsonapi.MarshalOnePayload(body, user); err != nil {
		t.Fatalf("Failed to marshal jsonapi request body: %s", err)
	}

	/* Running test */
	r, _ := http.NewRequest("POST", path, body)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, r)

	// Assert Headers
	checkHeaders(w, status, t)

	if status == http.StatusOK {
		/* Reading result */
		token := new(authentication.Token)
		if err := jsonapi.UnmarshalPayload(w.Body, token); err != nil {
			t.Fatalf("Failed to unmarshal response payload: %s\nReceived Body:%s\n", err, w.Body.String())
		}
		// assert token expiration date
		now := time.Now()
		exp, err := time.Parse(time.RFC3339, token.ExpiresAt)
		if err != nil {
			t.Logf(
				"Token date format unrecognized \nExpected Format: %s\nGot:%s\n",
				time.RFC3339,
				token.ExpiresAt,
			)
			t.Fail()
		} else if exp.Before(now) {
			t.Logf(
				"Token already expired\nNow: %s\nExpires:%s\n",
				now.Format(time.RFC3339),
				token.ExpiresAt,
			)
			t.Fail()
		}

		// assert token
		parsedToken, err := jwt.Parse(token.Subject, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(secret), nil
		})
		if err != nil {
			t.Fatalf("Failed to decode token: %s\n", err)
		}

		claims := parsedToken.Claims.(jwt.MapClaims)
		if err != nil {
			t.Fatalf("Failed getting claims from parsed token: %s\n", err)
		}

		if expClaim := claims["exp"].(float64); expClaim != float64(exp.Unix()) {
			t.Logf(
				"Expiration time mismatch between token and response attributes:\ntoken[\"exp\"]: %s\nattributes[\"ExpiresAt\"]:%s\n",
				expClaim, float64(exp.Unix()),
			)
			t.Fail()
		}

		if _, err := strconv.Atoi(claims["sub"].(string)); err != nil {
			t.Fatalf("User ID must be an integer; got %s\n", claims["sub"])
		}

		return token.Subject
	}
	return ""
}

func testCurrentGet(t *testing.T, token string, expected, status int) {
	/* Seting up test */
	s := NewServer()

	/* Running test */
	r, _ := http.NewRequest("GET", "/current", nil)
	r.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, r)

	// Assert headers
	checkHeaders(w, status, t)

	if status == http.StatusOK {
		/* Reading result */
		number := new(Number)
		if err := jsonapi.UnmarshalPayload(w.Body, number); err != nil {
			t.Fatalf("Failed to unmarshal response payload: %s\nReceived Body:%s\n", err, w.Body.String())
		}

		if got := number.Value; got != expected {
			t.Fatalf("Expected number = %d, but got = %d\n", expected, got)
		}
	}
}

func testCurrentUpdate(t *testing.T, token string, newValue interface{}, status int) {
	/* Seting up test */
	s := NewServer()

	type TestNumber struct {
		Value interface{} `jsonapi:"attr,value"`
	}

	number := &TestNumber{
		Value: newValue,
	}

	body := testutil.Body{}
	if err := jsonapi.MarshalOnePayload(body, number); err != nil {
		t.Fatalf("Failed to marshal jsonapi request body: %s", err)
	}

	/* Running test */
	r, _ := http.NewRequest("PUT", "/current", body)
	r.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, r)

	// Assert headers
	checkHeaders(w, status, t)

	if status == http.StatusOK {
		/* Reading result */
		number := new(Number)
		if err := jsonapi.UnmarshalPayload(w.Body, number); err != nil {
			t.Fatalf("Failed to unmarshal response payload: %s\nReceived Body:%s\n", err, w.Body.String())
		}

		if got := number.Value; got != newValue {
			t.Fatalf("Expected number = %d, but got = %d\n", newValue, got)
		}
	}
}

func testNext(t *testing.T, token string, expected, status int) {
	/* Seting up test */
	s := NewServer()

	/* Running test */
	r, _ := http.NewRequest("GET", "/next", nil)
	r.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, r)

	// Assert headers
	checkHeaders(w, status, t)

	if status == http.StatusOK {
		/* Reading result */
		number := new(Number)
		if err := jsonapi.UnmarshalPayload(w.Body, number); err != nil {
			t.Fatalf("Failed to unmarshal response payload: %s\nReceived Body:%s\n", err, w.Body.String())
		}

		if got := number.Value; got != expected {
			t.Fatalf("Expected number = %d, but got = %d\n", expected, got)
		}
	}
}

func checkHeaders(w *httptest.ResponseRecorder, status int, t *testing.T) {
	headers := w.Header()
	// Assert content-type
	if h := headers["Content-Type"][0]; h != "application/vnd.api+json" {
		t.Logf("\"Content-Type\" header must be application/vnd.api+json; but got %q\n", h)
		t.Fail()
	}
	// assert status
	if w.Code != status {
		t.Fatalf(
			"Expected status (%d); Got (%d):\nHeaders: {\n%q\n}\nBody: {\n%q}\n",
			status, w.Code, headers, w.Body.String(),
		)
	}
}
