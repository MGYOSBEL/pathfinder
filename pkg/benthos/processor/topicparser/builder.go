package topicparser

import (
	"fmt"
	"time"
)

// BuildMessage assembles a processed message from extracted metadata and payload
func BuildMessage(extractedMetadata, extractedPayload map[string]string, configVersion string) (*ProcessedMessage, error) {
	if configVersion == "" {
		return nil, fmt.Errorf("%w: config version cannot be empty", ErrInvalidTopicStructure)
	}

	// Check for key overlaps between metadata and payload
	if extractedMetadata != nil && extractedPayload != nil {
		for key := range extractedPayload {
			if _, exists := extractedMetadata[key]; exists {
				return nil, fmt.Errorf("%w: key %q exists in both metadata and payload", ErrPayloadExtractionFailed, key)
			}
		}
	}

	// Create a copy of the maps to avoid external modifications
	metadata := make(map[string]string)
	if extractedMetadata != nil {
		for k, v := range extractedMetadata {
			metadata[k] = v
		}
	}

	payload := make(map[string]string)
	if extractedPayload != nil {
		for k, v := range extractedPayload {
			payload[k] = v
		}
	}

	return &ProcessedMessage{
		ExtractedMetadata: metadata,
		ExtractedPayload:  payload,
		ConfigVersion:     configVersion,
		Timestamp:         time.Now().UTC().Format(time.RFC3339),
	}, nil
}

// BuildDLQMessage creates a dead-letter-queue message for failed processing
func BuildDLQMessage(topic string, payload []byte, errorMsg, configVersion string) (*DLQMessage, error) {
	// Copy the payload to avoid external modifications
	payloadCopy := make([]byte, len(payload))
	copy(payloadCopy, payload)

	return &DLQMessage{
		OriginalTopic:   topic,
		OriginalPayload: payloadCopy,
		ErrorMessage:    errorMsg,
		ConfigVersion:   configVersion,
		Timestamp:       time.Now().UTC().Format(time.RFC3339),
	}, nil
}
