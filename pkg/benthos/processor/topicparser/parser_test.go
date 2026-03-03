package topicparser_test

import (
	"testing"

	"github.com/MGYOSBEL/pathfinder/pkg/benthos/processor/topicparser"
	"github.com/stretchr/testify/assert"
)

var (
	enterprise string = "stark-industries"
	site       string = "barcelona"
	plant      string = "ironman-manufacuture"
	machine    string = "engine-A"
)

func constantValue(s string) *topicparser.ConstantValue {
	return &topicparser.ConstantValue{Value: s}
}

func topicSegmentValue(s string) *topicparser.TopicSegmentValue {
	return &topicparser.TopicSegmentValue{Value: s}
}

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
			name:     "single constant",
			topic:    "foo/bar/baz",
			pattern:  "foo/bar/baz",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName:  "enterprise",
					Constant: constantValue(enterprise),
				},
			},
			wantFields: map[string]string{
				"enterprise": enterprise,
			},
			wantErr: nil,
		},
		{
			name:     "multiple constants",
			topic:    "foo/bar/baz",
			pattern:  "foo/bar/baz",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName:  "enterprise",
					Constant: constantValue(enterprise),
				},
				{
					TagName:  "plant",
					Constant: constantValue(plant),
				},
				{
					TagName:  "site",
					Constant: constantValue(site),
				},
				{
					TagName:  "machine",
					Constant: constantValue(machine),
				},
			},
			wantFields: map[string]string{
				"enterprise": enterprise,
				"plant":      plant,
				"site":       site,
				"machine":    machine,
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
			name:     "extract segment 0",
			topic:    "plant-a/line-1/cell-x",
			pattern:  "plant-a/+/+",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "plant",
					TopicSegment: topicSegmentValue("0"),
				},
			},
			wantFields: map[string]string{
				"plant": "plant-a",
			},
			wantErr: nil,
		},
		{
			name:     "extract segment 1",
			topic:    "plant-a/line-1/cell-x",
			pattern:  "+/+/+",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "line",
					TopicSegment: topicSegmentValue("1"),
				},
			},
			wantFields: map[string]string{
				"line": "line-1",
			},
			wantErr: nil,
		},
		{
			name:     "extract segment 2",
			topic:    "plant-a/line-1/cell-x",
			pattern:  "+/+/+",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "cell",
					TopicSegment: topicSegmentValue("2"),
				},
			},
			wantFields: map[string]string{
				"cell": "cell-x",
			},
			wantErr: nil,
		},
		{
			name:     "index out of bounds",
			topic:    "plant-a/line-1",
			pattern:  "+/+",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "cell",
					TopicSegment: topicSegmentValue("5"),
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
			name:     "extract from segment 2 to end (single remaining)",
			topic:    "plant-a/line-1/sensor-data",
			pattern:  "plant-a/+/+",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "measurement",
					TopicSegment: topicSegmentValue("2:"),
				},
			},
			wantFields: map[string]string{
				"measurement": "sensor-data",
			},
			wantErr: nil,
		},
		{
			name:     "extract from segment 2 to end (multiple remaining)",
			topic:    "plant-a/line-1/sensor/temperature/reading",
			pattern:  "plant-a/+/#",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "measurement",
					TopicSegment: topicSegmentValue("2:"),
				},
			},
			wantFields: map[string]string{
				"measurement": "sensor/temperature/reading",
			},
			wantErr: nil,
		},
		{
			name:     "extract from segment 0 to end",
			topic:    "a/b/c/d",
			pattern:  "+/+/+/+",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "all",
					TopicSegment: topicSegmentValue("0:"),
				},
			},
			wantFields: map[string]string{
				"all": "a/b/c/d",
			},
			wantErr: nil,
		},
		{
			name:     "range out of bounds",
			topic:    "a/b",
			pattern:  "+/+",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "data",
					TopicSegment: topicSegmentValue("5:"),
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
			name:     "extract indices 0,2",
			topic:    "a/b/c/d",
			pattern:  "+/+/+/+",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "mixed",
					TopicSegment: topicSegmentValue("0,2"),
				},
			},
			wantFields: map[string]string{
				"mixed": "a/c",
			},
			wantErr: nil,
		},
		{
			name:     "extract indices 0,2,3",
			topic:    "plant/line/cell/measurement",
			pattern:  "+/+/+/+",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "hierarchy",
					TopicSegment: topicSegmentValue("0,2,3"),
				},
			},
			wantFields: map[string]string{
				"hierarchy": "plant/cell/measurement",
			},
			wantErr: nil,
		},
		{
			name:     "multi-index out of bounds",
			topic:    "a/b",
			pattern:  "+/+",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "data",
					TopicSegment: topicSegmentValue("0,5"),
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
			name:     "constants and single indices",
			topic:    "plant-a/line-1/cell-x",
			pattern:  "+/+/+",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName:  "enterprise",
					Constant: constantValue(enterprise),
				},
				{
					TagName: "plant",
					TopicSegment: topicSegmentValue("0"),
				},
				{
					TagName:  "site",
					Constant: constantValue(site),
				},
				{
					TagName: "line",
					TopicSegment: topicSegmentValue("1"),
				},
			},
			wantFields: map[string]string{
				"enterprise": enterprise,
				"plant":      "plant-a",
				"site":       site,
				"line":       "line-1",
			},
			wantErr: nil,
		},
		{
			name:     "constants and range extraction",
			topic:    "plant-a/line-1/temp/sensor-1",
			pattern:  "+/+/#",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName:  "enterprise",
					Constant: constantValue(enterprise),
				},
				{
					TagName: "plant",
					TopicSegment: topicSegmentValue("0"),
				},
				{
					TagName: "measurement_path",
					TopicSegment: topicSegmentValue("2:"),
				},
			},
			wantFields: map[string]string{
				"enterprise":       enterprise,
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
			name:     "invalid syntax in TopicSegment.Value (non-numeric)",
			topic:    "a/b/c",
			pattern:  "+/+/+",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "field",
					TopicSegment: topicSegmentValue("abc"),
				},
			},
			wantErr: topicparser.ErrInvalidTopicStructure,
		},
		{
			name:     "invalid syntax in TopicSegment.Value (colon in wrong place)",
			topic:    "a/b/c",
			pattern:  "+/+/+",
			metadataConfig: []topicparser.MetadataEntry{
				{
					TagName: "field",
					TopicSegment: topicSegmentValue(":2"),
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

