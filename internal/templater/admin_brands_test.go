package templater

import (
	"io"
	"strings"
	"testing"

	"github.com/leemartin77/mincepierank.co.uk/internal/storage"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
)

func renderPage(t *testing.T, page string, data PageData) string {
	t.Helper()
	tr := NewHtmlTemplater()
	r, err := tr.GeneratePage(page, data)
	if err != nil {
		t.Fatalf("failed to render %s: %v", page, err)
	}
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read %s: %v", page, err)
	}
	return string(out)
}

func adminHead() PageDataHead {
	return PageDataHead{Title: "Admin", IsAdminPage: true, MenuSettings: MenuSettings{ActiveYear: 2026}}
}

func TestAdminBrandsListPage(t *testing.T) {
	out := renderPage(t, "admin:brands", PageData{
		Head: adminHead(),
		PageData: map[string]interface{}{
			"Brands": []storage.AdminMakerRow{
				{Maker: types.Maker{Id: "bramwell", Name: "Bramwell's", Website: "https://bramwells.example"}, PieCount: 3},
				{Maker: types.Maker{Id: "greggs", Name: "Greggs"}, PieCount: 0},
			},
			"PageErrors": []string{"Cannot delete: 3 pies still reference this brand."},
		},
	})
	for _, want := range []string{
		"Bramwell&#39;s", "/admin/brands/bramwell", ">3<", "alert-danger",
		"still reference this brand", "/admin/brands/new",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("brands list missing %q", want)
		}
	}
}

func TestAdminBrandNewPage(t *testing.T) {
	out := renderPage(t, "admin:brandnew", PageData{
		Head: adminHead(),
		PageData: map[string]interface{}{
			"Form":       map[string]string{"id": "bram", "name": "Bram", "website": "https://bram.example"},
			"PageErrors": nil,
		},
	})
	for _, want := range []string{`action="/admin/brands"`, `value="bram"`, `value="Bram"`, `name="logo"`, `type="file"`} {
		if !strings.Contains(out, want) {
			t.Errorf("brandnew form missing %q", want)
		}
	}
}

func TestAdminBrandEditPage(t *testing.T) {
	out := renderPage(t, "admin:brand", PageData{
		Head: adminHead(),
		PageData: map[string]interface{}{
			"Maker":      &types.Maker{Id: "bramwell", Name: "Bramwell's", Logo: "/brands/bramwell/logo.png", Website: "https://bramwells.example"},
			"PieCount":   int64(2),
			"PageErrors": nil,
		},
	})
	for _, want := range []string{`hx-put="/admin/brands/bramwell"`, `value="https://bramwells.example"`, "/brands/bramwell/logo.png", `hx-delete="/admin/brands/bramwell"`, "2 pie(s)"} {
		if !strings.Contains(out, want) {
			t.Errorf("brand edit page missing %q", want)
		}
	}
}
