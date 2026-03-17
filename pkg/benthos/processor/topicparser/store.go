package topicparser

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MGYOSBEL/pathfinder/pkg/benthos/processor/topicparser/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConfigStore struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

func NewConfigStore(ctx context.Context, dsn string) (*ConfigStore, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &ConfigStore{
		queries: db.New(pool),
		pool:    pool,
	}, nil
}

// LoadConfigs loads all enabled configs from the database
func (cs *ConfigStore) LoadConfigs(ctx context.Context) ([]TopicParserConfig, error) {
	rows, err := cs.queries.ListEnabledConfigs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query configs: %w", err)
	}

	var configs []TopicParserConfig
	for _, row := range rows {
		config, err := rowToConfig(row)
		if err != nil {
			return nil, fmt.Errorf("failed to convert row: %w", err)
		}
		configs = append(configs, config)
	}

	return configs, nil
}

// rowToConfig converts a DB row to TopicParserConfig
func rowToConfig(row db.ListEnabledConfigsRow) (TopicParserConfig, error) {
	var metadataConfig []MetadataEntry
	var payloadConfig PayloadConfig

	// Unmarshal metadata_config JSONB
	if len(row.MetadataConfig) > 0 {
		if err := json.Unmarshal(row.MetadataConfig, &metadataConfig); err != nil {
			return TopicParserConfig{}, fmt.Errorf("failed to unmarshal metadata_config: %w", err)
		}
	}

	// Unmarshal payload_config JSONB
	if len(row.PayloadConfig) > 0 {
		if err := json.Unmarshal(row.PayloadConfig, &payloadConfig); err != nil {
			return TopicParserConfig{}, fmt.Errorf("failed to unmarshal payload_config: %w", err)
		}
	}

	return TopicParserConfig{
		ID:             int(row.ID),
		Name:           row.Name,
		Pattern:        row.Pattern,
		Version:        row.Version,
		Enabled:        row.Enabled,
		Priority:       int(row.Priority),
		MetadataConfig: metadataConfig,
		PayloadConfig:  payloadConfig,
	}, nil
}

func (cs *ConfigStore) Close() {
	cs.pool.Close()
}
