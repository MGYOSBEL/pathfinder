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

-- Insert test seed data matching the hardcoded config from plugin.go
INSERT INTO pathfinder_config (name, pattern, version, enabled, priority, metadata_config, payload_config)
VALUES (
    'default_parser',
    '#',
    '1.0',
    true,
    0,
    '[
        {"tag_name": "plant", "type": "Constant", "value": "Celsa"},
        {"tag_name": "site", "type": "Constant", "value": "Barcelona"},
        {"tag_name": "plant", "type": "TopicSegment", "value": "0"},
        {"tag_name": "line", "type": "TopicSegment", "value": "1"},
        {"tag_name": "machine", "type": "TopicSegment", "value": "2:"}
    ]'::jsonb,
    '{"variable": "$.name", "value": "$.value", "unit": "$.unit"}'::jsonb
) ON CONFLICT (name) DO NOTHING;
