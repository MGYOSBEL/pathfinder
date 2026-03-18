package topicparser

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redpanda-data/benthos/v4/public/service"
)

type topicParserProcessor struct {
	configs []TopicParserConfig
	store   *ConfigStore
}

func (tp *topicParserProcessor) Process(ctx context.Context, msg *service.Message) (service.MessageBatch, error) {
	topic, found := msg.MetaGet("mqtt_topic")
	if !found {
		return nil, fmt.Errorf("missing mqtt_topic in message metadata")
	}

	// Match configs by pattern
	matches := Match(topic, tp.configs)
	if len(matches) == 0 {
		// No matching config, drop message
		return nil, nil
	}

	// Use first match (highest priority)
	config := matches[0]

	// Parse topic metadata
	meta, err := ParseTopic(topic, config.Pattern, config.MetadataConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse topic: %w", err)
	}

	// Extract payload
	payload, err := msg.AsBytes()
	if err != nil {
		return nil, fmt.Errorf("failed to get payload: %w", err)
	}

	extracted, err := ExtractPayload(payload, config.PayloadConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to extract payload: %w", err)
	}

	// Build final message
	processed, err := BuildMessage(meta.Fields, extracted, config.Version)
	if err != nil {
		return nil, fmt.Errorf("failed to build message: %w", err)
	}

	// Marshal to JSON and return
	jsonBytes, err := json.Marshal(processed)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	outMsg := service.NewMessage(jsonBytes)
	_ = msg.MetaWalk(func(key string, value string) error {
		outMsg.MetaSet(key, value)
		return nil
	})

	return []*service.Message{outMsg}, nil
}

func (tp *topicParserProcessor) Close(ctx context.Context) error {
	if tp.store != nil {
		tp.store.Close()
	}
	return nil
}

func configSpec() *service.ConfigSpec {
	return service.NewConfigSpec().
		Summary("Topic Parser processor - extracts metadata and payload from MQTT topics").
		Field(service.NewStringField("dsn").Description("PostgreSQL connection DSN"))
}

func newProcessor(conf *service.ParsedConfig, mgr *service.Resources) (service.Processor, error) {
	dsn, err := conf.FieldString("dsn")
	if err != nil {
		return nil, fmt.Errorf("failed to read dsn field: %w", err)
	}

	store, err := NewConfigStore(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create config store: %w", err)
	}

	configs, err := store.LoadConfigs(context.Background())
	if err != nil {
		store.Close()
		return nil, fmt.Errorf("failed to load configs: %w", err)
	}

	return &topicParserProcessor{
		configs: configs,
		store:   store,
	}, nil
}

func init() {
	err := service.RegisterProcessor(
		"topic_parser",
		configSpec(),
		newProcessor,
	)
	if err != nil {
		panic(err)
	}
}
