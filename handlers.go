package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/jsonapi"
	"github.com/lib/pq"
	"github.com/mujz/i++/authentication"
	"github.com/mujz/i++/config"
	"github.com/mujz/i++/validator"
)

type key int

const userIDKey = key(1)
const userKey = key(2)

var (
	auth   *authentication.Authenticator
	fbAuth *authentication.FacebookAuthenticator
	model  Model
)

func init() {
	// Initialize Authenticator
	auth = authentication.NewAuthenticator(
		config.AuthSecretKey,
		config.AuthTokenExpirationInterval,
	)

	// Initialize the model
	model = Model{
		DBName:   config.DBName,
		User:     config.DBUser,
		Password: config.DBPassword,
		Host:     config.DBHost,
		Port:     config.DBPort,
		SSLMode:  config.DBSSLMode,
	}

	if err := model.Open(); err != nil {
		panic(err)
	}

	var Url *url.URL
	Url, err := url.Parse(config.WebURL)
	if err != nil {
		panic(err)
	}
	Url.Path += "/login"
	parameters := url.Values{}
	parameters.Add("error", "Failed to login with Facebook")
	Url.RawQuery = parameters.Encode()

	fbAuth = authentication.NewFacebookAuthenticator(
		config.FBAppID,
		config.FBAppSecret,
		config.BaseURL+"/login/facebook/callback",
		config.WebURL,
		Url.String(),
		nil,
		model,
		auth,
	)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	marshalError(w, http.StatusNotFound, r.URL.Path+" not found")
}

func CurrentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Get user ID from context
		ctx := r.Context()
		userID := ctx.Value(userIDKey).(int)

		// Get the number
		number, err := model.GetNumber(userID)
		if err != nil {
			marshalError(w, http.StatusInternalServerError, "Failed to retrieve user's number")
			log.Println(err)
			return
		}

		// Send number to client
		jsonapi.MarshalOnePayload(w, number)

	case "PUT":
		// Get user ID from context
		ctx := r.Context()
		userID := ctx.Value(userIDKey).(int)

		// Parse the body's JSON
		newNumber := new(Number)
		if err := jsonapi.UnmarshalPayload(r.Body, newNumber); err != nil {
			marshalError(w, http.StatusBadRequest, "Invalid request body.")
			log.Println(err)
			return
		}

		// Validate the number
		if !validator.ValidateNumber(newNumber.Value) {
			marshalError(
				w, http.StatusBadRequest,
				fmt.Sprintf(
					"Number out of range. Please make sure that it's between %d and %d",
					validator.MinNum, validator.MaxNum,
				),
			)
			return
		}

		// Update the number
		number, err := model.UpdateNumber(userID, newNumber.Value)
		if err != nil {
			marshalError(w, http.StatusInternalServerError, "Failed to update number.")
			log.Println(err)
			return
		}

		// Send updated number to client
		jsonapi.MarshalOnePayload(w, number)

	default:
		NotFoundHandler(w, r)
	}
}

func NextHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		NotFoundHandler(w, r)
		return
	}

	// Get user ID from context
	ctx := r.Context()
	userID := ctx.Value(userIDKey).(int)

	// Increment the number
	number, err := model.IncrementNumber(userID)
	if err != nil {
		marshalError(w, http.StatusInternalServerError, "Failed to increment number")
		log.Println(err)
		return
	}

	// Send number to client
	jsonapi.MarshalOnePayload(w, number)
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	ctx := r.Context()
	user := ctx.Value(userKey).(*User)

	// Create user
	token, err := auth.Signup(model, user.Email, user.Password)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			marshalError(w, http.StatusBadRequest, "Email address taken. Did you mean to log in instead?")
		} else {
			marshalError(w, http.StatusInternalServerError, "Failed to create user")
			log.Println(err)
		}
		return
	}

	// Send token to client
	jsonapi.MarshalOnePayload(w, token)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	ctx := r.Context()
	user := ctx.Value(userKey).(*User)

	// Log user in
	token, err := auth.Login(model, user.Email, user.Password)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			marshalError(w, http.StatusUnauthorized, "Password is incorrect.")
		} else if err == sql.ErrNoRows {
			marshalError(w, http.StatusUnauthorized, "Unknown email address. You'll need to sign up first.")
		} else {
			marshalError(w, http.StatusInternalServerError, "Error logging the user in.")
			log.Println(err)
		}
		return
	}

	// Send token to client
	jsonapi.MarshalOnePayload(w, token)
}

type globalHeadersHandler struct {
	handler http.Handler
}

func GlobalHeadersHandler(h http.Handler) http.Handler {
	return globalHeadersHandler{h}
}

func (h globalHeadersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", jsonapi.MediaType)

	// Allow cross domain
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Content-Length")

	// Intercept OPTIONS method
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
		w.WriteHeader(200)
	}
	h.handler.ServeHTTP(w, r)
}

func AuthDecorator(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get auth header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			marshalError(w, http.StatusUnauthorized, "No authentication token was present in request headers.")
			return
		}

		// Split the header into "Bearer" and token
		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			marshalError(w, http.StatusUnauthorized, "Invalid authorization header format; it must be Bearer {token}")
			return
		}

		// Parse userID from authtoken
		id, err := auth.Authenticate(authHeaderParts[1])
		if err != nil {
			marshalError(w, http.StatusUnauthorized, "Authentication token is invalid")
			log.Println(err)
			return
		}

		// Add the userID to the context
		ctx = context.WithValue(ctx, userIDKey, id)

		// Call handler function
		f(w, r.WithContext(ctx))
	}
}

func LoginValidationDecorator(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			NotFoundHandler(w, r)
			return
		}

		fmt.Println("The request body: ", r.Body)

		ctx := r.Context()

		// Parse the body's JSON
		user := new(User)
		if err := jsonapi.UnmarshalPayload(r.Body, user); err != nil {
			marshalError(w, http.StatusBadRequest, "Invalid request body.")
			log.Println(err)
			return
		}

		// Validate email address
		if !validator.ValidateEmail(user.Email) {
			marshalError(w, http.StatusBadRequest, "Invalid email address.")
			return
		}

		// Validate password
		if !validator.ValidatePassword(user.Password) {
			marshalError(w, http.StatusBadRequest,
				"Invalid password. Please make sure it's at least 8 characters long and doesn't contain an unsupported character. Supported characters are \"a-z, A-Z, 0-9, !#$%^&*_+=-}{|/.'`~\"")
			return
		}

		// Add the user to the context
		ctx = context.WithValue(ctx, userKey, user)

		// Call handler function
		f(w, r.WithContext(ctx))
	}
}

var errorTitles = map[int]string{
	http.StatusBadRequest:          "Bad Request",
	http.StatusUnauthorized:        "Unauthorized",
	http.StatusForbidden:           "Forbidden",
	http.StatusNotFound:            "Not Found",
	http.StatusInternalServerError: "Internal Server Error",
}

func marshalError(w http.ResponseWriter, status int, detail string) {
	w.WriteHeader(status)
	jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
		Title:  errorTitles[status],
		Detail: detail,
		Status: strconv.Itoa(status),
	}})
}
