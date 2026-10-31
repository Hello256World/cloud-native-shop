package usecase

import (
	"context"
	"errors"

	"github.com/Hello256World/cloud-native-shop/internal/catalog/domain"
)

type CategoryInput struct {
	ParentID  *int64
	Name      string
	Slug      string
	ImageName *string
	IsActive  bool
}

func (c CategoryInput) toDomain(id int64) *domain.Category {
	return &domain.Category{
		ID:        id,
		ParentID:  c.ParentID,
		Name:      c.Name,
		Slug:      c.Slug,
		ImageName: c.ImageName,
		IsActive:  c.IsActive,
	}
}

func (s *Service) CreateCategory(c context.Context, in CategoryInput) (*domain.Category, error) {
	if err := s.checkParent(c, 0, in.ParentID); err != nil {
		return nil, err
	}

	return s.categories.Create(c, in.toDomain(0))
}

func (s *Service) UpdateCategory(c context.Context, id int64, in CategoryInput) (*domain.Category, error) {
	if _, err := s.categories.GetByID(c, id); err != nil {
		return nil, err
	}
	if err := s.checkParent(c, id, in.ParentID); err != nil {
		return nil, err
	}
	return s.categories.Update(c, in.toDomain(id))
}

func (s *Service) DeleteCategory(c context.Context, id int64) error {
	return s.tx.WithinTx(c, func(ctx context.Context) error { 
		if _, err := s.categories.GetByID(ctx, id); err != nil {
			return err
		}

		if has, err := s.categories.HasChildren(ctx, id); err != nil {
			return err
		} else if has {
			return domain.ErrCategoryHasChildren
		}

		if has, err := s.categories.HasProducts(ctx, id); err != nil {
			return err
		} else if has {
			return domain.ErrCategoryHasChildren
		}

		return s.categories.Delete(ctx, id)
	})
}

func (s *Service) ListCategories(c context.Context) ([]domain.Category, error) {
	return s.categories.ListAll(c)
}

func (s *Service) ListPublicCategories(c context.Context) ([]domain.Category, error) { // TODO: ino Bekhon
	all, err := s.categories.ListActive(c)
	if err != nil {
		return nil, err
	}

	byID := make(map[int64]domain.Category, len(all))
	for _, cat := range all {
		byID[cat.ID] = cat
	}

	out := make([]domain.Category, 0, len(all))
	for _, cat := range all {
		visible := true
		for p := cat.ParentID; p != nil; {
			parent, ok := byID[*p]
			if !ok {
				visible = false
				break
			}
			p = parent.ParentID
		}
		if visible {
			out = append(out, cat)
		}
	}
	return out, nil
}

func (s *Service) checkParent(c context.Context, selfID int64, parentID *int64) error {
	for cur := parentID; cur != nil; {
		if *cur == selfID {
			return domain.ErrCategoryInvalidParent
		}

		p, err := s.categories.GetByID(c, *cur)
		if err != nil {
			if errors.Is(err, domain.ErrCategoryNotFound) {
				return domain.ErrCategoryInvalidParent
			}
			return err
		}
		cur = p.ParentID
	}
	return nil
}
