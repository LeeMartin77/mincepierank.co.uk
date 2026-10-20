package website

import (
	"net/url"
	"reflect"
	"testing"
)

func TestValidateUserPieForm(t *testing.T) {
	tests := []struct {
		name      string
		form      url.Values
		wantErrs  []string
		wantLabel []string
	}{
		{
			name: "all valid",
			form: url.Values{
				"displayname": {"Ginsters"},
				"maker":       {"Ginsters"},
				"webLink":     {"https://ginsters.co.uk"},
			},
			wantErrs: nil,
		},
		{
			name: "missing displayname",
			form: url.Values{
				"maker": {"Ginsters"},
				"webLink": {"https://ginsters.co.uk"},
			},
			wantErrs: []string{"Display name is required"},
		},
		{
			name: "missing maker",
			form: url.Values{
				"displayname": {"Ginsters"},
				"webLink": {"https://ginsters.co.uk"},
			},
			wantErrs: []string{"Maker/brand is required"},
		},
		{
			name: "missing webLink",
			form: url.Values{
				"displayname": {"Ginsters"},
				"maker": {"Ginsters"},
			},
			wantErrs: []string{"Website link is required"},
		},
		{
			name: "all empty",
			form: url.Values{},
			wantErrs: []string{
				"Display name is required",
				"Maker/brand is required",
				"Website link is required",
			},
		},
		{
			name: "valid with labels",
			form: url.Values{
				"displayname": {"Ginsters"},
				"maker":       {"Ginsters"},
				"webLink":     {"https://ginsters.co.uk"},
				"labels":      {"organic, traditional, cornish"},
			},
			wantErrs:  nil,
			wantLabel: []string{"organic", "traditional", "cornish"},
		},
		{
			name: "labels with extra spaces",
			form: url.Values{
				"displayname": {"Ginsters"},
				"maker":       {"Ginsters"},
				"webLink":     {"https://ginsters.co.uk"},
				"labels":      {" organic ,  traditional  "},
			},
			wantErrs:  nil,
			wantLabel: []string{"organic", "traditional"},
		},
		{
			name: "labels with empty entries skipped",
			form: url.Values{
				"displayname": {"Ginsters"},
				"maker":       {"Ginsters"},
				"webLink":     {"https://ginsters.co.uk"},
				"labels":      {"organic,,traditional"},
			},
			wantErrs:  nil,
			wantLabel: []string{"organic", "traditional"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs, labels := validateUserPieForm(tt.form.Get)
			if !reflect.DeepEqual(errs, tt.wantErrs) {
				t.Errorf("validateUserPieForm() errors = %v, want %v", errs, tt.wantErrs)
			}
			if tt.wantLabel != nil && !reflect.DeepEqual(labels, tt.wantLabel) {
				t.Errorf("validateUserPieForm() labels = %v, want %v", labels, tt.wantLabel)
			}
		})
	}
}

func TestParsePackCount(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    int32
		wantErr bool
	}{
		{"valid", "12", 12, false},
		{"zero", "0", 0, false},
		{"negative", "-1", 0, true},
		{"not a number", "abc", 0, true},
		{"empty", "", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePackCount(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePackCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parsePackCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParsePackPrice(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    int32
		wantErr bool
	}{
		{"valid", "3.60", 360, false},
		{"rounds up", "3.605", 361, false},
		{"whole number", "4", 400, false},
		{"zero", "0", 0, false},
		{"negative", "-1", 0, true},
		{"not a number", "abc", 0, true},
		{"empty", "", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePackPrice(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePackPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parsePackPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsFresh(t *testing.T) {
	tests := []struct {
		name string
		val  string
		want bool
	}{
		{"on", "on", true},
		{"true", "true", true},
		{"off", "off", false},
		{"false", "false", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isFresh(tt.val)
			if got != tt.want {
				t.Errorf("isFresh(%q) = %v, want %v", tt.val, got, tt.want)
			}
		})
	}
}
