// Package topicparser implements configuration structures for the MQTT topic
// parser Benthos processor plugin. It provides types for representing topic
// parsing rules, metadata extraction, and payload field mapping based on
// ISA-95 industrial hierarchy standards.
package topicparser

type MetadataEntry struct {
	TagName string `json:"tag_name"`
	Type    string `json:"type,omitempty"`
	Value   string `json:"value,omitempty"`
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
