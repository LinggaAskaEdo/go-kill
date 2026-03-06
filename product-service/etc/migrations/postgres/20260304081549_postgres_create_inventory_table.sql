-- +goose Up
CREATE TABLE inventory (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    product_id UUID UNIQUE NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    reserved_quantity INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS inventory;
