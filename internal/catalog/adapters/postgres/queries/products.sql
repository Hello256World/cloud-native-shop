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