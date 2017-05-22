package authentication

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const (
	secret                    = "my,secret,key"
	expirationPeriodInMinutes = 60
)

type Token struct {
	Subject       string `json:"sub" jsonapi:"primary,token"`
	ExpiresAt     string `jsonapi:"attr,exipres_at"`
	expiresAt     time.Time
	UnixExpiresAt int64  `json:"exp"`
	Type          string `json:"typ,omitempty" jsonapi:"attr,type,omitempty"`
}

type User struct {
	ID       int    `json:"id" jsonapi:"primary,User"`
	Username string `json:"username" jsonapi:"attr,email"`
	Password string `json:"password" jsonapi:"attr,password"`
}

type UserGetter interface {
	// Takes username (or email) and returns a user
	Get(string) (User, error)
}

type UserCreator interface {
	// Takes username (or email) and password and returns user
	Create(string, string) (User, error)
}

type Authenticator struct {
	secret             []byte        // Secret key
	expirationInterval time.Duration // Duration of time before token exires
}

func NewAuthenticator(secret []byte, expirationInterval time.Duration) *Authenticator {
	return &Authenticator{secret, expirationInterval}
}

func (a *Authenticator) Authenticate(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return a.secret, nil
	})

	if err != nil {
		return -1, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return strconv.Atoi(claims["sub"].(string))
	} else {
		return -1, fmt.Errorf("Token invalid")
	}
}

func (a *Authenticator) Login(model UserGetter, username, password string) (*Token, error) {
	// Get the user from the model
	user, err := model.Get(username)
	if err != nil {
		return nil, err
	}

	// Validate password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	// Generate token
	return a.generateToken(user.ID)
}

func (a *Authenticator) Signup(model UserCreator, username, password string) (*Token, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user, err := model.Create(username, string(hashedPassword))
	if err != nil {
		return nil, err
	}

	// Generate token
	return a.generateToken(user.ID)
}

func (a *Authenticator) generateToken(id int) (*Token, error) {
	expiresAt := time.Now().Add(a.expirationInterval)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Token{
		Subject:       strconv.Itoa(id),
		UnixExpiresAt: expiresAt.Unix(),
	})

	signedToken, err := token.SignedString(a.secret)
	if err != nil {
		return nil, err
	}

	return &Token{
		Subject:   signedToken,
		ExpiresAt: expiresAt.Format(time.RFC3339),
		expiresAt: expiresAt,
	}, nil
}

func (t Token) Valid() error {
	err := new(jwt.ValidationError)

	fmt.Printf("Now: %s\nExp: %s", time.Now(), t.expiresAt)
	if t.expiresAt.Before(time.Now()) {
		err.Inner = fmt.Errorf("token is expired")
		err.Errors |= jwt.ValidationErrorExpired
	}

	if err.Errors == 0 {
		return nil
	}

	return err
}
