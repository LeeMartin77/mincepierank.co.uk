package templater

import (
	"io"
	"strings"
	"testing"
)

func executeUserRanking(t *testing.T, data map[string]any) string {
	t.Helper()
	wrapper := NewHtmlTemplater()
	r, err := wrapper.GenerateTemplate("fragment:components:userranking", data)
	if err != nil {
		t.Fatalf("GenerateTemplate: %v", err)
	}
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	return string(b)
}

func rankingData(signedIn bool) map[string]any {
	return map[string]any{
		"Ranking":        map[string]any{},
		"Readonly":       false,
		"ActiveYear":     int64(2026),
		"Pie":            map[string]any{"Year": int32(2026), "MakerId": "testbrand", "Id": "testpie"},
		"SignedIn":       signedIn,
		"SignInRedirect": "/years/2026/brands/testbrand/testpie",
	}
}

func TestUserRankingShowsSignInCTAWhenSignedOut(t *testing.T) {
	out := executeUserRanking(t, rankingData(false))

	if !strings.Contains(out, "Sign in to rank with your account") {
		t.Error("expected sign-in CTA banner for signed-out users")
	}
	if !strings.Contains(out, `/auth/signin?redirect=%2fyears%2f2026%2fbrands%2ftestbrand%2ftestpie`) {
		t.Errorf("expected sign-in link with redirect back to the pie page, got:\n%s", out)
	}
	// Voting stays open: the rating form must still render for signed-out users.
	if !strings.Contains(out, `hx-put="/years/2026/brands/testbrand/testpie/rating"`) {
		t.Error("expected rating form to still render for signed-out users")
	}
}

func TestUserRankingHidesSignInCTAWhenSignedIn(t *testing.T) {
	out := executeUserRanking(t, rankingData(true))

	if strings.Contains(out, "Sign in to rank with your account") {
		t.Errorf("expected no sign-in CTA banner for signed-in users, got:\n%s", out)
	}
	if !strings.Contains(out, `hx-put="/years/2026/brands/testbrand/testpie/rating"`) {
		t.Error("expected rating form for signed-in users")
	}
}
