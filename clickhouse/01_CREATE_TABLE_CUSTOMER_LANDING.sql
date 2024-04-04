USE order_analytics;
CREATE TABLE customer_landing
(
  customer_id UInt64,
  customer_name String,
  created_at String,
  __op String,
  __table String,
  __lsn UInt64,
  __source_ts_ms String 
)
  ENGINE = Kafka('kafka:29092', 'order_db.order_schema.customer', 'clickhouse',
        'JSONEachRow') settings kafka_thread_per_consumer = 0, kafka_num_consumers = 1;
