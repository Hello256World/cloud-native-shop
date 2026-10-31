package postgres

import (
	"context"
	"fmt"

	"github.com/Hello256World/cloud-native-shop/internal/catalog/adapters/postgres/db"
	"github.com/Hello256World/cloud-native-shop/internal/catalog/domain"
	"github.com/Hello256World/cloud-native-shop/internal/catalog/ports"
	"github.com/Hello256World/cloud-native-shop/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type categoryRepository struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) ports.CategoryRepository {
	return &categoryRepository{
		pool: pool,
	}
}

func (r *categoryRepository) q(c context.Context) *db.Queries {
	return db.New(postgres.Executor(c, r.pool))
}

func (r *categoryRepository) Create(c context.Context, in *domain.Category) (*domain.Category, error) {
	q := r.q(c)

	row, err := q.InsertCategory(c, db.InsertCategoryParams{
		ParentID: in.ParentID, Name: in.Name, Slug: in.Slug, ImageName: in.ImageName, IsActive: in.IsActive,
	})

	if err != nil {
		if constraintContains(err, "slug") {
			return nil, domain.ErrCategorySlugTaken
		}
		return nil, fmt.Errorf("catalog.repo.Category.Create: %w", err)
	}
	return categoryToDomain(row), nil
}

func (r *categoryRepository) GetByID(c context.Context, id int64) (*domain.Category, error) {
	row, err := r.q(c).GetCategoryByID(c, id)

	if err != nil {
		if isNoRows(err) {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, fmt.Errorf("catalog.repo.Category.GetByID id=%d: %w", id, err)
	}

	return categoryToDomain(row), nil
}

func (r *categoryRepository) ListAll(c context.Context) ([]domain.Category, error) {
	rows, err := r.q(c).ListCategories(c)

	if err != nil {
		return nil, fmt.Errorf("catalog.repo.Category.ListAll: %w", err)
	}

	return categoriesToDomain(rows), nil
}

func (r *categoryRepository) ListActive(c context.Context) ([]domain.Category, error) {
	rows, err := r.q(c).ListActiveCategories(c)

	if err != nil {
		return nil, fmt.Errorf("catalog.repo.Category.ListActive: %w", err)
	}

	return categoriesToDomain(rows), nil
}

func (r *categoryRepository) Update(c context.Context, in *domain.Category) (*domain.Category, error) {
	row, err := r.q(c).UpdateCategory(c, db.UpdateCategoryParams{
		ID:        in.ID,
		ParentID:  in.ParentID,
		Name:      in.Name,
		Slug:      in.Slug,
		ImageName: in.ImageName,
		IsActive:  in.IsActive,
	})

	if err != nil {
		switch {
		case isNoRows(err):
			return nil, domain.ErrCategoryNotFound
		case constraintContains(err, "slug"):
			return nil, domain.ErrCategorySlugTaken
		}

		return nil, fmt.Errorf("catalog.repo.Category.Update id=%d: %w", in.ID, err)
	}

	return categoryToDomain(row), err
}

func (r *categoryRepository) Delete(c context.Context, id int64) error {
	n, err := r.q(c).SoftDeleteCategory(c, id)

	if err != nil {
		return fmt.Errorf("catalog.repo.Category.Delete id=%d: %w", id, err)
	}

	if n == 0 {
		return domain.ErrCategoryNotFound
	}

	return nil
}

func (r *categoryRepository) HasChildren(c context.Context, id int64) (bool, error) {
	ok, err := r.q(c).CategoryHasChildren(c, &id)
	if err != nil {
		return false, fmt.Errorf("catalog.repo.Category.HasChildren id=%d: %w", id, err)
	}
	return ok, nil
}

func (r *categoryRepository) HasProducts(c context.Context, id int64) (bool, error) {
	ok, err := r.q(c).CategoryHasProducts(c, id)

	if err != nil {
		return false, fmt.Errorf("catalog.repo.Category.HasProducts id=%d: %w", id, err)
	}

	return ok, nil
}
