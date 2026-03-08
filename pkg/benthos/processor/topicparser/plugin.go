package topicparser

import (
	"context"
	"fmt"

	"github.com/redpanda-data/benthos/v4/public/service"
)

type topicParserProcessor struct {
	config TopicParserConfig
}

func (tp *topicParserProcessor) Process(ctx context.Context, msg *service.Message) (service.MessageBatch, error) {
	topic, found := msg.MetaGet("mqtt_topic")
	if !found {
		return nil, fmt.Errorf("missing mqtt_topic in message metadata")
	}
	meta, err := ParseTopic(topic, tp.config.Pattern, tp.config.MetadataConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse topic: %w", err)
	}
	fmt.Println(meta)
	return []*service.Message{msg}, nil
}

func (tp *topicParserProcessor) Close(ctx context.Context) error {
	return nil
}

func configSpec() *service.ConfigSpec {
	return service.NewConfigSpec().
		Summary("Topic Parser processor - extracts metadata and payload from MQTT topics")
}

func newProcessor(conf *service.ParsedConfig, mgr *service.Resources) (service.Processor, error) {
	config := TopicParserConfig{
		ID:      1,
		Name:    "default_parser",
		Pattern: "#", // This guarantees that all topics match this config
		Version: "1.0",
		MetadataConfig: []MetadataEntry{
			{
				TagName: "plant",
				Type:    MetadataTypeConstant,
				Value:   "Celsa",
			},
			{
				TagName: "site",
				Type:    MetadataTypeConstant,
				Value:   "Barcelona",
			},
			{
				TagName: "plant",
				Type:    MetadataTypeTopicSegment,
				Value:   "0",
			},
			{
				TagName: "line",
				Type:    MetadataTypeTopicSegment,
				Value:   "1",
			},
			{
				TagName: "machine",
				Type:    MetadataTypeTopicSegment,
				Value:   "2:",
			},
		},
		Enabled:  true,
		Priority: 0,
	}
	return &topicParserProcessor{config: config}, nil
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
