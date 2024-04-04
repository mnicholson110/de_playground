USE order_analytics;
CREATE TABLE order_status_landing (
    order_status_id UInt64,
    order_status_desc String, 
    created_at String,
    __op String,
    __table String,
    __lsn UInt64,
    __source_ts_ms String 
)
  ENGINE = Kafka('kafka:29092', 'order_db.order_schema.order_status', 'clickhouse',
        'JSONEachRow') settings kafka_thread_per_consumer = 0, kafka_num_consumers = 1;
