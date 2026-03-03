package topicparser_test

import (
	"testing"

	"github.com/MGYOSBEL/pathfinder/pkg/benthos/processor/topicparser"
	"github.com/stretchr/testify/assert"
)

func TestExtractPayload_ValidExtractions(t *testing.T) {
	tests := []struct {
		name           string
		payload        []byte
		payloadConfig  topicparser.PayloadConfig
		wantVariable   string
		wantValue      string
		wantUnit       string
		wantErr        error
	}{
		{
			name:    "simple top-level fields",
			payload: []byte(`{"variable": "temp-01", "value": "23.5", "unit": "celsius"}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.variable",
				Value:    "$.value",
				Unit:     "$.unit",
			},
			wantVariable: "temp-01",
			wantValue:    "23.5",
			wantUnit:     "celsius",
			wantErr:      nil,
		},
		{
			name:    "nested object fields",
			payload: []byte(`{"sensor": {"name": "temp-01", "reading": 23.5, "unit": "celsius"}}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.sensor.name",
				Value:    "$.sensor.reading",
				Unit:     "$.sensor.unit",
			},
			wantVariable: "temp-01",
			wantValue:    "23.5",
			wantUnit:     "celsius",
			wantErr:      nil,
		},
		{
			name:    "deeply nested fields",
			payload: []byte(`{"data": {"device": {"sensor": {"value": 42.1}}}}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.data.device.sensor.value",
				Value:    "$.data.device.sensor.value",
				Unit:     "$.data.device.sensor.value",
			},
			wantVariable: "42.1",
			wantValue:    "42.1",
			wantUnit:     "42.1",
			wantErr:      nil,
		},
		{
			name:    "string conversion from number",
			payload: []byte(`{"value": 123}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.value",
				Value:    "$.value",
				Unit:     "$.value",
			},
			wantVariable: "123",
			wantValue:    "123",
			wantUnit:     "123",
			wantErr:      nil,
		},
		{
			name:    "string conversion from float",
			payload: []byte(`{"reading": 23.456}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.reading",
				Value:    "$.reading",
				Unit:     "$.reading",
			},
			wantVariable: "23.456",
			wantValue:    "23.456",
			wantUnit:     "23.456",
			wantErr:      nil,
		},
		{
			name:    "string conversion from boolean true",
			payload: []byte(`{"active": true}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.active",
				Value:    "$.active",
				Unit:     "$.active",
			},
			wantVariable: "true",
			wantValue:    "true",
			wantUnit:     "true",
			wantErr:      nil,
		},
		{
			name:    "string conversion from boolean false",
			payload: []byte(`{"active": false}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.active",
				Value:    "$.active",
				Unit:     "$.active",
			},
			wantVariable: "false",
			wantValue:    "false",
			wantUnit:     "false",
			wantErr:      nil,
		},
		{
			name:    "array indexing - first element",
			payload: []byte(`{"readings": [{"value": 10}, {"value": 20}]}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.readings[0].value",
				Value:    "$.readings[0].value",
				Unit:     "$.readings[0].value",
			},
			wantVariable: "10",
			wantValue:    "10",
			wantUnit:     "10",
			wantErr:      nil,
		},
		{
			name:    "array indexing - second element",
			payload: []byte(`{"readings": [{"value": 10}, {"value": 20}]}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.readings[1].value",
				Value:    "$.readings[1].value",
				Unit:     "$.readings[1].value",
			},
			wantVariable: "20",
			wantValue:    "20",
			wantUnit:     "20",
			wantErr:      nil,
		},
		{
			name:    "array of primitives",
			payload: []byte(`{"values": [100, 200, 300]}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.values[0]",
				Value:    "$.values[1]",
				Unit:     "$.values[2]",
			},
			wantVariable: "100",
			wantValue:    "200",
			wantUnit:     "300",
			wantErr:      nil,
		},
		{
			name:    "mix of constants and extractions",
			payload: []byte(`{"sensor": "temp-01", "reading": 25.3}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.sensor",
				Value:    "$.reading",
				Unit:     "$.sensor",
			},
			wantVariable: "temp-01",
			wantValue:    "25.3",
			wantUnit:     "temp-01",
			wantErr:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := topicparser.ExtractPayload(tt.payload, tt.payloadConfig)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantVariable, got["variable"])
				assert.Equal(t, tt.wantValue, got["value"])
				assert.Equal(t, tt.wantUnit, got["unit"])
			}
		})
	}
}

func TestExtractPayload_MissingFields(t *testing.T) {
	tests := []struct {
		name          string
		payload       []byte
		payloadConfig topicparser.PayloadConfig
		wantErr       error
	}{
		{
			name:    "missing variable field",
			payload: []byte(`{"value": "test", "unit": "m"}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.missing_variable",
				Value:    "$.value",
				Unit:     "$.unit",
			},
			wantErr: topicparser.ErrPayloadExtractionFailed,
		},
		{
			name:    "missing value field",
			payload: []byte(`{"variable": "test", "unit": "m"}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.variable",
				Value:    "$.missing_value",
				Unit:     "$.unit",
			},
			wantErr: topicparser.ErrPayloadExtractionFailed,
		},
		{
			name:    "missing unit field",
			payload: []byte(`{"variable": "test", "value": "10"}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.variable",
				Value:    "$.value",
				Unit:     "$.missing_unit",
			},
			wantErr: topicparser.ErrPayloadExtractionFailed,
		},
		{
			name:    "nested field missing",
			payload: []byte(`{"sensor": {"reading": 25.3}}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.sensor.missing_field",
				Value:    "$.sensor.reading",
				Unit:     "$.sensor.unit",
			},
			wantErr: topicparser.ErrPayloadExtractionFailed,
		},
		{
			name:    "intermediate object missing",
			payload: []byte(`{"data": {}}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.data.device.sensor.value",
				Value:    "$.data.device.sensor.value",
				Unit:     "$.data.device.sensor.value",
			},
			wantErr: topicparser.ErrPayloadExtractionFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := topicparser.ExtractPayload(tt.payload, tt.payloadConfig)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestExtractPayload_NullValues(t *testing.T) {
	tests := []struct {
		name          string
		payload       []byte
		payloadConfig topicparser.PayloadConfig
		wantErr       error
	}{
		{
			name:    "null variable field",
			payload: []byte(`{"variable": null, "value": "10", "unit": "m"}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.variable",
				Value:    "$.value",
				Unit:     "$.unit",
			},
			wantErr: topicparser.ErrPayloadExtractionFailed,
		},
		{
			name:    "null value field",
			payload: []byte(`{"variable": "temp", "value": null, "unit": "m"}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.variable",
				Value:    "$.value",
				Unit:     "$.unit",
			},
			wantErr: topicparser.ErrPayloadExtractionFailed,
		},
		{
			name:    "null unit field",
			payload: []byte(`{"variable": "temp", "value": "10", "unit": null}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.variable",
				Value:    "$.value",
				Unit:     "$.unit",
			},
			wantErr: topicparser.ErrPayloadExtractionFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := topicparser.ExtractPayload(tt.payload, tt.payloadConfig)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestExtractPayload_InvalidJSON(t *testing.T) {
	tests := []struct {
		name          string
		payload       []byte
		payloadConfig topicparser.PayloadConfig
		wantErr       error
	}{
		{
			name:    "malformed JSON",
			payload: []byte(`{invalid json}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.variable",
				Value:    "$.value",
				Unit:     "$.unit",
			},
			wantErr: topicparser.ErrPayloadExtractionFailed,
		},
		{
			name:    "empty payload",
			payload: []byte(``),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.variable",
				Value:    "$.value",
				Unit:     "$.unit",
			},
			wantErr: topicparser.ErrPayloadExtractionFailed,
		},
		{
			name:    "not JSON (plain text)",
			payload: []byte(`just some text`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.variable",
				Value:    "$.value",
				Unit:     "$.unit",
			},
			wantErr: topicparser.ErrPayloadExtractionFailed,
		},
		{
			name:    "JSON array instead of object",
			payload: []byte(`[1, 2, 3]`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.variable",
				Value:    "$.value",
				Unit:     "$.unit",
			},
			wantErr: topicparser.ErrPayloadExtractionFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := topicparser.ExtractPayload(tt.payload, tt.payloadConfig)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestExtractPayload_ArrayIndexing(t *testing.T) {
	tests := []struct {
		name          string
		payload       []byte
		payloadConfig topicparser.PayloadConfig
		wantErr       error
	}{
		{
			name:    "index out of bounds",
			payload: []byte(`{"readings": [{"value": 10}]}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.readings[5].value",
				Value:    "$.readings[5].value",
				Unit:     "$.readings[5].value",
			},
			wantErr: topicparser.ErrPayloadExtractionFailed,
		},
		{
			name:    "negative index",
			payload: []byte(`{"readings": [{"value": 10}]}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.readings[-1].value",
				Value:    "$.readings[-1].value",
				Unit:     "$.readings[-1].value",
			},
			wantErr: topicparser.ErrPayloadExtractionFailed,
		},
		{
			name:    "invalid array index syntax",
			payload: []byte(`{"readings": [{"value": 10}]}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.readings[abc].value",
				Value:    "$.readings[abc].value",
				Unit:     "$.readings[abc].value",
			},
			wantErr: topicparser.ErrPayloadExtractionFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := topicparser.ExtractPayload(tt.payload, tt.payloadConfig)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestExtractPayload_EdgeCases(t *testing.T) {
	tests := []struct {
		name           string
		payload        []byte
		payloadConfig  topicparser.PayloadConfig
		wantVariable   string
		wantValue      string
		wantUnit       string
		wantErr        error
	}{
		{
			name:    "empty object",
			payload: []byte(`{}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.variable",
				Value:    "$.value",
				Unit:     "$.unit",
			},
			wantErr: topicparser.ErrPayloadExtractionFailed,
		},
		{
			name:    "empty string values are valid",
			payload: []byte(`{"variable": "", "value": "", "unit": ""}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.variable",
				Value:    "$.value",
				Unit:     "$.unit",
			},
			wantVariable: "",
			wantValue:    "",
			wantUnit:     "",
			wantErr:      nil,
		},
		{
			name:    "zero values are valid",
			payload: []byte(`{"variable": 0, "value": 0, "unit": 0}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.variable",
				Value:    "$.value",
				Unit:     "$.unit",
			},
			wantVariable: "0",
			wantValue:    "0",
			wantUnit:     "0",
			wantErr:      nil,
		},
		{
			name:    "special characters in strings",
			payload: []byte(`{"variable": "temp/sensor#1", "value": "23.5°C", "unit": "℃"}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.variable",
				Value:    "$.value",
				Unit:     "$.unit",
			},
			wantVariable: "temp/sensor#1",
			wantValue:    "23.5°C",
			wantUnit:     "℃",
			wantErr:      nil,
		},
		{
			name:    "unicode characters",
			payload: []byte(`{"variable": "温度", "value": "湿度", "unit": "压力"}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.variable",
				Value:    "$.value",
				Unit:     "$.unit",
			},
			wantVariable: "温度",
			wantValue:    "湿度",
			wantUnit:     "压力",
			wantErr:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := topicparser.ExtractPayload(tt.payload, tt.payloadConfig)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantVariable, got["variable"])
				assert.Equal(t, tt.wantValue, got["value"])
				assert.Equal(t, tt.wantUnit, got["unit"])
			}
		})
	}
}

func TestExtractPayload_RealWorldScenarios(t *testing.T) {
	tests := []struct {
		name           string
		payload        []byte
		payloadConfig  topicparser.PayloadConfig
		wantVariable   string
		wantValue      string
		wantUnit       string
		wantErr        error
	}{
		{
			name: "MQTT sensor message",
			payload: []byte(`{
				"device_id": "sensor-001",
				"timestamp": "2026-03-03T06:30:00Z",
				"data": {
					"temperature": 23.5,
					"humidity": 65.2,
					"pressure": 1013.25
				}
			}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.device_id",
				Value:    "$.data.temperature",
				Unit:     "$.data.humidity",
			},
			wantVariable: "sensor-001",
			wantValue:    "23.5",
			wantUnit:     "65.2",
			wantErr:      nil,
		},
		{
			name: "Industrial machine metrics",
			payload: []byte(`{
				"metrics": [
					{"name": "temperature", "value": 78.5},
					{"name": "vibration", "value": 2.3},
					{"name": "runtime", "value": 145}
				]
			}`),
			payloadConfig: topicparser.PayloadConfig{
				Variable: "$.metrics[0].name",
				Value:    "$.metrics[0].value",
				Unit:     "$.metrics[1].value",
			},
			wantVariable: "temperature",
			wantValue:    "78.5",
			wantUnit:     "2.3",
			wantErr:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := topicparser.ExtractPayload(tt.payload, tt.payloadConfig)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantVariable, got["variable"])
				assert.Equal(t, tt.wantValue, got["value"])
				assert.Equal(t, tt.wantUnit, got["unit"])
			}
		})
	}
}
