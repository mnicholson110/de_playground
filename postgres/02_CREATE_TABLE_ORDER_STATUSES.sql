CREATE TABLE order_schema.order_statuses (
  order_status_id SERIAL PRIMARY KEY,
  order_status_desc TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO order_schema.order_statuses (order_status_desc)
VALUES
  ('Created'),
  ('Processing'),
  ('Shipped'),
  ('Delivered'),
  ('Cancelled');
