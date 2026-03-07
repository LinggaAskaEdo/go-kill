-- +goose Up
CREATE TABLE order_status_history (
    id UUID PRIMARY KEY DEFAULT uuid_v7(),
    order_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL,
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    INDEX idx_order_id (order_id)
) ENGINE=InnoDB;

-- +goose Down
DROP TABLE IF EXISTS order_status_history;
