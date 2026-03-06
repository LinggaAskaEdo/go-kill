-- +goose Up
INSERT INTO public.categories (name, slug, parent_id) VALUES
    ('Electronics', 'electronics', NULL),
    ('Computers & Laptops', 'computers-laptops', NULL),
    ('Computer Accessories', 'computer-accessories', NULL),
    ('Clothing & Fashion', 'clothing-fashion', NULL),
    ('Home & Garden', 'home-garden', NULL),
    ('Sports & Outdoors', 'sports-outdoors', NULL),
    ('Books & Media', 'books-media', NULL),
    ('Health & Beauty', 'health-beauty', NULL),
    ('Toys & Games', 'toys-games', NULL),
    ('Office Supplies', 'office-supplies', NULL);

-- +goose Down
DELETE FROM public.categories WHERE slug IN (
    'electronics',
    'computers-laptops',
    'computer-accessories',
    'clothing-fashion',
    'home-garden',
    'sports-outdoors',
    'books-media',
    'health-beauty',
    'toys-games',
    'office-supplies'
);