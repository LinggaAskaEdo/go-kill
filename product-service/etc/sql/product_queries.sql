-- name: CreateProduct
INSERT INTO products (name, description, price, sku, is_active, created_at, updated_at)
VALUES($1, $2, $3, $4, $5, NOW(), NOW())
RETURNING id;

-- name: CreateProductCategories
INSERT INTO product_categories (product_id, category_id, created_at)
SELECT $1, unnest($2::text[])::uuid, NOW();

-- name: CreateProductInventory
INSERT INTO inventory (product_id, quantity, reserved_quantity, updated_at)
VALUES($1, $2, $3, NOW());

-- name: GetListProducts
SELECT id, name, description, price, sku, is_active 
FROM products 
WHERE is_active = true 
LIMIT 50;

-- name: GetProductByID
SELECT id, name, description, price, sku, is_active 
FROM products 
WHERE id = $1 AND is_active = true;

-- name: GetListCategories
SELECT id, name, slug 
FROM categories;

-- name: GetCategoriesByProductID
SELECT c.id, c.name, c.slug 
FROM categories c
INNER JOIN product_categories pc ON c.id = pc.category_id
WHERE pc.product_id = $1;

-- name: GetProductsByCategoryID
SELECT p.id, p.name, p.description, p.price, p.sku 
FROM products p
INNER JOIN product_categories pc ON p.id = pc.product_id
WHERE pc.category_id = $1 AND p.is_active = true;