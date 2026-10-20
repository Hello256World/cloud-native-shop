package auth

import (
	"testing"
	"time"
)

const testSecret = "0123456789abcdef0123456789abcdef"

func TestIssueVerify(t *testing.T) {
	m := NewTokenManager(testSecret, "test", 15*time.Minute)

	tok, _, err := m.Issue(42, []string{"customer"})
	if err != nil {
		t.Fatal(err)
	}
	p, err := m.Verify(tok)
	if err != nil {
		t.Fatal(err)
	}
	if p.UserID != 42 || !p.HasAnyRole("customer") {
		t.Fatalf("unexpected principal: %+v", p)
	}
}

func TestVerifyRejects(t *testing.T) {
	m := NewTokenManager(testSecret, "test", 15*time.Minute)
	tok, _, _ := m.Issue(1, []string{"customer"})

	t.Run("wrong secret", func(t *testing.T) {
		other := NewTokenManager("ffffffffffffffffffffffffffffffff", "test", time.Minute)
		if _, err := other.Verify(tok); err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("wrong issuer", func(t *testing.T) {
		other := NewTokenManager(testSecret, "someone-else", time.Minute)
		if _, err := other.Verify(tok); err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("expired", func(t *testing.T) {
		past := NewTokenManager(testSecret, "test", time.Minute)
		past.now = func() time.Time { return time.Now().Add(-time.Hour) }
		old, _, _ := past.Issue(1, nil)
		if _, err := m.Verify(old); err == nil {
			t.Fatal("expected expiry error")
		}
	})

	t.Run("garbage", func(t *testing.T) {
		if _, err := m.Verify("not.a.jwt"); err == nil {
			t.Fatal("expected error")
		}
	})
}