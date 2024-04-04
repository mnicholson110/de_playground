USE order_analytics;
CREATE MATERIALIZED VIEW customer_mv TO customer_history AS
SELECT
  customer_id,
  customer_name,
  parseDateTimeBestEffort(created_at) AS created_at,
  __op,
  __table,
  __lsn,
  parseDateTimeBestEffort(__source_ts_ms) AS __source_ts_ms
FROM customer_landing;
