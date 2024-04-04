USE order_analytics;
CREATE MATERIALIZED VIEW order_status_mv TO order_status_history AS
SELECT
  order_status_id,
  order_status_desc,
  parseDateTimeBestEffort(created_at) AS created_at,
  __op,
  __table,
  __lsn,
  parseDateTimeBestEffort(__source_ts_ms) AS __source_ts_ms
FROM order_status_landing;
