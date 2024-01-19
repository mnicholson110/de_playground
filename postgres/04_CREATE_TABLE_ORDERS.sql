CREATE TABLE order_schema.orders (
  order_id INTEGER PRIMARY KEY,
  order_date DATE NOT NULL,
  order_status_id INTEGER NOT NULL,
  customer_id INTEGER NOT NULL,
  FOREIGN KEY (order_status_id) REFERENCES order_schema.order_statuses(order_status_id),
  FOREIGN KEY (customer_id) REFERENCES order_schema.customers(customer_id)
);
