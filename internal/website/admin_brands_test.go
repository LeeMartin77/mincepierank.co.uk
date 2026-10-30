package website

import (
	"reflect"
	"testing"
)

func TestBrandIDFromName(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"simple", "Bramwell", "bramwell"},
		{"spaces to dashes", "Mrs Smith Pies", "mrs-smith-pies"},
		{"apostrophes stripped", "Bramwell's", "bramwell-s"},
		{"punctuation collapsed", "Co-op!! Cakes", "co-op-cakes"},
		{"trims dashes", "  --Greggs--  ", "greggs"},
		{"already slugged", "bramwell-pies", "bramwell-pies"},
		{"numbers kept", "9 Lives Pies", "9-lives-pies"},
		{"no alphanumerics", "!!!", ""},
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := brandIDFromName(tt.in); got != tt.want {
				t.Errorf("brandIDFromName(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestValidateBrandForm(t *testing.T) {
	tests := []struct {
		name     string
		formName string
		formID   string
		wantID   string
		wantErrs int
	}{
		{"id derived from name", "Bramwell's", "", "bramwell-s", 0},
		{"explicit id kept", "Bramwell's", "bramwells", "bramwells", 0},
		{"explicit valid id", "Anything", "a1-2", "a1-2", 0},
		{
			name:     "name required",
			formName: "", formID: "x",
			wantErrs: 1,
		},
		{
			name:     "uppercase id rejected (path safety)",
			formName: "Bramwell", formID: "Bramwell",
			wantErrs: 1,
		},
		{
			name:     "path traversal id rejected",
			formName: "Bramwell", formID: "../evil",
			wantErrs: 1,
		},
		{
			name:     "spaces in id rejected",
			formName: "Bramwell", formID: "bram well",
			wantErrs: 1,
		},
		{
			name:     "id starting with dash rejected",
			formName: "Bramwell", formID: "-bramwell",
			wantErrs: 1,
		},
		{
			name:     "name with no alphanumerics cannot derive an id",
			formName: "!!!", formID: "",
			wantErrs: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, errs := validateBrandForm(tt.formName, tt.formID)
			if len(errs) != tt.wantErrs {
				t.Fatalf("validateBrandForm(%q, %q) errs = %v, want %d errors", tt.formName, tt.formID, errs, tt.wantErrs)
			}
			if tt.wantErrs == 0 && id != tt.wantID {
				t.Errorf("validateBrandForm(%q, %q) id = %q, want %q", tt.formName, tt.formID, id, tt.wantID)
			}
		})
	}
}

func TestValidateBrandFormReportsAllProblems(t *testing.T) {
	_, errs := validateBrandForm("", "Bad ID")
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %v", errs)
	}
	if !reflect.DeepEqual(errs[0], "Brand name is required") {
		t.Errorf("first error = %q", errs[0])
	}
}
