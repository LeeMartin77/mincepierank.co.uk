// Package templatefuncs holds the HTML template functions that need nothing
// beyond the standard library. internal/templater adds the storage-aware ones
// on top; cmd/template-preview uses these directly so the two can't drift.
package templatefuncs

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"strings"
	"time"
)

// Common returns a fresh FuncMap of the standard-library-only functions.
func Common() template.FuncMap {
	return template.FuncMap{
		"add":         Add,
		"dict":        Dict,
		"jsonserial":  JSONSerial,
		"ppencediv":   PenceDiv,
		"sub":         Sub,
		"timetohuman": TimeToHuman,
		"nicefloat":   NiceFloat,
		"starpercent": StarPercent,
		"starvalues":  StarValues,
	}
}

// Add returns the sum of two integers.
func Add(a, b int) int { return a + b }

// Sub returns the difference of two integers.
func Sub(a, b int) int { return a - b }

func Dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func JSONSerial(value interface{}) (string, error) {
	b, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return "Error serialising", err
	}
	return string(b), nil
}

func PenceDiv(value interface{}) (float32, error) {
	switch v := value.(type) {
	case float32:
		return v / 100, nil
	case int32:
		return float32(v) / 100, nil
	}
	return 0, fmt.Errorf("first value not floatable")
}

func TimeToHuman(value interface{}) (string, error) {
	tm, ok := value.(time.Time)
	if !ok {
		return "", fmt.Errorf("is not a time")
	}
	return tm.Format(time.RFC1123), nil
}

// toFloat accepts the numeric types ratings arrive as (float summaries,
// int32 user rankings).
func toFloat(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case int:
		return float64(v), nil
	}
	return 0, fmt.Errorf("value %v (%T) not floatable", value, value)
}

// NiceFloat formats to 2dp, dropping a trailing ".00".
func NiceFloat(value interface{}) (string, error) {
	f, err := toFloat(value)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(fmt.Sprintf("%.2f", float32(f)), ".00"), nil
}

// StarPercent converts a 0-5 rating into a 0-100 fill width, clamped.
func StarPercent(value interface{}) (string, error) {
	f, err := toFloat(value)
	if err != nil {
		return "", err
	}
	p := f / 5 * 100
	if p < 0 {
		p = 0
	} else if p > 100 {
		p = 100
	}
	return fmt.Sprintf("%.1f", p), nil
}

// StarValues lists the selectable ratings for the star-input template.
func StarValues() []int32 { return []int32{1, 2, 3, 4, 5} }
