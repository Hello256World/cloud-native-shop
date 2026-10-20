package auth

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid or expired access token")

type TokenManager struct {
	secret []byte
	issuer string
	ttl    time.Duration
	now    func() time.Time
}

func NewTokenManager(secret, issuer string, ttl time.Duration) *TokenManager {
	return &TokenManager{secret: []byte(secret), issuer: issuer, ttl: ttl, now: time.Now}
}

type claims struct {
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}

func (m *TokenManager) Issue(userID int64, roles []string) (string, time.Time, error) {
	now := m.now()
	exp := now.Add(m.ttl)

	jti, err := randomID()
	if err != nil {
		return "", time.Time{}, fmt.Errorf("auth.Issue jti: %w", err)
	}

	c := claims{
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userID, 10),
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
			ID:        jti,
		},
	}

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(m.secret)

	if err != nil {
		return "", time.Time{}, fmt.Errorf("auth.Issue sign: %w", err)
	}
	return signed, exp, nil
}

func (m *TokenManager) Verify(token string) (Principal, error) {
	var c claims

	_, err := jwt.ParseWithClaims(token, &c, func(t *jwt.Token) (any, error) {
		return m.secret, nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(m.issuer),
		jwt.WithExpirationRequired(),
		jwt.WithTimeFunc(m.now),
	)
	if err != nil {
		return Principal{}, fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}

	id, err := strconv.ParseInt(c.Subject, 10, 64)
	if err != nil {
		return Principal{}, fmt.Errorf("%w: bad subject", ErrInvalidToken)
	}
	return Principal{UserID: id, Roles: c.Roles}, nil
}
