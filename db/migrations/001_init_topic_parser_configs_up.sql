-- Create the pathfinder_config table
CREATE TABLE IF NOT EXISTS pathfinder_config (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    pattern TEXT NOT NULL,
    version TEXT NOT NULL default '1.0',
    enabled BOOLEAN NOT NULL default true,
    priority INTEGER NOT NULL default 0,
    metadata_config JSONB default '[]',
    payload_config JSONB default '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create index on pattern and enabled for faster lookups
CREATE INDEX IF NOT EXISTS idx_topic_parser_configs_enabled_priority 
ON pathfinder_config (enabled DESC, priority DESC);
