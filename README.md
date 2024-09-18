# Ecommerce API

## Overview

The `ecommerce-api` is a RESTful API built with Go and Gin that provides endpoints for managing an e-commerce system. It supports CRUD operations for categories, products, variants, and orders, with JWT authentication for secure access.

## Features

- **Health Check**: Endpoint to check the server status.
- **Categories**: Create, update, delete, and retrieve categories.
- **Products**: Create, update, delete, and retrieve products.
- **Variants**: Create, update, delete, and retrieve variants.
- **Orders**: Create and retrieve orders.

## Swagger Documentation

The API is documented using Swagger. You can access the Swagger documentation to view and interact with the API endpoints.

- **Swagger Documentation URL**: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

The Swagger documentation provides a user-friendly interface to explore the API, view request and response examples, and test endpoints interactively.

## Getting Started

### Set up your .env file:

 * Create a .env file in the root directory of the project and add the following environment variables:

- ECOMMERCE_DEBUG=true
- ECOMMERCE_PORT=8080
- ECOMMERCE_DB_PORT=5432
- ECOMMERCE_DB_USER=user
- ECOMMERCE_DB_PASSWORD=password
- ECOMMERCE_DB_HOST=localhost
- ECOMMERCE_DB_DATABASE=ecommerce
- ECOMMERCE_DB_SCHEMA=public
- ECOMMERCE_ACCEPTED_VERSIONS=v1,v2,v3a
- ECOMMERCE_JWT_SECRET_KEY=your_jwt_secret_key

## Running the Application

- To run the server, use the following command:

    `go run main.go --runserver`
This will start the server and make the API available at http://localhost:8080.

## Running Migrations

To manage database migrations, you can use the following commands:

- Apply Migrations Up

    `go run main.go -migration -up`
This command will apply any new migrations to the database.

- Apply Migrations Down
    `go run main.go -migration -down`
This command will revert the most recent migrations.

## API Endpoints

### Health Check

- GET /:version/health

### Categories

- POST /:version/category - Create a category
- PATCH /:version/category/:id - Update a category
- DELETE /:version/category/:id - Delete a category
- GET /:version/category/:id - Get category by ID
- GET /:version/category - Get all categories

### Products

- POST /:version/products - Create a product
- PATCH /:version/products/:id - Update a product
- DELETE /:version/products/:id - Delete a product
- GET /:version/products/:id - Get product by ID
- GET /:version/products - Get all products

### Variants

- POST /:version/variants - Create a variant
- PATCH /:version/variants/:id - Update a variant
- DELETE /:version/variants/:id - Delete a variant
- GET /:version/variants/:id - Get variant by ID
- GET /:version/variants - Get all variants

### Orders

- POST /:version/orders - Create an order
- GET /:version/orders/:id - Get order by ID
- GET /:version/orders - Get all orders
