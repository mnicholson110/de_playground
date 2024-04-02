CREATE TABLE order_schema.order (
  order_id SERIAL PRIMARY KEY,
  order_amount NUMERIC(10,2) NOT NULL,
  order_status_id INTEGER NOT NULL,
  customer_id INTEGER NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (order_status_id) REFERENCES order_schema.order_status(order_status_id),
  FOREIGN KEY (customer_id) REFERENCES order_schema.customer(customer_id)
);
