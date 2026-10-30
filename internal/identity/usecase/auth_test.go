package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Hello256World/cloud-native-shop/internal/identity/domain"
)

// ---------- fakes ----------

type fakeUsers struct {
	byPhone map[string]*domain.User
	nextID  int64
}

func newFakeUsers() *fakeUsers { return &fakeUsers{byPhone: map[string]*domain.User{}} }

func (f *fakeUsers) Create(_ context.Context, u *domain.User) (*domain.User, error) {
	if _, ok := f.byPhone[u.Phone]; ok {
		return nil, domain.ErrPhoneTaken
	}
	f.nextID++
	c := *u
	c.ID = f.nextID
	f.byPhone[c.Phone] = &c
	return &c, nil
}

func (f *fakeUsers) GetByID(_ context.Context, id int64) (*domain.User, error) {
	for _, u := range f.byPhone {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

func (f *fakeUsers) GetByPhone(_ context.Context, phone string) (*domain.User, error) {
	if u, ok := f.byPhone[phone]; ok {
		return u, nil
	}
	return nil, domain.ErrUserNotFound
}

type fakeRefresh struct {
	byHash map[string]*domain.RefreshToken
	nextID int64
}

func newFakeRefresh() *fakeRefresh { return &fakeRefresh{byHash: map[string]*domain.RefreshToken{}} }

func (f *fakeRefresh) Create(_ context.Context, userID int64, hash string, exp time.Time) error {
	f.nextID++
	f.byHash[hash] = &domain.RefreshToken{ID: f.nextID, UserID: userID, TokenHash: hash, ExpiresAt: exp}
	return nil
}

func (f *fakeRefresh) GetByHash(_ context.Context, hash string) (*domain.RefreshToken, error) {
	if t, ok := f.byHash[hash]; ok {
		return t, nil
	}
	return nil, domain.ErrRefreshTokenInvalid
}

func (f *fakeRefresh) Revoke(_ context.Context, id int64) error {
	for _, t := range f.byHash {
		if t.ID == id && t.RevokedAt == nil {
			now := time.Now()
			t.RevokedAt = &now
		}
	}
	return nil
}

type fakeHasher struct{}

func (fakeHasher) Hash(p string) (string, error)    { return "h:" + p, nil }
func (fakeHasher) Verify(p, h string) (bool, error) { return h == "h:"+p, nil }

type fakeIssuer struct{}

func (fakeIssuer) Issue(id int64, _ []domain.Role) (string, time.Time, error) {
	return "access", time.Now().Add(15 * time.Minute), nil
}

type fakeTx struct{}

func (fakeTx) WithinTx(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }

// ---------- helpers ----------

var t0 = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

const refreshTTL = 30 * 24 * time.Hour

func newTestService(now *time.Time) (*Service, *fakeUsers) {
	users := newFakeUsers()
	s := NewService(users, newFakeRefresh(), nil, fakeHasher{}, fakeIssuer{}, fakeTx{}, refreshTTL)
	s.now = func() time.Time { return *now }
	return s, users
}

func register(t *testing.T, s *Service) {
	t.Helper()
	_, err := s.Register(context.Background(), RegisterInput{
		FullName: "Ali", Phone: "09120000000", Password: "12345678",
	})
	if err != nil {
		t.Fatal(err)
	}
}

// ---------- tests ----------

func TestRegister(t *testing.T) {
	now := t0
	s, users := newTestService(&now)
	register(t, s)

	u := users.byPhone["09120000000"]
	if !u.IsActive || !u.HasRole(domain.RoleCustomer) || u.HasRole(domain.RoleAdmin) {
		t.Fatalf("unexpected user: %+v", u)
	}
	if u.PasswordHash == nil || *u.PasswordHash == "12345678" {
		t.Fatal("password must be stored hashed")
	}

	_, err := s.Register(context.Background(), RegisterInput{FullName: "X", Phone: "09120000000", Password: "abcdefgh"})
	if !errors.Is(err, domain.ErrPhoneTaken) {
		t.Fatalf("want ErrPhoneTaken, got %v", err)
	}
}

func TestLoginErrors(t *testing.T) {
	now := t0
	s, users := newTestService(&now)
	register(t, s)
	ctx := context.Background()

	// unknown phone and wrong password must be indistinguishable
	if _, err := s.Login(ctx, "09999999999", "12345678"); !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("unknown phone: got %v", err)
	}
	if _, err := s.Login(ctx, "09120000000", "wrong"); !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("wrong password: got %v", err)
	}

	// inactive is checked AFTER the password
	users.byPhone["09120000000"].IsActive = false
	if _, err := s.Login(ctx, "09120000000", "wrong"); !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("inactive + wrong password must not leak: got %v", err)
	}
	if _, err := s.Login(ctx, "09120000000", "12345678"); !errors.Is(err, domain.ErrUserInactive) {
		t.Fatalf("inactive: got %v", err)
	}
}

func TestRefreshKeepsAbsoluteExpiry(t *testing.T) {
	now := t0
	s, _ := newTestService(&now)
	register(t, s)
	ctx := context.Background()

	pair, err := s.Login(ctx, "09120000000", "12345678")
	if err != nil {
		t.Fatal(err)
	}

	// day 29: still valid; the new token must NOT extend the session
	now = t0.Add(29 * 24 * time.Hour)
	p2, err := s.Refresh(ctx, pair.RefreshToken)
	if err != nil {
		t.Fatal(err)
	}
	if !p2.RefreshExpiresAt.Equal(pair.RefreshExpiresAt) {
		t.Fatalf("expiry extended: %v -> %v", pair.RefreshExpiresAt, p2.RefreshExpiresAt)
	}

	// old token was rotated (revoked)
	if _, err := s.Refresh(ctx, pair.RefreshToken); !errors.Is(err, domain.ErrRefreshTokenInvalid) {
		t.Fatalf("reused token: got %v", err)
	}

	// day 31: the session hit its absolute cap
	now = t0.Add(31 * 24 * time.Hour)
	if _, err := s.Refresh(ctx, p2.RefreshToken); !errors.Is(err, domain.ErrRefreshTokenInvalid) {
		t.Fatalf("expired session: got %v", err)
	}
}

func TestLogoutIsIdempotent(t *testing.T) {
	now := t0
	s, _ := newTestService(&now)
	register(t, s)
	ctx := context.Background()

	pair, _ := s.Login(ctx, "09120000000", "12345678")
	for i := 0; i < 2; i++ {
		if err := s.Logout(ctx, pair.RefreshToken); err != nil {
			t.Fatalf("logout #%d: %v", i+1, err)
		}
	}
	if err := s.Logout(ctx, "never-existed"); err != nil {
		t.Fatalf("unknown token must be a no-op: %v", err)
	}
	if _, err := s.Refresh(ctx, pair.RefreshToken); !errors.Is(err, domain.ErrRefreshTokenInvalid) {
		t.Fatalf("revoked token must not refresh: got %v", err)
	}
}

func TestEnsureSuperAdminIsIdempotent(t *testing.T) {
	now := t0
	s, users := newTestService(&now)
	in := RegisterInput{FullName: "Root", Phone: "09121111111", Password: "a-long-password"}

	created, err := s.EnsureSuperAdmin(context.Background(), in)
	if err != nil || !created {
		t.Fatalf("first call: created=%v err=%v", created, err)
	}
	u := users.byPhone[in.Phone]
	if !u.HasRole(domain.RoleAdmin) || !u.HasRole(domain.RoleSuperAdmin) {
		t.Fatalf("missing admin roles: %+v", u.Roles)
	}

	created, err = s.EnsureSuperAdmin(context.Background(), in)
	if err != nil || created {
		t.Fatalf("second call: created=%v err=%v", created, err)
	}
}
