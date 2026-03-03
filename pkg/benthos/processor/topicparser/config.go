// Package topicparser implements configuration structures for the MQTT topic
// parser Benthos processor plugin. It provides types for representing topic
// parsing rules, metadata extraction, and payload field mapping based on
// ISA-95 industrial hierarchy standards.
package topicparser

import "fmt"

// MetadataType defines valid types for metadata extraction
type MetadataType string

const (
	MetadataTypeConstant     MetadataType = "Constant"
	MetadataTypeTopicSegment MetadataType = "TopicSegment"
)

// IsValid checks if the metadata type is valid
func (mt MetadataType) IsValid() bool {
	return mt == MetadataTypeConstant || mt == MetadataTypeTopicSegment
}

type MetadataEntry struct {
	TagName string       `json:"tag_name"`
	Type    MetadataType `json:"type,omitempty"`
	Value   string       `json:"value,omitempty"`
}

// Validate checks if the metadata entry is valid
func (me MetadataEntry) Validate() error {
	if me.TagName == "" {
		return fmt.Errorf("tag_name cannot be empty")
	}

	if !me.Type.IsValid() {
		return fmt.Errorf("invalid type %q, must be %q or %q", me.Type, MetadataTypeConstant, MetadataTypeTopicSegment)
	}

	if me.Value == "" {
		return fmt.Errorf("value cannot be empty")
	}

	return nil
}

type PayloadConfig struct {
	Variable string `json:"variable"`
	Unit     string `json:"unit"`
	Value    string `json:"value"`
}

type ExtractedMetadata struct {
	Fields map[string]string
}

type TopicParserConfig struct {
	ID             int
	Name           string          `json:"name"`
	Pattern        string          `json:"pattern"`
	Version        string          `json:"version"`
	Enabled        bool            `json:"enabled"`
	Priority       int             `json:"priority"`
	MetadataConfig []MetadataEntry `json:"metadata_config"`
	PayloadConfig  PayloadConfig   `json:"payload_config"`
}
