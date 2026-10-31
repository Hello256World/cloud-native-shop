-- name: InsertCategory :one
INSERT INTO catalog.categories (parent_id, name, slug, image_name, is_active)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetCategoryByID :one
SELECT * FROM catalog.categories
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListCategories :many
SELECT * FROM catalog.categories
WHERE deleted_at IS NULL
ORDER BY id;

-- name: ListActiveCategories :many
SELECT * FROM catalog.categories
WHERE deleted_at IS NULL AND is_active
ORDER BY id;

-- name: UpdateCategory :one
UPDATE catalog.categories
SET parent_id = $2, name = $3, slug = $4, image_name = $5, is_active = $6
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteCategory :execrows
UPDATE catalog.categories SET deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: CategoryHasChildren :one
SELECT EXISTS (
  SELECT 1 FROM catalog.categories WHERE parent_id = $1 AND deleted_at IS NULL
);

-- name: CategoryHasProducts :one
SELECT EXISTS (
  SELECT 1 FROM catalog.products WHERE category_id = $1 AND deleted_at IS NULL
);