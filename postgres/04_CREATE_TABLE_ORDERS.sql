CREATE TABLE orders_schema.orders (
  order_id SERIAL PRIMARY KEY,
  order_amount NUMERIC(10,2) NOT NULL,
  order_status_id INTEGER NOT NULL,
  customer_id INTEGER NOT NULL,
  FOREIGN KEY (order_status_id) REFERENCES orders_schema.order_statuses(order_status_id),
  FOREIGN KEY (customer_id) REFERENCES orders_schema.customers(customer_id),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
