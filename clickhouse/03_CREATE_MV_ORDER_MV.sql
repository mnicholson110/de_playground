USE order_analytics;
CREATE MATERIALIZED VIEW order_mv TO order_history AS
SELECT *
FROM order_landing;
