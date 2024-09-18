-- Drop OrderItems Table
DROP TABLE IF EXISTS order_items;

-- Drop Orders Table
DROP TABLE IF EXISTS orders;

-- Drop Product-Variant Mapping Table
DROP TABLE IF EXISTS product_variant_mapping;

-- Drop Variant Table
DROP TABLE IF EXISTS variants;

-- Drop Product-Category Mapping Table
DROP TABLE IF EXISTS product_category_mapping;

-- Drop Product Table
DROP TABLE IF EXISTS products;

-- Drop Category Mapping Table
DROP TABLE IF EXISTS category_mapping;

-- Drop Category Table
DROP TABLE IF EXISTS categories;

-- Drop the UUID extension
DROP EXTENSION IF EXISTS "uuid-ossp";
