package topicparser

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseTopic(topic, pattern string, metadataConfig []MetadataEntry) (*ExtractedMetadata, error) {
	fields := make(map[string]string)

	for _, entry := range metadataConfig {
		if entry.Constant != nil {
			fields[entry.TagName] = entry.Constant.Value
		} else if entry.TopicSegment != nil {
			value, err := extractSegment(topic, entry.TopicSegment.Value)
			if err != nil {
				return nil, err
			}
			fields[entry.TagName] = value
		}
	}

	return &ExtractedMetadata{Fields: fields}, nil
}

func extractSegment(topic, segmentValue string) (string, error) {
	topicSegments := strings.Split(topic, "/")

	if strings.Contains(segmentValue, ",") {
		return extractMultiIndex(topicSegments, segmentValue)
	}

	if strings.Contains(segmentValue, ":") {
		return extractRange(topicSegments, segmentValue)
	}

	return extractSingleIndex(topicSegments, segmentValue)
}

func extractSingleIndex(segments []string, indexStr string) (string, error) {
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return "", fmt.Errorf("%w: invalid index syntax %q", ErrInvalidTopicStructure, indexStr)
	}

	if index < 0 || index >= len(segments) {
		return "", fmt.Errorf("%w: index %d out of bounds for %d segments", ErrInvalidTopicStructure, index, len(segments))
	}

	return segments[index], nil
}

func extractRange(segments []string, rangeStr string) (string, error) {
	parts := strings.Split(rangeStr, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("%w: invalid range syntax %q", ErrInvalidTopicStructure, rangeStr)
	}

	startStr := parts[0]
	start, err := strconv.Atoi(startStr)
	if err != nil {
		return "", fmt.Errorf("%w: invalid range start %q", ErrInvalidTopicStructure, startStr)
	}

	if start < 0 || start >= len(segments) {
		return "", fmt.Errorf("%w: range start %d out of bounds for %d segments", ErrInvalidTopicStructure, start, len(segments))
	}

	return strings.Join(segments[start:], "/"), nil
}

func extractMultiIndex(segments []string, indicesStr string) (string, error) {
	indexStrs := strings.Split(indicesStr, ",")
	var result []string

	for _, indexStr := range indexStrs {
		index, err := strconv.Atoi(strings.TrimSpace(indexStr))
		if err != nil {
			return "", fmt.Errorf("%w: invalid index in multi-index %q", ErrInvalidTopicStructure, indexStr)
		}

		if index < 0 || index >= len(segments) {
			return "", fmt.Errorf("%w: index %d out of bounds for %d segments", ErrInvalidTopicStructure, index, len(segments))
		}

		result = append(result, segments[index])
	}

	return strings.Join(result, "/"), nil
}
