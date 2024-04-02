CREATE TABLE orders_schema.order_statuses (
  order_status_id SERIAL PRIMARY KEY,
  order_status_desc TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO orders_schema.order_statuses (order_status_desc)
VALUES
  ('Created'),
  ('Processing'),
  ('Shipped'),
  ('Delivered'),
  ('Cancelled');
