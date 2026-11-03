package usecase

import (
	"context"
	"errors"
	"sort"
	"testing"

	"github.com/Hello256World/cloud-native-shop/internal/catalog/domain"
)

type fakeCategories struct {
	byID        map[int64]*domain.Category
	nextID      int64
	withProduct map[int64]bool
}

func newFakeCategories() *fakeCategories {
	return &fakeCategories{byID: map[int64]*domain.Category{}, withProduct: map[int64]bool{}}
}

func (f *fakeCategories) Create(_ context.Context, c *domain.Category) (*domain.Category, error) {
	f.nextID++
	cp := *c
	cp.ID = f.nextID
	f.byID[cp.ID] = &cp
	return &cp, nil
}

func (f *fakeCategories) GetByID(_ context.Context, id int64) (*domain.Category, error) {
	if c, ok := f.byID[id]; ok {
		return c, nil
	}
	return nil, domain.ErrCategoryNotFound
}

func (f *fakeCategories) list(onlyActive bool) []domain.Category {
	out := []domain.Category{}
	for _, c := range f.byID {
		if !onlyActive || c.IsActive {
			out = append(out, *c)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (f *fakeCategories) ListAll(context.Context) ([]domain.Category, error) {
	return f.list(false), nil
}
func (f *fakeCategories) ListActive(context.Context) ([]domain.Category, error) {
	return f.list(true), nil
}

func (f *fakeCategories) Update(_ context.Context, c *domain.Category) (*domain.Category, error) {
	cp := *c
	f.byID[c.ID] = &cp
	return &cp, nil
}

func (f *fakeCategories) Delete(_ context.Context, id int64) error {
	delete(f.byID, id)
	return nil
}

func (f *fakeCategories) HasChildren(_ context.Context, id int64) (bool, error) {
	for _, c := range f.byID {
		if c.ParentID != nil && *c.ParentID == id {
			return true, nil
		}
	}
	return false, nil
}

func (f *fakeCategories) HasProducts(_ context.Context, id int64) (bool, error) {
	return f.withProduct[id], nil
}

type fakeTx struct{}

func (fakeTx) WithinTx(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }

func newSvc() (*Service, *fakeCategories) {
	repo := newFakeCategories()
	return NewService(repo, nil, nil, fakeTx{}), repo
}

func mk(t *testing.T, s *Service, name string, parent *int64, active bool) *domain.Category {
	t.Helper()
	c, err := s.CreateCategory(context.Background(), CategoryInput{
		Name: name, Slug: name, ParentID: parent, IsActive: active,
	})
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestCreateCategory_ParentMustExist(t *testing.T) {
	s, _ := newSvc()
	missing := int64(99)
	_, err := s.CreateCategory(context.Background(), CategoryInput{Name: "x", Slug: "x", ParentID: &missing})
	if !errors.Is(err, domain.ErrCategoryInvalidParent) {
		t.Fatalf("got %v", err)
	}
}

func TestUpdateCategory_RejectsCycles(t *testing.T) {
	s, _ := newSvc()
	ctx := context.Background()
	a := mk(t, s, "a", nil, true)
	b := mk(t, s, "b", &a.ID, true)
	c := mk(t, s, "c", &b.ID, true)

	for name, parent := range map[string]int64{"self": a.ID, "child": b.ID, "grandchild": c.ID} {
		p := parent
		_, err := s.UpdateCategory(ctx, a.ID, CategoryInput{Name: "a", Slug: "a", ParentID: &p, IsActive: true})
		if !errors.Is(err, domain.ErrCategoryInvalidParent) {
			t.Fatalf("%s: got %v", name, err)
		}
	}

	if _, err := s.UpdateCategory(ctx, 12345, CategoryInput{Name: "x", Slug: "x"}); !errors.Is(err, domain.ErrCategoryNotFound) {
		t.Fatalf("missing id must be not found, got %v", err)
	}
}

func TestDeleteCategory_Guards(t *testing.T) {
	s, repo := newSvc()
	ctx := context.Background()
	a := mk(t, s, "a", nil, true)
	b := mk(t, s, "b", &a.ID, true)

	if err := s.DeleteCategory(ctx, a.ID); !errors.Is(err, domain.ErrCategoryHasChildren) {
		t.Fatalf("with child: got %v", err)
	}

	repo.withProduct[b.ID] = true
	if err := s.DeleteCategory(ctx, b.ID); !errors.Is(err, domain.ErrCategoryHasChildren) {
		t.Fatalf("with products: got %v", err)
	}

	repo.withProduct[b.ID] = false
	if err := s.DeleteCategory(ctx, b.ID); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteCategory(ctx, a.ID); err != nil {
		t.Fatalf("empty now: %v", err)
	}
	if err := s.DeleteCategory(ctx, a.ID); !errors.Is(err, domain.ErrCategoryNotFound) {
		t.Fatalf("twice: got %v", err)
	}
}

func TestListPublic_HidesChildrenOfInactive(t *testing.T) {
	s, _ := newSvc()
	root := mk(t, s, "root", nil, false) // inactive
	mk(t, s, "child", &root.ID, true)    // active, but parent is hidden
	mk(t, s, "other", nil, true)

	got, err := s.ListPublicCategories(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Name != "other" {
		t.Fatalf("unexpected: %+v", got)
	}
}
