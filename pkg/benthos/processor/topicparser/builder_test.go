package topicparser_test

import (
	"testing"

	"github.com/MGYOSBEL/pathfinder/pkg/benthos/processor/topicparser"
	"github.com/stretchr/testify/assert"
)

func TestBuildMessage_ValidMessages(t *testing.T) {
	tests := []struct {
		name               string
		extractedMetadata  map[string]string
		extractedPayload   map[string]string
		configVersion      string
		wantMetadata       map[string]string
		wantPayload        map[string]string
		wantConfigVersion  string
		wantErr            error
	}{
		{
			name:          "simple message with both metadata and payload",
			extractedMetadata: map[string]string{
				"enterprise": "celsa",
				"plant":      "cml03",
				"line":       "ple",
				"cell":       "furnace",
			},
			extractedPayload: map[string]string{
				"variable": "temperature",
				"value":    "23.5",
				"unit":     "celsius",
			},
			configVersion: "v0.0.1",
			wantMetadata: map[string]string{
				"enterprise": "celsa",
				"plant":      "cml03",
				"line":       "ple",
				"cell":       "furnace",
			},
			wantPayload: map[string]string{
				"variable": "temperature",
				"value":    "23.5",
				"unit":     "celsius",
			},
			wantConfigVersion: "v0.0.1",
			wantErr:           nil,
		},
		{
			name: "message with empty metadata",
			extractedMetadata: map[string]string{},
			extractedPayload: map[string]string{
				"variable": "humidity",
				"value":    "65.2",
				"unit":     "percent",
			},
			configVersion: "v0.0.2",
			wantMetadata:  map[string]string{},
			wantPayload: map[string]string{
				"variable": "humidity",
				"value":    "65.2",
				"unit":     "percent",
			},
			wantConfigVersion: "v0.0.2",
			wantErr:           nil,
		},
		{
			name: "message with many metadata fields",
			extractedMetadata: map[string]string{
				"enterprise":  "celsa",
				"site":        "barcelona",
				"plant":       "cml03",
				"area":        "stamping",
				"line":        "ple",
				"cell":        "furnace",
				"equipment":   "heating_element",
				"measurement": "temperature/sensor1",
			},
			extractedPayload: map[string]string{
				"variable": "temp-01",
				"value":    "78.5",
				"unit":     "celsius",
			},
			configVersion: "v0.1.0",
			wantMetadata: map[string]string{
				"enterprise":  "celsa",
				"site":        "barcelona",
				"plant":       "cml03",
				"area":        "stamping",
				"line":        "ple",
				"cell":        "furnace",
				"equipment":   "heating_element",
				"measurement": "temperature/sensor1",
			},
			wantPayload: map[string]string{
				"variable": "temp-01",
				"value":    "78.5",
				"unit":     "celsius",
			},
			wantConfigVersion: "v0.1.0",
			wantErr:           nil,
		},
		{
			name: "message with numeric and special characters",
			extractedMetadata: map[string]string{
				"device_id": "sensor-001",
				"location":  "room/A-2",
			},
			extractedPayload: map[string]string{
				"variable": "temp#sensor",
				"value":    "23.456",
				"unit":     "°C",
			},
			configVersion: "v0.0.1",
			wantMetadata: map[string]string{
				"device_id": "sensor-001",
				"location":  "room/A-2",
			},
			wantPayload: map[string]string{
				"variable": "temp#sensor",
				"value":    "23.456",
				"unit":     "°C",
			},
			wantConfigVersion: "v0.0.1",
			wantErr:           nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := topicparser.BuildMessage(tt.extractedMetadata, tt.extractedPayload, tt.configVersion)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantMetadata, msg.ExtractedMetadata)
				assert.Equal(t, tt.wantPayload, msg.ExtractedPayload)
				assert.Equal(t, tt.wantConfigVersion, msg.ConfigVersion)
				assert.NotEmpty(t, msg.Timestamp)
			}
		})
	}
}

func TestBuildMessage_MetadataOverlap(t *testing.T) {
	tests := []struct {
		name              string
		extractedMetadata map[string]string
		extractedPayload  map[string]string
		configVersion     string
		wantErr           error
	}{
		{
			name: "no overlap between metadata and payload keys",
			extractedMetadata: map[string]string{
				"enterprise": "celsa",
				"plant":      "cml03",
			},
			extractedPayload: map[string]string{
				"variable": "temperature",
				"value":    "23.5",
				"unit":     "celsius",
			},
			configVersion: "v0.0.1",
			wantErr:       nil,
		},
		{
			name: "overlap in keys (variable also in metadata)",
			extractedMetadata: map[string]string{
				"variable": "sensor_id",
				"plant":    "cml03",
			},
			extractedPayload: map[string]string{
				"variable": "temperature",
				"value":    "23.5",
				"unit":     "celsius",
			},
			configVersion: "v0.0.1",
			wantErr:       topicparser.ErrPayloadExtractionFailed,
		},
		{
			name: "overlap on value key",
			extractedMetadata: map[string]string{
				"value": "some_value",
			},
			extractedPayload: map[string]string{
				"variable": "temperature",
				"value":    "23.5",
				"unit":     "celsius",
			},
			configVersion: "v0.0.1",
			wantErr:       topicparser.ErrPayloadExtractionFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := topicparser.BuildMessage(tt.extractedMetadata, tt.extractedPayload, tt.configVersion)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBuildMessage_InvalidInputs(t *testing.T) {
	tests := []struct {
		name              string
		extractedMetadata map[string]string
		extractedPayload  map[string]string
		configVersion     string
		wantErr           error
	}{
		{
			name:              "empty config version",
			extractedMetadata: map[string]string{"plant": "cml03"},
			extractedPayload: map[string]string{
				"variable": "temperature",
				"value":    "23.5",
				"unit":     "celsius",
			},
			configVersion: "",
			wantErr:       topicparser.ErrInvalidTopicStructure,
		},
		{
			name:              "nil metadata map",
			extractedMetadata: nil,
			extractedPayload: map[string]string{
				"variable": "temperature",
				"value":    "23.5",
				"unit":     "celsius",
			},
			configVersion: "v0.0.1",
			wantErr:       nil, // nil is acceptable, treated as empty
		},
		{
			name:              "nil payload map",
			extractedMetadata: map[string]string{"plant": "cml03"},
			extractedPayload:  nil,
			configVersion:     "v0.0.1",
			wantErr:           nil, // nil is acceptable, treated as empty
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := topicparser.BuildMessage(tt.extractedMetadata, tt.extractedPayload, tt.configVersion)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBuildMessage_TimestampGeneration(t *testing.T) {
	msg, err := topicparser.BuildMessage(
		map[string]string{"plant": "cml03"},
		map[string]string{"variable": "temp", "value": "23.5", "unit": "celsius"},
		"v0.0.1",
	)

	assert.NoError(t, err)
	assert.NotEmpty(t, msg.Timestamp)
	// Verify timestamp is in RFC3339 format (ISO 8601)
	// Should look like: 2026-03-03T08:50:00Z
	assert.Regexp(t, `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`, msg.Timestamp)
}

func TestBuildDLQMessage_ExtractionError(t *testing.T) {
	tests := []struct {
		name            string
		topic           string
		payload         []byte
		errMsg          string
		configVersion   string
		wantTopic       string
		wantPayload     string
		wantError       string
		wantConfigVersion string
		wantErr         error
	}{
		{
			name:            "pattern match error",
			topic:           "factory/line/cell",
			payload:         []byte(`{"variable": "temp", "value": 23.5, "unit": "celsius"}`),
			errMsg:          "no matching pattern for topic",
			configVersion:   "v0.0.1",
			wantTopic:       "factory/line/cell",
			wantPayload:     `{"variable": "temp", "value": 23.5, "unit": "celsius"}`,
			wantError:       "no matching pattern for topic",
			wantConfigVersion: "v0.0.1",
			wantErr:         nil,
		},
		{
			name:            "payload extraction error",
			topic:           "cml03/ple/furnace",
			payload:         []byte(`{invalid json}`),
			errMsg:          "invalid JSON payload",
			configVersion:   "v0.0.1",
			wantTopic:       "cml03/ple/furnace",
			wantPayload:     `{invalid json}`,
			wantError:       "invalid JSON payload",
			wantConfigVersion: "v0.0.1",
			wantErr:         nil,
		},
		{
			name:            "missing required field",
			topic:           "cml03/ple/furnace",
			payload:         []byte(`{"variable": "temp"}`),
			errMsg:          "missing field 'value' in payload",
			configVersion:   "v0.0.2",
			wantTopic:       "cml03/ple/furnace",
			wantPayload:     `{"variable": "temp"}`,
			wantError:       "missing field 'value' in payload",
			wantConfigVersion: "v0.0.2",
			wantErr:         nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dlq, err := topicparser.BuildDLQMessage(tt.topic, tt.payload, tt.errMsg, tt.configVersion)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantTopic, dlq.OriginalTopic)
				assert.Equal(t, tt.wantPayload, string(dlq.OriginalPayload))
				assert.Equal(t, tt.wantError, dlq.ErrorMessage)
				assert.Equal(t, tt.wantConfigVersion, dlq.ConfigVersion)
				assert.NotEmpty(t, dlq.Timestamp)
			}
		})
	}
}

func TestBuildDLQMessage_PreservesInput(t *testing.T) {
	originalTopic := "cml03/ple/furnace/temp/sensor1"
	originalPayload := []byte(`{"device_id": "sensor-001", "reading": 23.5}`)
	errorMsg := "unexpected error occurred"
	configVersion := "v0.1.0"

	dlq, err := topicparser.BuildDLQMessage(originalTopic, originalPayload, errorMsg, configVersion)

	assert.NoError(t, err)
	assert.Equal(t, originalTopic, dlq.OriginalTopic)
	assert.Equal(t, originalPayload, dlq.OriginalPayload)
	assert.Equal(t, errorMsg, dlq.ErrorMessage)
	assert.Equal(t, configVersion, dlq.ConfigVersion)
	assert.NotEmpty(t, dlq.Timestamp)
}

func TestBuildDLQMessage_EmptyErrorMessage(t *testing.T) {
	dlq, err := topicparser.BuildDLQMessage(
		"topic/path",
		[]byte(`{"data": "test"}`),
		"",
		"v0.0.1",
	)

	assert.NoError(t, err)
	assert.Equal(t, "", dlq.ErrorMessage)
}
