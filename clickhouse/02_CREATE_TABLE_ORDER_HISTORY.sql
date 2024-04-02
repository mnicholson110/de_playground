USE order_analytics;
CREATE TABLE order_history
(
	order_id UInt64,
	order_amount String,
	order_status_id UInt64,
	customer_id UInt64,
	created_at DateTime,
	updated_at DateTime,
	__op String,
	__table String,
	__lsn UInt64,
	__source_ts_ms DateTime
)
	ENGINE = MergeTree ORDER BY (updated_at, order_status_id);
