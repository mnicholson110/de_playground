USE order_analytics;
CREATE TABLE customer_history
(
  customer_id UInt64,
  customer_name String,
  created_at DateTime,
  __op String,
  __table String,
  __lsn UInt64,
  __source_ts_ms DateTime 
)
  ENGINE = MergeTree ORDER BY (created_at, customer_id);
