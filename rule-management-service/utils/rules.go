package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

type Rule struct {
	ID           int     `json:"id"`
	MetricName   string  `json:"metricName" binding:"required"`
	Condition    string  `json:"condition" binding:"required"`
	Threshold    float64 `json:"threshold" binding:"required"`
	AlertMessage string  `json:"alertMessage" binding:"required"`
}

var NextRuleID = 4 // Track the next available rule ID

// ReadRulesFromJSON reads the rules from the JSON file.
func ReadRulesFromJSON(filename string) ([]Rule, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []Rule{}, nil // No file, return an empty slice
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var rules []Rule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Update NextRuleID based on the rules in the file
	for _, rule := range rules {
		if rule.ID >= NextRuleID {
			NextRuleID = rule.ID + 1
		}
	}

	return rules, nil
}

// WriteRulesToJSON writes the given rules to the JSON file.
func WriteRulesToJSON(filename string, rules []Rule) error {
	data, err := json.MarshalIndent(rules, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}
