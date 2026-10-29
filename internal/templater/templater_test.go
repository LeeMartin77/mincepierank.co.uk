package templater

import (
	"io"
	"strings"
	"testing"
)

func renderHead(t *testing.T, head PageDataHead) string {
	t.Helper()
	tr := NewHtmlTemplater()
	r, err := tr.GenerateTemplate("fragment:global:head", head)
	if err != nil {
		t.Fatalf("failed to render head: %v", err)
	}
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read head: %v", err)
	}
	return string(out)
}

func TestHeadBodyClassPublic(t *testing.T) {
	out := renderHead(t, PageDataHead{Title: "Home"})
	if !strings.Contains(out, `class="mpr-public"`) {
		t.Errorf("public head should carry body class mpr-public, got:\n%s", out)
	}
	if strings.Contains(out, `mpr-admin`) {
		t.Errorf("public head should not carry mpr-admin class, got:\n%s", out)
	}
}

func TestHeadBodyClassAdmin(t *testing.T) {
	out := renderHead(t, PageDataHead{Title: "Admin", IsAdminPage: true})
	if !strings.Contains(out, `class="mpr-admin"`) {
		t.Errorf("admin head should carry body class mpr-admin, got:\n%s", out)
	}
}

func TestHeadIncludesThemeStylesheet(t *testing.T) {
	out := renderHead(t, PageDataHead{Title: "Home"})
	if !strings.Contains(out, `/static/theme.css`) {
		t.Errorf("head should link the theme stylesheet, got:\n%s", out)
	}
}
