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

func TestHeadPublicUsesDrawer(t *testing.T) {
	out := renderHead(t, PageDataHead{
		Title:        "Home",
		MenuSettings: MenuSettings{ActiveYear: 2026},
	})
	if !strings.Contains(out, `id="mpr-menu-pill"`) || !strings.Contains(out, `id="mpr-menu"`) {
		t.Errorf("public head should render the drawer nav, got:\n%s", out)
	}
	if strings.Contains(out, "navbar-toggler") {
		t.Errorf("public head should not render the bootstrap navbar, got:\n%s", out)
	}
	for _, want := range []string{"All Pies", "Brands", "Categories", "About", "My Rankings", "/years/2026/all-pies"} {
		if !strings.Contains(out, want) {
			t.Errorf("public drawer is missing nav item %q", want)
		}
	}
}

func TestHeadSnowCanvas(t *testing.T) {
	out := renderHead(t, PageDataHead{Title: "Home", Snow: true, MenuSettings: MenuSettings{ActiveYear: 2026}})
	if !strings.Contains(out, `<canvas id="canvas"`) {
		t.Errorf("snowing public head should render the canvas, got:\n%s", out)
	}

	off := renderHead(t, PageDataHead{Title: "Home", MenuSettings: MenuSettings{ActiveYear: 2026}})
	if strings.Contains(off, "<canvas") {
		t.Errorf("non-snowing public head should not render the canvas, got:\n%s", off)
	}

	admin := renderHead(t, PageDataHead{Title: "Admin", IsAdminPage: true, Snow: true})
	if strings.Contains(admin, "<canvas") {
		t.Errorf("admin head should never render the canvas, got:\n%s", admin)
	}
}

func TestHeadAdminUsesNavbar(t *testing.T) {
	out := renderHead(t, PageDataHead{
		Title:        "Admin",
		IsAdminPage:  true,
		MenuSettings: MenuSettings{ActiveYear: 2026, SignedIn: true},
	})
	if !strings.Contains(out, "navbar-toggler") {
		t.Errorf("admin head should keep the bootstrap navbar, got:\n%s", out)
	}
	if strings.Contains(out, "mpr-menu-pill") {
		t.Errorf("admin head should not render the drawer pill, got:\n%s", out)
	}
}

func TestHeadConditionalNavItems(t *testing.T) {
	cases := []struct {
		name       string
		settings   MenuSettings
		mustHave   []string
		mustNotHas []string
	}{
		{"pies disabled", MenuSettings{ActiveYear: 2026}, nil, []string{"Community Pies", "Submit a Pie", "My Pies"}},
		{"enabled, signed out", MenuSettings{ActiveYear: 2026, UserPiesEnabled: true}, []string{"Community Pies"}, []string{"Submit a Pie", "My Pies"}},
		{"enabled, signed in", MenuSettings{ActiveYear: 2026, UserPiesEnabled: true, SignedIn: true}, []string{"Community Pies", "Submit a Pie", "My Pies", "Log Out"}, []string{"Sign In"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := renderHead(t, PageDataHead{Title: "Home", MenuSettings: tc.settings})
			for _, want := range tc.mustHave {
				if !strings.Contains(out, want) {
					t.Errorf("expected nav to contain %q, got:\n%s", want, out)
				}
			}
			for _, notWant := range tc.mustNotHas {
				if strings.Contains(out, notWant) {
					t.Errorf("expected nav NOT to contain %q, got:\n%s", notWant, out)
				}
			}
		})
	}
}
