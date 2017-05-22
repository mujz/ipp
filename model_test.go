package main

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/mujz/ipp/config"
)

func TestOpen(t *testing.T) {
	model := Model{
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
	model.Close()
}

func TestCreate(t *testing.T) {
	email := fmt.Sprintf(e, time.Now().UnixNano())
	user, err := model.Create(email, p)
	if err != nil {
		t.Fatal(err)
	}

	if got := user.Username; got != email {
		t.Fatalf("Expected email = %s, Got = %s", email, got)
	}
}

func TestCreateFacebookUser(t *testing.T) {
	fbID := time.Now().UnixNano()
	_, err := model.CreateFacebookUser(strconv.FormatInt(fbID, 10))
	if err != nil {
		t.Fatal(err)
	}
}

func TestGet(t *testing.T) {
	email := fmt.Sprintf(e, time.Now().UnixNano())
	createdUser, err := model.Create(email, p)
	if err != nil {
		t.Fatal(err)
	}

	user, err := model.Get(email)
	if err != nil {
		t.Fatal(err)
	}

	if got := user.Username; got != email {
		t.Fatalf("Expected email = %s, Got = %s", email, got)
	}

	if got := user.ID; got != createdUser.ID {
		t.Fatalf("Expected ID = %d, Got = %d", createdUser.ID, got)
	}
}

func TestGetNumber(t *testing.T) {
	user, err := model.Create(fmt.Sprintf(e, time.Now().UnixNano()), p)
	if err != nil {
		t.Fatal(err)
	}

	number, err := model.GetNumber(user.ID)
	if err != nil {
		t.Fatal(err)
	}

	if expected := 1; number.Value != expected {
		t.Fatalf("Expected Number = %d, Got = %d", expected, number.Value)
	}
}

func TestIncrementNumber(t *testing.T) {
	user, err := model.Create(fmt.Sprintf(e, time.Now().UnixNano()), p)
	if err != nil {
		t.Fatal(err)
	}

	number, err := model.IncrementNumber(user.ID)
	if err != nil {
		t.Fatal(err)
	}

	if expected := 2; number.Value != expected {
		t.Fatalf("Expected Number = %d, Got = %d", expected, number.Value)
	}
}

func TestUpdateNumber(t *testing.T) {
	user, err := model.Create(fmt.Sprintf(e, time.Now().UnixNano()), p)
	if err != nil {
		t.Fatal(err)
	}

	newNumber := 123
	number, err := model.UpdateNumber(user.ID, newNumber)
	if err != nil {
		t.Fatal(err)
	}

	if got := number.Value; got != newNumber {
		t.Fatalf("Expected Number = %d, Got = %d", newNumber, got)
	}
}
