package topicparser

import (
	"slices"
	"strings"
)

func Match(topic string, configs []TopicParserConfig) []TopicParserConfig {
	matching := []TopicParserConfig{}

	// find all matching configs
	for _, config := range configs {
		if config.Enabled && matchPattern(config.Pattern, topic) {
			matching = append(matching, config)
		}
	}
	if len(matching) == 0 {
		return matching
	}

	// Sort by priority
	slices.SortFunc(matching,
		func(a TopicParserConfig, b TopicParserConfig) int { return b.Priority - a.Priority },
	)

	maxPriority := matching[0].Priority
	result := []TopicParserConfig{}

	for _, config := range matching {
		if config.Priority != maxPriority {
			break
		}
		result = append(result, config)
	}

	return result
}

func matchPattern(pattern, topic string) bool {
	var topicSegments []string
	if topic == "" {
		topicSegments = []string{}
	} else {
		topicSegments = strings.Split(topic, "/")
	}
	patternSegments := strings.Split(pattern, "/")

	for i, segment := range patternSegments {
		if segment == "#" {
			return i == len(patternSegments)-1
		}
		if i >= len(topicSegments) {
			return false
		}
		if segment == "+" {
			continue
		}
		if segment != topicSegments[i] {
			return false
		}
	}
	return len(topicSegments) == len(patternSegments)
}
