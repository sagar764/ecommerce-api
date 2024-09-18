-- Enable extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Define ENUM type for order status
CREATE TYPE order_status AS ENUM (
    'Pending',       -- Order has been placed but not yet processed
    'Accepted',      -- Order has been accepted and is being processed
    'Shipped',       -- Order has been shipped to the customer
    'Delivered',     -- Order has been delivered to the customer
    'Cancelled',     -- Order has been cancelled
    'Returned'       -- Order has been returned by the customer
);

-- Create Category Table with UUID
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID,
    updated_by UUID
);

-- Create Product Table with UUID
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    image_url VARCHAR(255),
    is_active BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID,
    updated_by UUID
);

-- Create Variant Table with UUID
CREATE TABLE variants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255),
    mrp DECIMAL(10, 2) NOT NULL,
    discount_price DECIMAL(10, 2),
    size VARCHAR(50),
    color VARCHAR(50),
    quantity INT NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID,
    updated_by UUID
);

-- Create Category Mapping Table (for hierarchical child-parent category relationships)
CREATE TABLE category_mapping (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    parent_category_id UUID REFERENCES categories(id) ON DELETE CASCADE,
    child_category_id UUID REFERENCES categories(id) ON DELETE CASCADE
);

-- Create Product-Category Mapping Table (for many-to-many relationship between products and categories)
CREATE TABLE product_category_mapping (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID REFERENCES products(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id) ON DELETE CASCADE
);

-- Create Product-Variant Mapping Table
CREATE TABLE product_variant_mapping (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID REFERENCES products(id) ON DELETE CASCADE,
    variant_id UUID REFERENCES variants(id) ON DELETE CASCADE
);

-- Create Order Table with ENUM for status
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    status order_status DEFAULT 'Accepted', 
    order_total DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID,
    updated_by UUID
);

-- Create OrderItems Table to track ordered variants
CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID REFERENCES orders(id) ON DELETE CASCADE,
    variant_id UUID REFERENCES variants(id),
    quantity INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL
);
