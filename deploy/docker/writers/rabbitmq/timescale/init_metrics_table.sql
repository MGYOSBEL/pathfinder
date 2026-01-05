CREATE TABLE IF NOT EXISTS metrics_data (
    topic TEXT,
    plant text,
    line text,
    cell text,    
    measurement text,
    variable text,
    value DOUBLE precision,
    unit text,
    timestamp TIMESTAMPTZ
);

SELECT create_hypertable('metrics_data', 'timestamp');
