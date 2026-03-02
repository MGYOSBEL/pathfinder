package topicparser

import "errors"

var (
	ErrMissingTopic            = errors.New("mqtt_topic metadata is missing")
	ErrNoMatchingPattern       = errors.New("topic doesn't match any config pattern")
	ErrInvalidTopicStructure   = errors.New("topic matches a pattern, but can't be parsed by that config's rule")
	ErrPayloadExtractionFailed = errors.New("json payload extraction failed")
)
