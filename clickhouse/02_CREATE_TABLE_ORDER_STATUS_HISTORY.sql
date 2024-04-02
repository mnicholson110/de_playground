USE order_analytics;
CREATE TABLE order_status_history
(
  order_status_id UInt64,
  order_status_desc String,
  created_at DateTime,
  __op String,
  __table String,
  __lsn UInt64,
  __source_ts_ms DateTime
)
  ENGINE = MergeTree ORDER BY (created_at, order_status_id);
