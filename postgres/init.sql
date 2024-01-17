CREATE SCHEMA order_schema
  CREATE TABLE order_statuses (
    orderStatusId INTEGER PRIMARY KEY,
    orderStatus TEXT NOT NULL
  )
  CREATE TABLE orders (
    orderId INTEGER PRIMARY KEY,
    orderStatusId INTEGER NOT NULL,
    orderStatus TEXT NOT NULL
  );
