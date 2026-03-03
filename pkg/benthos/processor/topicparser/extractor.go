package topicparser

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func ExtractPayload(payload []byte, payloadConfig PayloadConfig) (map[string]string, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, fmt.Errorf("%w: invalid JSON payload", ErrPayloadExtractionFailed)
	}

	result := make(map[string]string)

	// Extract variable
	variable, err := extractField(data, payloadConfig.Variable)
	if err != nil {
		return nil, err
	}
	result["variable"] = variable

	// Extract value
	value, err := extractField(data, payloadConfig.Value)
	if err != nil {
		return nil, err
	}
	result["value"] = value

	// Extract unit
	unit, err := extractField(data, payloadConfig.Unit)
	if err != nil {
		return nil, err
	}
	result["unit"] = unit

	return result, nil
}

func extractField(data map[string]interface{}, path string) (string, error) {
	if !strings.HasPrefix(path, "$.") {
		return "", fmt.Errorf("%w: path must start with $.", ErrPayloadExtractionFailed)
	}

	// Remove leading $.
	path = path[2:]

	// Navigate through the path
	var current interface{} = data
	parts := parsePathSegments(path)

	for _, part := range parts {
		if current == nil {
			return "", fmt.Errorf("%w: path component does not exist", ErrPayloadExtractionFailed)
		}

		if arrayIndex, isArray := parseArrayIndex(part); isArray {
			// Handle array indexing
			arr, ok := current.([]interface{})
			if !ok {
				return "", fmt.Errorf("%w: expected array but got different type", ErrPayloadExtractionFailed)
			}

			if arrayIndex < 0 || arrayIndex >= len(arr) {
				return "", fmt.Errorf("%w: array index out of bounds", ErrPayloadExtractionFailed)
			}

			current = arr[arrayIndex]
		} else {
			// Handle object field access
			obj, ok := current.(map[string]interface{})
			if !ok {
				return "", fmt.Errorf("%w: expected object but got different type", ErrPayloadExtractionFailed)
			}

			var exists bool
			current, exists = obj[part]
			if !exists {
				return "", fmt.Errorf("%w: field %q not found", ErrPayloadExtractionFailed, part)
			}
		}
	}

	// Convert final value to string
	if current == nil {
		return "", fmt.Errorf("%w: value is null", ErrPayloadExtractionFailed)
	}

	return valueToString(current)
}

func parsePathSegments(path string) []string {
	var segments []string
	var current strings.Builder

	for i := 0; i < len(path); i++ {
		char := path[i]

		if char == '[' {
			// Save the field name before [
			if current.Len() > 0 {
				segments = append(segments, current.String())
				current.Reset()
			}

			// Find matching ]
			closeIdx := strings.Index(path[i:], "]")
			if closeIdx == -1 {
				break
			}

			// Extract array index part including brackets
			indexPart := path[i : i+closeIdx+1]
			segments = append(segments, indexPart)
			i += closeIdx
		} else if char == '.' {
			if current.Len() > 0 {
				segments = append(segments, current.String())
				current.Reset()
			}
		} else {
			current.WriteByte(char)
		}
	}

	if current.Len() > 0 {
		segments = append(segments, current.String())
	}

	return segments
}

func parseArrayIndex(segment string) (int, bool) {
	if !strings.HasPrefix(segment, "[") || !strings.HasSuffix(segment, "]") {
		return 0, false
	}

	indexStr := segment[1 : len(segment)-1]
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return 0, false
	}

	return index, true
}

func valueToString(val interface{}) (string, error) {
	switch v := val.(type) {
	case string:
		return v, nil
	case float64:
		// Check if it's a whole number
		if v == float64(int64(v)) {
			return fmt.Sprintf("%.0f", v), nil
		}
		return fmt.Sprintf("%g", v), nil
	case bool:
		return fmt.Sprintf("%v", v), nil
	case nil:
		return "", fmt.Errorf("%w: value is null", ErrPayloadExtractionFailed)
	default:
		// For complex types like objects/arrays, return error
		return "", fmt.Errorf("%w: unsupported value type", ErrPayloadExtractionFailed)
	}
}
