package topicparser_test

import (
	"testing"

	"github.com/MGYOSBEL/pathfinder/pkg/benthos/processor/topicparser"
	"github.com/stretchr/testify/assert"
)

func TestParseTopic_ConstantExtraction(t *testing.T) {
	tests := []struct {
		name           string
		topic          string
		pattern        string
		metadataConfig []topicparser.MetadataEntry
		wantFields     map[string]string
		wantErr        error
	}{
		{
			name:           "empty config",
			topic:          "foo/bar/baz",
			pattern:        "foo/bar/baz",
			metadataConfig: []topicparser.MetadataEntry{},
			wantFields:     map[string]string{},
			wantErr:        nil,
		},
		{
			name:    "single constant",
			topic:   "foo/bar/baz",
			pattern: "foo/bar/baz",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "enterprise",
					Type:    topicparser.MetadataTypeConstant,
					Value:   "stark-industries",
				},
			},
			wantFields: map[string]string{
				"enterprise": "stark-industries",
			},
			wantErr: nil,
		},
		{
			name:    "multiple constants",
			topic:   "foo/bar/baz",
			pattern: "foo/bar/baz",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "enterprise",
					Type:    topicparser.MetadataTypeConstant,
					Value:   "stark-industries",
				},
				{
					TagName: "plant",
					Type:    topicparser.MetadataTypeConstant,
					Value:   "ironman-manufacuture",
				},
				{
					TagName: "site",
					Type:    topicparser.MetadataTypeConstant,
					Value:   "barcelona",
				},
				{
					TagName: "machine",
					Type:    topicparser.MetadataTypeConstant,
					Value:   "engine-A",
				},
			},
			wantFields: map[string]string{
				"enterprise": "stark-industries",
				"plant":      "ironman-manufacuture",
				"site":       "barcelona",
				"machine":    "engine-A",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := topicparser.ParseTopic(tt.topic, tt.pattern, tt.metadataConfig)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantFields, got.Fields)
			}
		})
	}
}

func TestParseTopic_SingleIndexExtraction(t *testing.T) {
	tests := []struct {
		name           string
		topic          string
		pattern        string
		metadataConfig []topicparser.MetadataEntry
		wantFields     map[string]string
		wantErr        error
	}{
		{
			name:    "extract segment 0",
			topic:   "plant-a/line-1/cell-x",
			pattern: "plant-a/+/+",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "plant",
					Type:    topicparser.MetadataTypeTopicSegment,
					Value:   "0",
				},
			},
			wantFields: map[string]string{
				"plant": "plant-a",
			},
			wantErr: nil,
		},
		{
			name:    "extract segment 1",
			topic:   "plant-a/line-1/cell-x",
			pattern: "+/+/+",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "line",
					Type:    topicparser.MetadataTypeTopicSegment,
					Value:   "1",
				},
			},
			wantFields: map[string]string{
				"line": "line-1",
			},
			wantErr: nil,
		},
		{
			name:    "extract segment 2",
			topic:   "plant-a/line-1/cell-x",
			pattern: "+/+/+",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "cell",
					Type:    topicparser.MetadataTypeTopicSegment,
					Value:   "2",
				},
			},
			wantFields: map[string]string{
				"cell": "cell-x",
			},
			wantErr: nil,
		},
		{
			name:    "index out of bounds",
			topic:   "plant-a/line-1",
			pattern: "+/+",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "cell",
					Type:    topicparser.MetadataTypeTopicSegment,
					Value:   "5",
				},
			},
			wantFields: nil,
			wantErr:    topicparser.ErrInvalidTopicStructure,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := topicparser.ParseTopic(tt.topic, tt.pattern, tt.metadataConfig)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantFields, got.Fields)
			}
		})
	}
}

func TestParseTopic_RangeExtraction(t *testing.T) {
	tests := []struct {
		name           string
		topic          string
		pattern        string
		metadataConfig []topicparser.MetadataEntry
		wantFields     map[string]string
		wantErr        error
	}{
		{
			name:    "extract from segment 2 to end (single remaining)",
			topic:   "plant-a/line-1/sensor-data",
			pattern: "plant-a/+/+",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "measurement",
					Type:    topicparser.MetadataTypeTopicSegment,
					Value:   "2:",
				},
			},
			wantFields: map[string]string{
				"measurement": "sensor-data",
			},
			wantErr: nil,
		},
		{
			name:    "extract from segment 2 to end (multiple remaining)",
			topic:   "plant-a/line-1/sensor/temperature/reading",
			pattern: "plant-a/+/#",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "measurement",
					Type:    topicparser.MetadataTypeTopicSegment,
					Value:   "2:",
				},
			},
			wantFields: map[string]string{
				"measurement": "sensor/temperature/reading",
			},
			wantErr: nil,
		},
		{
			name:    "extract from segment 0 to end",
			topic:   "a/b/c/d",
			pattern: "+/+/+/+",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "all",
					Type:    topicparser.MetadataTypeTopicSegment,
					Value:   "0:",
				},
			},
			wantFields: map[string]string{
				"all": "a/b/c/d",
			},
			wantErr: nil,
		},
		{
			name:    "range out of bounds",
			topic:   "a/b",
			pattern: "+/+",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "data",
					Type:    topicparser.MetadataTypeTopicSegment,
					Value:   "5:",
				},
			},
			wantFields: nil,
			wantErr:    topicparser.ErrInvalidTopicStructure,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := topicparser.ParseTopic(tt.topic, tt.pattern, tt.metadataConfig)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantFields, got.Fields)
			}
		})
	}
}

func TestParseTopic_MultiIndexExtraction(t *testing.T) {
	tests := []struct {
		name           string
		topic          string
		pattern        string
		metadataConfig []topicparser.MetadataEntry
		wantFields     map[string]string
		wantErr        error
	}{
		{
			name:    "extract indices 0,2",
			topic:   "a/b/c/d",
			pattern: "+/+/+/+",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "mixed",
					Type:    topicparser.MetadataTypeTopicSegment,
					Value:   "0,2",
				},
			},
			wantFields: map[string]string{
				"mixed": "a/c",
			},
			wantErr: nil,
		},
		{
			name:    "extract indices 0,2,3",
			topic:   "plant/line/cell/measurement",
			pattern: "+/+/+/+",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "hierarchy",
					Type:    topicparser.MetadataTypeTopicSegment,
					Value:   "0,2,3",
				},
			},
			wantFields: map[string]string{
				"hierarchy": "plant/cell/measurement",
			},
			wantErr: nil,
		},
		{
			name:    "multi-index out of bounds",
			topic:   "a/b",
			pattern: "+/+",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "data",
					Type:    topicparser.MetadataTypeTopicSegment,
					Value:   "0,5",
				},
			},
			wantFields: nil,
			wantErr:    topicparser.ErrInvalidTopicStructure,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := topicparser.ParseTopic(tt.topic, tt.pattern, tt.metadataConfig)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantFields, got.Fields)
			}
		})
	}
}

func TestParseTopic_MixedConstantAndSegment(t *testing.T) {
	tests := []struct {
		name           string
		topic          string
		pattern        string
		metadataConfig []topicparser.MetadataEntry
		wantFields     map[string]string
		wantErr        error
	}{
		{
			name:    "constants and single indices",
			topic:   "plant-a/line-1/cell-x",
			pattern: "+/+/+",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "enterprise",
					Type:    topicparser.MetadataTypeConstant,
					Value:   "stark-industries",
				},
				{
					TagName: "plant",
					Type:    topicparser.MetadataTypeTopicSegment,
					Value:   "0",
				},
				{
					TagName: "site",
					Type:    topicparser.MetadataTypeConstant,
					Value:   "barcelona",
				},
				{
					TagName: "line",
					Type:    topicparser.MetadataTypeTopicSegment,
					Value:   "1",
				},
			},
			wantFields: map[string]string{
				"enterprise": "stark-industries",
				"plant":      "plant-a",
				"site":       "barcelona",
				"line":       "line-1",
			},
			wantErr: nil,
		},
		{
			name:    "constants and range extraction",
			topic:   "plant-a/line-1/temp/sensor-1",
			pattern: "+/+/#",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "enterprise",
					Type:    topicparser.MetadataTypeConstant,
					Value:   "stark-industries",
				},
				{
					TagName: "plant",
					Type:    topicparser.MetadataTypeTopicSegment,
					Value:   "0",
				},
				{
					TagName: "measurement_path",
					Type:    topicparser.MetadataTypeTopicSegment,
					Value:   "2:",
				},
			},
			wantFields: map[string]string{
				"enterprise":       "stark-industries",
				"plant":            "plant-a",
				"measurement_path": "temp/sensor-1",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := topicparser.ParseTopic(tt.topic, tt.pattern, tt.metadataConfig)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantFields, got.Fields)
			}
		})
	}
}

func TestParseTopic_InvalidSyntax(t *testing.T) {
	tests := []struct {
		name           string
		topic          string
		pattern        string
		metadataConfig []topicparser.MetadataEntry
		wantErr        error
	}{
		{
			name:    "invalid syntax in TopicSegment.Value (non-numeric)",
			topic:   "a/b/c",
			pattern: "+/+/+",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "field",
					Type:    topicparser.MetadataTypeTopicSegment,
					Value:   "abc",
				},
			},
			wantErr: topicparser.ErrInvalidTopicStructure,
		},
		{
			name:    "invalid syntax in TopicSegment.Value (colon in wrong place)",
			topic:   "a/b/c",
			pattern: "+/+/+",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "field",
					Type:    topicparser.MetadataTypeTopicSegment,
					Value:   ":2",
				},
			},
			wantErr: topicparser.ErrInvalidTopicStructure,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := topicparser.ParseTopic(tt.topic, tt.pattern, tt.metadataConfig)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
