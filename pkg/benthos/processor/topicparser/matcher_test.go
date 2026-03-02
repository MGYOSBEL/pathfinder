package topicparser_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/MGYOSBEL/pathfinder/pkg/benthos/processor/topicparser"
)

// Helper to create test configs
func configWithPattern(id int, pattern string, priority int) topicparser.TopicParserConfig {
	return topicparser.TopicParserConfig{
		ID:       id,
		Name:     "config-" + strconv.Itoa(id),
		Pattern:  pattern,
		Priority: priority,
		Enabled:  true,
		Version:  "v0.0.1",
	}
}

// TestMatch_ExactMatching tests literal topic/pattern matching (no wildcards)
func TestMatch_ExactMatching(t *testing.T) {
	tests := []struct {
		name    string
		topic   string
		configs []topicparser.TopicParserConfig
		wantLen int
		wantIDs []int
	}{
		{
			name:  "single segment exact match",
			topic: "factory",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "factory", 10),
			},
			wantLen: 1,
			wantIDs: []int{1},
		},
		{
			name:  "single segment mismatch",
			topic: "factory",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "plant", 10),
			},
			wantLen: 0,
			wantIDs: []int{},
		},
		{
			name:  "multi-segment exact match",
			topic: "factory/line/cell",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "factory/line/cell", 10),
			},
			wantLen: 1,
			wantIDs: []int{1},
		},
		{
			name:  "multi-segment exact mismatch (different segment)",
			topic: "factory/line/cell",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "factory/line/equipment", 10),
			},
			wantLen: 0,
			wantIDs: []int{},
		},
		{
			name:  "length mismatch (topic has more segments)",
			topic: "factory/line/cell/equipment",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "factory/line/cell", 10),
			},
			wantLen: 0,
			wantIDs: []int{},
		},
		{
			name:  "length mismatch (pattern has more segments)",
			topic: "factory/line",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "factory/line/cell", 10),
			},
			wantLen: 0,
			wantIDs: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := topicparser.Match(tt.topic, tt.configs)
			assert.Equal(t, tt.wantLen, len(got), "expected %d matches", tt.wantLen)
			gotIDs := extractIDs(got)
			assert.Equal(t, tt.wantIDs, gotIDs)
		})
	}
}

// TestMatch_SingleWildcard tests + (single level) wildcard matching
func TestMatch_SingleWildcard(t *testing.T) {
	tests := []struct {
		name    string
		topic   string
		configs []topicparser.TopicParserConfig
		wantLen int
		wantIDs []int
	}{
		{
			name:  "single + matches any single segment",
			topic: "factory",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "+", 10),
			},
			wantLen: 1,
			wantIDs: []int{1},
		},
		{
			name:  "single + does not match multi-segment",
			topic: "factory/line",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "+", 10),
			},
			wantLen: 0,
			wantIDs: []int{},
		},
		{
			name:  "+ in middle matches",
			topic: "factory/line/cell",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "factory/+/cell", 10),
			},
			wantLen: 1,
			wantIDs: []int{1},
		},
		{
			name:  "+ in middle does not match mismatch",
			topic: "factory/line/equipment",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "factory/+/cell", 10),
			},
			wantLen: 0,
			wantIDs: []int{},
		},
		{
			name:  "multiple + matches exact segment count",
			topic: "a/b/c",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "+/+/+", 10),
			},
			wantLen: 1,
			wantIDs: []int{1},
		},
		{
			name:  "multiple + does not match fewer segments",
			topic: "a/b",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "+/+/+", 10),
			},
			wantLen: 0,
			wantIDs: []int{},
		},
		{
			name:  "multiple + does not match more segments",
			topic: "a/b/c/d",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "+/+/+", 10),
			},
			wantLen: 0,
			wantIDs: []int{},
		},
		{
			name:  "+ at start",
			topic: "factory/line/cell",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "+/line/cell", 10),
			},
			wantLen: 1,
			wantIDs: []int{1},
		},
		{
			name:  "+ at end",
			topic: "factory/line/cell",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "factory/line/+", 10),
			},
			wantLen: 1,
			wantIDs: []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := topicparser.Match(tt.topic, tt.configs)
			assert.Equal(t, tt.wantLen, len(got), "expected %d matches", tt.wantLen)
			gotIDs := extractIDs(got)
			assert.Equal(t, tt.wantIDs, gotIDs)
		})
	}
}

// TestMatch_HashWildcard tests # (multi-level) wildcard matching
func TestMatch_HashWildcard(t *testing.T) {
	tests := []struct {
		name    string
		topic   string
		configs []topicparser.TopicParserConfig
		wantLen int
		wantIDs []int
	}{
		{
			name:  "# alone matches single segment",
			topic: "factory",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "#", 10),
			},
			wantLen: 1,
			wantIDs: []int{1},
		},
		{
			name:  "# alone matches multiple segments",
			topic: "factory/line/cell/equipment",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "#", 10),
			},
			wantLen: 1,
			wantIDs: []int{1},
		},
		{
			name:  "# at end matches zero extra levels",
			topic: "factory",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "factory/#", 10),
			},
			wantLen: 1,
			wantIDs: []int{1},
		},
		{
			name:  "# at end matches one extra level",
			topic: "factory/line",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "factory/#", 10),
			},
			wantLen: 1,
			wantIDs: []int{1},
		},
		{
			name:  "# at end matches multiple extra levels",
			topic: "factory/line/cell/equipment/sensor",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "factory/#", 10),
			},
			wantLen: 1,
			wantIDs: []int{1},
		},
		{
			name:  "# does not match if prefix differs",
			topic: "plant/line/cell",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "factory/#", 10),
			},
			wantLen: 0,
			wantIDs: []int{},
		},
		{
			name:  "# after multiple segments",
			topic: "factory/line/cell/equipment/sensor/data",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "factory/line/#", 10),
			},
			wantLen: 1,
			wantIDs: []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := topicparser.Match(tt.topic, tt.configs)
			assert.Equal(t, tt.wantLen, len(got), "expected %d matches", tt.wantLen)
			gotIDs := extractIDs(got)
			assert.Equal(t, tt.wantIDs, gotIDs)
		})
	}
}

// TestMatch_Priority tests priority-based filtering (only highest priority matches returned)
func TestMatch_Priority(t *testing.T) {
	tests := []struct {
		name    string
		topic   string
		configs []topicparser.TopicParserConfig
		wantIDs []int // Expected IDs: only configs with highest priority
	}{
		{
			name:  "single config",
			topic: "factory/line/cell",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "+/+/+", 10),
			},
			wantIDs: []int{1},
		},
		{
			name:  "two configs, different priorities (only highest returned)",
			topic: "factory/line/cell",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "+/+/+", 10),
				configWithPattern(2, "factory/+/+", 20), // Higher priority
			},
			wantIDs: []int{2}, // Only config 2 (priority 20), config 1 filtered out
		},
		{
			name:  "three configs, mixed priorities (only highest returned)",
			topic: "factory/line/cell",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "factory/line/cell", 5), // Lowest - filtered
				configWithPattern(2, "factory/+/cell", 20),   // Highest - returned
				configWithPattern(3, "factory/line/+", 15),   // Middle - filtered
			},
			wantIDs: []int{2}, // Only config 2 (priority 20)
		},
		{
			name:  "tied priorities (broadcast only ties)",
			topic: "factory/line/cell",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "factory/+/cell", 10),
				configWithPattern(2, "factory/line/+", 10),   // Same priority as 1
				configWithPattern(3, "factory/line/cell", 5), // Lower priority - filtered
			},
			wantIDs: []int{1, 2}, // Both configs 1 and 2 (tied at priority 10), config 3 filtered
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := topicparser.Match(tt.topic, tt.configs)
			assert.Equal(t, len(tt.wantIDs), len(got), "expected %d matches", len(tt.wantIDs))
			gotIDs := extractIDs(got)
			assert.Equal(t, tt.wantIDs, gotIDs)
		})
	}
}

// TestMatch_Broadcast tests multiple matching configs with same priority (tied for highest)
func TestMatch_Broadcast(t *testing.T) {
	tests := []struct {
		name    string
		topic   string
		configs []topicparser.TopicParserConfig
		wantLen int
		wantIDs []int
	}{
		{
			name:  "two patterns both match but different priorities (highest returned)",
			topic: "factory/line/cell",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "+/+/+", 10),     // Lower priority - filtered
				configWithPattern(2, "factory/#", 20), // Higher priority - returned
			},
			wantLen: 1,
			wantIDs: []int{2},
		},
		{
			name:  "three patterns, only highest priority returned",
			topic: "factory/line/cell/equipment",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "factory/+/cell/equipment", 30), // Highest - returned
				configWithPattern(2, "factory/#", 20),                // Middle - filtered
				configWithPattern(3, "+/+/+/+", 10),                  // Lowest - filtered
			},
			wantLen: 1,
			wantIDs: []int{1},
		},
		{
			name:  "multiple patterns with same highest priority (broadcast ties)",
			topic: "a/b/c",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "+/+/+", 20), // Highest - returned
				configWithPattern(2, "a/+/+", 20), // Highest (tie) - returned
				configWithPattern(3, "+/b/+", 20), // Highest (tie) - returned
				configWithPattern(4, "a/b/c", 10), // Lower - filtered
			},
			wantLen: 3,
			wantIDs: []int{1, 2, 3}, // All three tied configs broadcast
		},
		{
			name:  "partial matches, only highest priority returned",
			topic: "factory/line/cell",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "factory/line/cell", 15),      // Matches, lower - filtered
				configWithPattern(2, "factory/line/equipment", 20), // Does not match
				configWithPattern(3, "factory/+/cell", 25),         // Matches, highest - returned
			},
			wantLen: 1,
			wantIDs: []int{3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := topicparser.Match(tt.topic, tt.configs)
			assert.Equal(t, tt.wantLen, len(got), "expected %d matches", tt.wantLen)
			gotIDs := extractIDs(got)
			assert.Equal(t, tt.wantIDs, gotIDs)
		})
	}
}

// TestMatch_EdgeCases tests edge cases and boundary conditions
func TestMatch_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		topic   string
		configs []topicparser.TopicParserConfig
		wantLen int
	}{
		{
			name:    "empty topic",
			topic:   "",
			configs: []topicparser.TopicParserConfig{configWithPattern(1, "+", 10)},
			wantLen: 0,
		},
		{
			name:    "empty configs list",
			topic:   "factory",
			configs: []topicparser.TopicParserConfig{},
			wantLen: 0,
		},
		{
			name:  "disabled config should not match",
			topic: "factory/line/cell",
			configs: []topicparser.TopicParserConfig{
				{
					Name:     "disabled",
					Pattern:  "factory/+/+",
					Priority: 10,
					Enabled:  false, // Disabled
					Version:  "v0.0.1",
				},
			},
			wantLen: 0,
		},
		{
			name:  "very deep topic",
			topic: "a/b/c/d/e/f/g/h/i/j",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "#", 10),
				configWithPattern(2, "a/#", 10),
			},
			wantLen: 2,
		},
		{
			name:  "single segment wildcard pattern vs multi-segment topic",
			topic: "a/b/c",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "+", 10),
			},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := topicparser.Match(tt.topic, tt.configs)
			assert.Equal(t, tt.wantLen, len(got), "expected %d matches", tt.wantLen)
		})
	}
}

// TestMatch_ComplexScenarios tests real-world-like scenarios with proper priority filtering
func TestMatch_ComplexScenarios(t *testing.T) {
	tests := []struct {
		name    string
		topic   string
		configs []topicparser.TopicParserConfig
		wantIDs []int
	}{
		{
			name:  "real scenario: factory data ingestion with priority-based selection",
			topic: "factory_a/line_1/cell_x/temperature",
			configs: []topicparser.TopicParserConfig{
				// Only the highest priority config is used, lower ones filtered
				configWithPattern(1, "factory_a/line_1/+/temperature", 30), // Highest - returned
				configWithPattern(2, "factory_a/+/+/temperature", 20),      // Filtered
				configWithPattern(3, "factory_a/#", 10),                    // Filtered
			},
			wantIDs: []int{1}, // Only highest priority
		},
		{
			name:  "multiple factories, correct factory + highest priority",
			topic: "factory_b/line/cell/humidity",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "factory_a/#", 10),
				configWithPattern(2, "factory_b/#", 20),      // Highest matching - returned
				configWithPattern(3, "factory_b/line/+", 15), // Filtered (lower priority)
				configWithPattern(4, "factory_c/#", 10),
			},
			wantIDs: []int{2}, // Only config 2 (highest matching priority)
		},
		{
			name:  "multiple high-priority configs (tied, so broadcast)",
			topic: "sensor/data/reading/value",
			configs: []topicparser.TopicParserConfig{
				configWithPattern(1, "#", 5),                      // Filtered (low priority)
				configWithPattern(2, "sensor/data/reading/+", 15), // Highest - returned
				configWithPattern(3, "sensor/+/+/value", 15),      // Tied highest - returned
				configWithPattern(4, "sensor/data/+/value", 10),   // Filtered (low priority)
			},
			wantIDs: []int{2, 3}, // Broadcast the two tied highest-priority configs
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := topicparser.Match(tt.topic, tt.configs)
			require.Equal(t, len(tt.wantIDs), len(got), "expected %d matches", len(tt.wantIDs))
			gotIDs := extractIDs(got)
			assert.Equal(t, tt.wantIDs, gotIDs)
		})
	}
}

// Helper function to extract IDs from matched configs
func extractIDs(configs []topicparser.TopicParserConfig) []int {
	ids := make([]int, len(configs))
	for i, cfg := range configs {
		ids[i] = cfg.ID
	}
	return ids
}
