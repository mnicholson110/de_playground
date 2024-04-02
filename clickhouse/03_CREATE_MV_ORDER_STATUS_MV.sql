USE order_analytics;
CREATE MATERIALIZED VIEW order_status_mv TO order_status_history AS
SELECT *
FROM order_status_landing;
