package website

import (
	"testing"
)

func TestValidateCategoryForm(t *testing.T) {
	tests := []struct {
		name      string
		formLabel string
		formSlug  string
		wantLabel string
		wantSlug  string
		wantErrs  int
	}{
		{
			name:      "missing slug rejected, never derived",
			formLabel: "Gluten Free",
			wantLabel: "Gluten Free", wantErrs: 1,
		},
		{"explicit slug kept", "Dairy Free", "nofreefrom", "Dairy Free", "nofreefrom", 0},
		{"explicit valid slug", "Anything", "a1-2", "Anything", "a1-2", 0},
		{"label is trimmed", "  Vegan  ", "vegan", "Vegan", "vegan", 0},
		{
			name:     "label required",
			formSlug: "x", wantSlug: "x",
			wantErrs: 1,
		},
		{
			name:      "missing slug rejected even when label is unusable",
			formLabel: "!!!", wantLabel: "!!!",
			wantErrs: 1,
		},
		{
			name:      "uppercase slug rejected",
			formLabel: "Vegan", formSlug: "Vegan", wantLabel: "Vegan", wantSlug: "Vegan",
			wantErrs: 1,
		},
		{
			name:      "path traversal slug rejected",
			formLabel: "Vegan", formSlug: "../evil", wantLabel: "Vegan", wantSlug: "../evil",
			wantErrs: 1,
		},
		{
			name:      "spaces in slug rejected",
			formLabel: "Vegan", formSlug: "free from", wantLabel: "Vegan", wantSlug: "free from",
			wantErrs: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			label, slug, errs := validateCategoryForm(tt.formLabel, tt.formSlug)
			if label != tt.wantLabel || slug != tt.wantSlug {
				t.Errorf("validateCategoryForm(%q, %q) = (%q, %q), want (%q, %q)",
					tt.formLabel, tt.formSlug, label, slug, tt.wantLabel, tt.wantSlug)
			}
			if len(errs) != tt.wantErrs {
				t.Errorf("validateCategoryForm(%q, %q) errs = %v, want %d problems",
					tt.formLabel, tt.formSlug, errs, tt.wantErrs)
			}
		})
	}
}
