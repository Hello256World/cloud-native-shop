-- name: InsertProduct :one
INSERT INTO catalog.products (category_id, name, slug, description, thumbnail, shipment_weight_g, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetProductByID :one
SELECT * FROM catalog.products WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateProduct :one
UPDATE catalog.products
SET category_id = $2, name = $3, slug = $4, description = $5,
    thumbnail = $6, shipment_weight_g = $7, is_active = $8
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteProduct :execrows
UPDATE catalog.products SET deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: InsertVariant :one
INSERT INTO catalog.product_variants (product_id, sku, price, is_active)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateVariant :one
UPDATE catalog.product_variants SET price = $3, is_active = $4
WHERE id = $1 AND product_id = $2
RETURNING *;

-- name: InsertVariantOption :exec
INSERT INTO catalog.product_variant_values (product_variant_id, variant_value_id)
VALUES ($1, $2);

-- name: ListVariantsByProduct :many
SELECT * FROM catalog.product_variants WHERE product_id = $1 ORDER BY id;

-- name: ListVariantOptionsByProduct :many
SELECT pvv.product_variant_id,
       vt.id   AS type_id,
       vt.name AS type_name,
       vv.id   AS value_id,
       vv.value
FROM catalog.product_variant_values pvv
JOIN catalog.product_variants pv ON pv.id = pvv.product_variant_id
JOIN catalog.variant_values  vv ON vv.id = pvv.variant_value_id
JOIN catalog.variant_types   vt ON vt.id = vv.variant_type_id
WHERE pv.product_id = $1
ORDER BY pvv.product_variant_id, vt.display_order, vt.id;

-- name: InsertImage :one
INSERT INTO catalog.product_images (product_id, image_name, priority, is_active)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteImage :execrows
DELETE FROM catalog.product_images WHERE id = $1 AND product_id = $2;

-- name: ListImagesByProduct :many
SELECT * FROM catalog.product_images WHERE product_id = $1 ORDER BY priority, id;

-- name: GetProductBySlug :one
SELECT * FROM catalog.products WHERE slug = $1 AND deleted_at IS NULL;

-- name: ListPublicProducts :many
SELECT p.id, p.category_id, p.name, p.slug, p.thumbnail, p.is_active,
       COALESCE(p.min_price, 0)::bigint AS min_price
FROM catalog.products p
JOIN catalog.categories c
  ON c.id = p.category_id AND c.is_active AND c.deleted_at IS NULL
WHERE p.deleted_at IS NULL
  AND p.is_active
  AND p.min_price IS NOT NULL
  AND (sqlc.arg(after_id)::bigint = 0 OR p.id < sqlc.arg(after_id)::bigint)
  AND (sqlc.arg(category_id)::bigint = 0 OR p.category_id = sqlc.arg(category_id)::bigint)
  AND (sqlc.arg(search)::text = '' OR p.name ILIKE '%' || sqlc.arg(search)::text || '%' ESCAPE '\')
  AND (sqlc.arg(price_min)::bigint = 0 OR p.min_price >= sqlc.arg(price_min)::bigint)
  AND (sqlc.arg(price_max)::bigint = 0 OR p.min_price <= sqlc.arg(price_max)::bigint)
  AND (sqlc.arg(attr_id)::bigint = 0 OR EXISTS (
    SELECT 1 FROM catalog.product_attributes pa
    WHERE pa.product_id = p.id
      AND pa.attribute_id = sqlc.arg(attr_id)::bigint
      AND (pa.value_number = sqlc.narg(attr_number)::numeric
           OR pa.value_text = sqlc.narg(attr_text)::text)
  ))
ORDER BY p.id DESC
LIMIT sqlc.arg(page_size)::int;

-- name: ListPublicProductsByPriceAsc :many
SELECT p.id, p.category_id, p.name, p.slug, p.thumbnail, p.is_active,
       COALESCE(p.min_price, 0)::bigint AS min_price
FROM catalog.products p
JOIN catalog.categories c
  ON c.id = p.category_id AND c.is_active AND c.deleted_at IS NULL
WHERE p.deleted_at IS NULL
  AND p.is_active
  AND p.min_price IS NOT NULL
  AND (sqlc.arg(after_id)::bigint = 0
       OR (p.min_price, p.id) > (sqlc.arg(after_price)::bigint, sqlc.arg(after_id)::bigint))
  AND (sqlc.arg(category_id)::bigint = 0 OR p.category_id = sqlc.arg(category_id)::bigint)
  AND (sqlc.arg(search)::text = '' OR p.name ILIKE '%' || sqlc.arg(search)::text || '%' ESCAPE '\')
  AND (sqlc.arg(price_min)::bigint = 0 OR p.min_price >= sqlc.arg(price_min)::bigint)
  AND (sqlc.arg(price_max)::bigint = 0 OR p.min_price <= sqlc.arg(price_max)::bigint)
  AND (sqlc.arg(attr_id)::bigint = 0 OR EXISTS (
    SELECT 1 FROM catalog.product_attributes pa
    WHERE pa.product_id = p.id
      AND pa.attribute_id = sqlc.arg(attr_id)::bigint
      AND (pa.value_number = sqlc.narg(attr_number)::numeric
           OR pa.value_text = sqlc.narg(attr_text)::text)
  ))
ORDER BY p.min_price ASC, p.id ASC
LIMIT sqlc.arg(page_size)::int;

-- name: ListPublicProductsByPriceDesc :many
SELECT p.id, p.category_id, p.name, p.slug, p.thumbnail, p.is_active,
       COALESCE(p.min_price, 0)::bigint AS min_price
FROM catalog.products p
JOIN catalog.categories c
  ON c.id = p.category_id AND c.is_active AND c.deleted_at IS NULL
WHERE p.deleted_at IS NULL
  AND p.is_active
  AND p.min_price IS NOT NULL
  AND (sqlc.arg(after_id)::bigint = 0
       OR (p.min_price, p.id) < (sqlc.arg(after_price)::bigint, sqlc.arg(after_id)::bigint))
  AND (sqlc.arg(category_id)::bigint = 0 OR p.category_id = sqlc.arg(category_id)::bigint)
  AND (sqlc.arg(search)::text = '' OR p.name ILIKE '%' || sqlc.arg(search)::text || '%' ESCAPE '\')
  AND (sqlc.arg(price_min)::bigint = 0 OR p.min_price >= sqlc.arg(price_min)::bigint)
  AND (sqlc.arg(price_max)::bigint = 0 OR p.min_price <= sqlc.arg(price_max)::bigint)
  AND (sqlc.arg(attr_id)::bigint = 0 OR EXISTS (
    SELECT 1 FROM catalog.product_attributes pa
    WHERE pa.product_id = p.id
      AND pa.attribute_id = sqlc.arg(attr_id)::bigint
      AND (pa.value_number = sqlc.narg(attr_number)::numeric
           OR pa.value_text = sqlc.narg(attr_text)::text)
  ))
ORDER BY p.min_price DESC, p.id DESC
LIMIT sqlc.arg(page_size)::int;

-- name: ListAdminProducts :many
SELECT p.id, p.category_id, p.name, p.slug, p.thumbnail, p.is_active,
       COALESCE((
         SELECT MIN(pv.price) FROM catalog.product_variants pv
         WHERE pv.product_id = p.id AND pv.is_active
       ), 0)::bigint AS min_price
FROM catalog.products p
WHERE p.deleted_at IS NULL
  AND (sqlc.arg(after_id)::bigint = 0 OR p.id < sqlc.arg(after_id)::bigint)
  AND (sqlc.arg(category_id)::bigint = 0 OR p.category_id = sqlc.arg(category_id)::bigint)
  AND (sqlc.arg(search)::text = '' OR p.name ILIKE '%' || sqlc.arg(search)::text || '%' ESCAPE '\')
ORDER BY p.id DESC
LIMIT sqlc.arg(page_size)::int;

-- name: RecalcProductMinPrice :exec
UPDATE catalog.products p
SET min_price = (
  SELECT MIN(pv.price) FROM catalog.product_variants pv
  WHERE pv.product_id = p.id AND pv.is_active
)
WHERE p.id = $1;