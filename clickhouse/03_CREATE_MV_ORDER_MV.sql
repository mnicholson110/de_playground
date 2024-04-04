USE order_analytics;
CREATE MATERIALIZED VIEW order_mv TO order_history AS
SELECT
order_id,
  order_amount,
  order_status_id,
  customer_id,
  parseDateTimeBestEffort(created_at) AS created_at,
  parseDateTimeBestEffort(updated_at) AS updated_at,
  __op,
  __table,
  __lsn,
  parseDateTimeBestEffort(__source_ts_ms) AS __source_ts_ms
FROM order_landing;
