USE order_analytics;
CREATE MATERIALIZED VIEW customer_mv TO customer_history AS
SELECT *
FROM customer_landing;
