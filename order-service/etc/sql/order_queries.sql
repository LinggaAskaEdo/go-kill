-- name: CreateOrder
INSERT INTO orders (user_id, order_number, status, total_amount, shipping_address_id, billing_address_id, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
RETURN id;

-- name: CreateOrderItemsNamed
INSERT INTO order_items (order_id, product_id, product_name, quantity, unit_price, subtotal, created_at) 
VALUES (:order_id, :product_id, :product_name, :quantity, :unit_price, :subtotal, NOW());

-- name: CreatePayment
INSERT INTO payments (order_id, payment_method, amount, status, created_at, updated_at)
VALUES (?, ?, ?, 'pending', NOW(), NOW());

-- name: CreateStatusHistory
INSERT INTO order_status_history (order_id, status, note, created_at)
VALUES (?, ?, ?, NOW());

-- name: GetOrder
SELECT id, order_number, status, total_amount 
FROM orders 
WHERE id = ? AND user_id = ?;

-- name: GetOrderItem
SELECT id, product_id, product_name, quantity, unit_price, subtotal 
FROM order_items 
WHERE order_id = ?;

-- name: GetOrderLimit
SELECT id, order_number, status, total_amount
FROM orders
WHERE user_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: GetOrderTotal
SELECT COUNT(*) 
FROM orders 
WHERE user_id = ?;

-- name: UpdateOrderStatus
UPDATE orders 
SET status = 'cancelled', updated_at = NOW() 
WHERE id = ?;