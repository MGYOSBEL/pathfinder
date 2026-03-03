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
	fmt.Printf("receive message from topic %s\n", topic)
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
		ID:       1,
		Name:     "default_parser",
		Pattern:  "+/+/+",
		Version:  "1.0",
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
