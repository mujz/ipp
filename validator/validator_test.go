package validator

import "testing"

func TestValidateEmail(t *testing.T) {
	if e := "good@email.com"; !ValidateEmail(e) {
		t.Fatalf("Expected %s to be valid", e)
	}
	if e := "two@atsigns@email.com"; ValidateEmail(e) {
		t.Fatalf("Expected %s to be invalid", e)
	}
	if e := "baddomain@e%ail.com"; ValidateEmail(e) {
		t.Fatalf("Expected %s to be invalid", e)
	}
	if e := "baddomain@email.c^m"; ValidateEmail(e) {
		t.Fatalf("Expected %s to be invalid", e)
	}
}

func TestValidatePassword(t *testing.T) {
	if p := "AllGood13!#$%^&*_+=-}{|/.'`~"; !ValidatePassword(p) {
		t.Fatalf("Expected %s to be invalid", p)
	}
	if p := "short"; ValidatePassword(p) {
		t.Fatalf("Expected %s to be invalid", p)
	}
	if p := "inv@lid"; ValidatePassword(p) {
		t.Fatalf("Expected %s to be invalid", p)
	}
	if p := "inv alid"; ValidatePassword(p) {
		t.Fatalf("Expected %s to be invalid", p)
	}
}

func TestValidateNumber(t *testing.T) {
	if n := 1; !ValidateNumber(n) {
		t.Fatalf("Expected %d to be valid", n)
	}
	if n := 2147483647; !ValidateNumber(n) {
		t.Fatalf("Expected %d to be valid", n)
	}
	if n := 0; ValidateNumber(n) {
		t.Fatalf("Expected %d to be invalid", n)
	}
	if n := 2147483648; ValidateNumber(n) {
		t.Fatalf("Expected %d to be invalid", n)
	}
}
