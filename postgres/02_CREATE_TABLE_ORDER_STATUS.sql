CREATE TABLE order_schema.order_status (
  order_status_id SERIAL PRIMARY KEY,
  order_status_desc TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO order_schema.order_status (order_status_desc)
VALUES
  ('Created'),
  ('Processing'),
  ('Shipped'),
  ('Delivered'),
  ('Cancelled');
