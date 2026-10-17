package templatefuncs

import (
	"bytes"
	"html/template"
	"strings"
	"testing"
)

func TestNiceFloat(t *testing.T) {
	for in, want := range map[interface{}]string{
		float64(3.21): "3.21", float64(4): "4", float32(2.5): "2.50", int32(3): "3",
	} {
		got, err := NiceFloat(in)
		if err != nil || got != want {
			t.Errorf("NiceFloat(%v) = %q, %v; want %q", in, got, err, want)
		}
	}
	if _, err := NiceFloat("x"); err == nil {
		t.Error("NiceFloat(string) should error")
	}
}

func TestStarPercent(t *testing.T) {
	for in, want := range map[interface{}]string{
		float64(0): "0.0", float64(3.21): "64.2", float64(5): "100.0",
		int32(2): "40.0", float64(7): "100.0", float64(-1): "0.0",
	} {
		got, err := StarPercent(in)
		if err != nil || got != want {
			t.Errorf("StarPercent(%v) = %q, %v; want %q", in, got, err, want)
		}
	}
}

// Renders the real star templates and checks radio order and checked state.
func TestStarTemplates(t *testing.T) {
	tmpl := template.Must(template.New("").Funcs(Common()).
		ParseFiles("../../assets/templates/html/fragments/components/stars.gohtml"))

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "star-input", map[string]interface{}{
		"Name": "pastry", "Label": "Pastry", "Value": int32(2),
	}); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	last := -1
	for _, v := range []string{"1", "2", "3", "4", "5"} {
		i := strings.Index(out, `id="pastry-`+v+`"`)
		if i < last {
			t.Fatalf("radio %s out of order; must be 1..5", v)
		}
		last = i
	}
	if strings.Count(out, "checked") != 1 || !strings.Contains(out, `value="2" id="pastry-2" checked`) {
		t.Errorf("expected only value 2 checked:\n%s", out)
	}

	buf.Reset()
	if err := tmpl.ExecuteTemplate(&buf, "star-display", float64(3.21)); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "width: 64.2%") || !strings.Contains(buf.String(), `aria-label="3.21 out of 5"`) {
		t.Errorf("star-display output unexpected:\n%s", buf.String())
	}
}
