package website

import (
	"net/url"
	"reflect"
	"testing"
)

func TestParseRatings(t *testing.T) {
	valid := func() url.Values {
		return url.Values{"pastry": {"5"}, "filling": {"3"}, "topping": {"4"}, "looks": {"2"}, "value": {"1"}}
	}
	with := func(k, v string) url.Values {
		f := valid()
		if v == "" {
			f.Del(k)
		} else {
			f.Set(k, v)
		}
		return f
	}

	tests := []struct {
		name     string
		form     url.Values
		want     map[string]int32
		wantErrs []string
	}{
		{
			name: "all valid",
			form: valid(),
			want: map[string]int32{"pastry": 5, "filling": 3, "topping": 4, "looks": 2, "value": 1},
		},
		{name: "missing pastry", form: with("pastry", ""), wantErrs: []string{"Pastry is required"}},
		{name: "above range", form: with("pastry", "6"), wantErrs: []string{"Pastry must be a whole number from 1 to 5"}},
		{name: "zero rejected", form: with("looks", "0"), wantErrs: []string{"Looks must be a whole number from 1 to 5"}},
		{name: "negative", form: with("topping", "-1"), wantErrs: []string{"Topping must be a whole number from 1 to 5"}},
		{name: "not a number", form: with("value", "abc"), wantErrs: []string{"Value for Money must be a whole number from 1 to 5"}},
		{name: "decimal", form: with("filling", "2.5"), wantErrs: []string{"Filling must be a whole number from 1 to 5"}},
		{
			name: "empty form reports every field",
			form: url.Values{},
			wantErrs: []string{
				"Filling is required", "Pastry is required", "Topping is required",
				"Looks is required", "Value for Money is required",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, errs := parseRatings(tt.form.Get)
			if !reflect.DeepEqual(errs, tt.wantErrs) {
				t.Fatalf("errs = %q, want %q", errs, tt.wantErrs)
			}
			if tt.wantErrs == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsed = %v, want %v", got, tt.want)
			}
		})
	}
}
