package auth

import "context"

type Principal struct {
	UserID int64
	Roles  []string
}

func (p Principal) HasAnyRole(roles ...string) bool {
	for _, have := range p.Roles {
		for _, want := range roles {
			if have == want {
				return true
			}
		}
	}
	return false
}

type principalKey struct{}

func WithPrincipal(c context.Context, p Principal) context.Context {
	return context.WithValue(c, principalKey{}, p)
}

func FromContext(c context.Context) (Principal, bool) {
	p, ok := c.Value(principalKey{}).(Principal)
	return p, ok
}
