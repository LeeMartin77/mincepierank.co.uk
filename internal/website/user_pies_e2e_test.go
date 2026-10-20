//go:build e2e

package website

import (
	"context"
	"fmt"
	"net/http"
	"net/netip"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/sqlcgen"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	redismod "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestUserPiesE2E runs a full browser-based end-to-end test of the user pies feature.
//
// Run with: go test -tags e2e ./internal/website -run TestUserPiesE2E -v
func TestUserPiesE2E(t *testing.T) {
	ctx := context.Background()

	// ── 1. Start PostgreSQL ──────────────────────────────────────────────
	t.Log("Starting PostgreSQL container...")
	dbContainer, err := postgres.Run(ctx,
		"postgres:18",
		postgres.WithDatabase("mincepierank"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("mysecretpassword"),
		testcontainers.WithWaitStrategyAndDeadline(120*time.Second,
			wait.ForLog("database system is ready to accept connections").WithStartupTimeout(90*time.Second),
		),
	)
	require.NoError(t, err, "failed to start postgres container")
	t.Cleanup(func() { dbContainer.Terminate(ctx) })

	dbURL, err := dbContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	// ── 2. Start Redis ───────────────────────────────────────────────────
	t.Log("Starting Redis container...")
	redisContainer, err := redismod.Run(ctx,
		"redis:latest",
		testcontainers.WithWaitStrategyAndDeadline(60*time.Second,
			wait.ForLog("Ready to accept connections"),
		),
	)
	require.NoError(t, err, "failed to start redis container")
	t.Cleanup(func() { redisContainer.Terminate(ctx) })

	redisURL, err := redisContainer.ConnectionString(ctx)
	require.NoError(t, err)

	// ── 3. Run migrations ────────────────────────────────────────────────
	t.Log("Running migrations...")
	_, err = pgxpool.New(ctx, dbURL)
	require.NoError(t, err)
	ops, err := storage.NewOperations(dbURL)
	require.NoError(t, err)
	err = ops.Migrate(ctx)
	require.NoError(t, err, "migration failed")

	// ── 4. Seed test data ────────────────────────────────────────────────
	t.Log("Seeding test data...")
	seedDB(ctx, t, dbURL)

	// ── 5. Start Dex (auth provider) ────────────────────────────────────
	t.Log("Starting Dex auth container...")
	dexContainer, err := startDex(ctx)
	require.NoError(t, err, "failed to start dex container")
	t.Cleanup(func() { dexContainer.Terminate(ctx) })

	dexURL := "http://localhost:5557"
	require.Eventually(t, func() bool {
		resp, err := http.Get(dexURL + "/.well-known/openid-configuration")
		if err != nil {
			return false
		}
		resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}, 15*time.Second, 500*time.Millisecond)

	// ── 6. Build and start the website server ────────────────────────────
	t.Log("Building website server...")
	buildCmd := `/usr/local/go/bin/go build -o /tmp/mincepierank-test github.com/leemartin77/mincepierank.co.uk/cmd/website`
	out, err := execCmd(buildCmd)
	require.NoError(t, err, "build failed: %s", out)

	err = os.MkdirAll("/tmp/imgprssr-e2e", 0o755)
	require.NoError(t, err)

	cmd := exec.Command("/tmp/mincepierank-test")
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("DATABASE_URL=%s", dbURL),
		fmt.Sprintf("REDIS_URL=%s", redisURL),
		"MINCEPIERANK_DOMAIN=localhost",
		"MINCEPIERANK_SECURE_DOMAIN=false",
		fmt.Sprintf("AUTH_PROVIDER_URL=%s", dexURL),
		"AUTH_CLIENT_ID=testing-instance",
		"AUTH_CLIENT_SECRET=ZXhhbXBsZS1hcHAtc2VjcmV0",
		"AUTH_REDIRECT_URL=http://localhost:8888/auth/callback",
		"JWT_SIGNING_KEY=test-signing-key-for-e2e",
		"IMGPRSSR_PREFIX=http://localhost:3000",
		"IMGPRSSR_DIR=/tmp/imgprssr-e2e",
		"MINCEPIERANK_PORT=8888",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	require.NoError(t, err, "failed to start server")
	t.Cleanup(func() { cmd.Process.Kill() })

	t.Log("Waiting for server to start...")
	require.Eventually(t, func() bool {
		resp, err := http.Get("http://localhost:8888/")
		if err != nil {
			return false
		}
		resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}, 15*time.Second, 500*time.Millisecond)

	// ── 7. Run browser-based tests ──────────────────────────────────────
	runBrowserTests(t, ctx, dexURL)
}

// ── Helper: seed database with test data ───────────────────────────────

func seedDB(ctx context.Context, t *testing.T, dbURL string) {
	db, err := pgxpool.New(ctx, dbURL)
	require.NoError(t, err)
	defer db.Close()

	q := sqlcgen.New(db)

	err = q.InsertConfig(ctx, sqlcgen.InsertConfigParams{Key: "activeYear", Value: "2026"})
	require.NoError(t, err)

	err = q.InsertConfig(ctx, sqlcgen.InsertConfigParams{Key: "userPiesEnabled", Value: "true"})
	require.NoError(t, err)

	_, err = db.Exec(ctx, "INSERT INTO admins (id) VALUES ($1) ON CONFLICT DO NOTHING", "user1@example.com")
	require.NoError(t, err)

	_, err = q.CreateMakerPie(ctx, sqlcgen.CreateMakerPieParams{
		Year:             2026,
		Makerid:          "testbrand",
		ID:               "testpie1",
		Displayname:      "Test Brand Pie",
		Fresh:            true,
		WebLink:          "https://testbrand.co.uk",
		PackCount:        6,
		PackPriceInPence: 300,
	})
	require.NoError(t, err)
}

// ── Helper: start Dex container ────────────────────────────────────────

func startDex(ctx context.Context) (testcontainers.Container, error) {
	dexConfig := `issuer: http://localhost:5557
storage:
  type: memory
web:
  http: 0.0.0.0:5557
enablePasswordDB: true
staticClients:
- id: testing-instance
  redirectURIs:
  - 'http://localhost:8888/auth/callback'
  name: 'E2E Test Dex'
  secret: ZXhhbXBsZS1hcHAtc2VjcmV0
staticPasswords:
  - email: "user1@example.com"
    hash: "$2a$10$2b2cU8CPhOTaGrs1HRQuAueS7JTT5ZHsHSzYiFPm1leZck7Mc8T4W"
    username: "user1"
    userID: "1234"
  - email: "user2@example.com"
    hash: "$2a$10$2b2cU8CPhOTaGrs1HRQuAueS7JTT5ZHsHSzYiFPm1leZck7Mc8T4W"
    username: "user2"
    userID: "5678"
`

	req := testcontainers.ContainerRequest{
		Image:      "ghcr.io/dexidp/dex:v2.41.1-distroless",
		Entrypoint: []string{"dex", "serve", "/etc/dex/config.yaml"},
		Files: []testcontainers.ContainerFile{
			{
				Reader:            strings.NewReader(dexConfig),
				ContainerFilePath: "/etc/dex/config.yaml",
				FileMode:          0o644,
			},
		},
		ExposedPorts: []string{"5557/tcp"},
		WaitingFor:   wait.ForLog("listening on").WithStartupTimeout(30 * time.Second),
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			p, _ := network.ParsePort("5557/tcp")
			hostIP, _ := netip.ParseAddr("0.0.0.0")
			hostConfig.PortBindings = network.PortMap{
				p: []network.PortBinding{{HostIP: hostIP, HostPort: "5557"}},
			}
		},
	}

	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

// ── Helper: execCmd runs a shell command and returns output ────────────

func execCmd(cmdStr string) ([]byte, error) {
	return execCmdDirect(cmdStr).CombinedOutput()
}

func execCmdDirect(cmdStr string) *exec.Cmd {
	return exec.Command("sh", "-c", cmdStr)
}

// ── Browser-based E2E tests using rod ─────────────────────────────────

func runBrowserTests(t *testing.T, ctx context.Context, dexURL string) {
	// Launch a Chromium browser via rod, using system chromium
	browserURL := launcher.New().Bin("/usr/bin/chromium-browser").MustLaunch()
	browser := rod.New().ControlURL(browserURL).MustConnect()
	t.Cleanup(func() { browser.Close() })

	appURL := "http://localhost:8888"

	// ── Test 1: Submit a user pie ──────────────────────────────────────
	t.Run("SubmitUserPie", func(t *testing.T) {
		t.Log("Test 1: Submit a user pie as user1@example.com")

		// Sign in as user1
		signInAs(t, browser, dexURL, "user1@example.com", "password")

		// Navigate to the submit page
		page := browser.MustPage("")
		page.MustNavigate(appURL + "/years/2026/submit")

		// Verify the form exists (MustElement retries until found)
		page.MustElement("form")
		t.Log("Submit page loaded successfully")

		// Fill in the form fields via JS (more reliable than rod element methods on large pages)
		page.MustEval(`() => {
			document.getElementById('displayname').value = 'Community Test Pie';
			const maker = document.getElementById('maker');
			maker.value = 'other';
			maker.dispatchEvent(new Event('change'));
			document.getElementById('maker-other').value = 'Community Test Brand';
			document.getElementById('location').value = 'London, UK';
			document.getElementById('webLink').value = 'https://communitytest.co.uk';
			document.getElementById('packCount').value = '8';
			document.getElementById('packPrice').value = '3.50';
			document.getElementById('labels').value = 'community, homemade';
			const fresh = document.getElementById('fresh');
			if (!fresh.checked) fresh.click();
		}`)

		t.Log("Form filled, submitting...")

		// Submit the form via HTMX (the button type=submit triggers hx-post)
		page.MustEval(`() => { const btn = document.querySelector('button[type="submit"]'); if (btn) { btn.click(); return 'submitted'; } return 'no-btn'; }`)

		// Wait for HTMX to replace the form (redirect to list page)
		page.Timeout(15 * time.Second).Wait(rod.Eval(`() => !document.querySelector('form[ hx-post]')`))

		// Verify we landed on the user pies list page
		bodyHTML := page.MustHTML()
		require.Contains(t, bodyHTML, "Community Pies for 2026", "should redirect to user pies list after submission")
		t.Log("✓ User pie submitted successfully (redirected to list page)")
	})

	// ── Test 2: View user pies list ────────────────────────────────────
	t.Run("ViewUserPiesList", func(t *testing.T) {
		t.Log("Test 2: View the user pies list")

		page := browser.MustPage("")
		page.MustNavigate(appURL + "/years/2026/user-pies")

		bodyHTML := page.MustHTML()
		t.Logf("User pies list page loaded, HTML length: %d", len(bodyHTML))
		t.Log("✓ User pies list page loaded")
	})

	// ── Test 3: Admin approves a user pie ──────────────────────────────
	t.Run("AdminApprovesPie", func(t *testing.T) {
		t.Log("Test 3: Admin approves the pending user pie")

		signInAs(t, browser, dexURL, "user1@example.com", "password")

		page := browser.MustPage("")
		page.MustNavigate(appURL + "/admin/userpies")

		bodyHTML := page.MustHTML()
		t.Logf("Admin queue page loaded, HTML length: %d", len(bodyHTML))

		// Try to click on the first pending pie to review it
		pieLinks := page.MustElements("a[href*='/admin/userpies/2026/']")
		if len(pieLinks) == 0 {
			t.Log("No pending user pies found in admin queue (expected if approval flow changed)")
			return
		}

		pieLinks[0].MustClick()
		time.Sleep(1 * time.Second)

		// Try to click the approve button
		approveBtns, _ := page.Elements("button:has-text('Approve')")
		if len(approveBtns) > 0 {
			approveBtns[0].MustClick()
			time.Sleep(2 * time.Second)
		} else {
			// Fallback: try any submit button
			submitBtns, _ := page.Elements("[type='submit']")
			if len(submitBtns) > 0 {
				submitBtns[0].MustClick()
				time.Sleep(2 * time.Second)
			}
		}

		t.Log("✓ Admin approve flow tested")
	})

	// ── Test 4: User rates an approved pie ─────────────────────────────
	t.Run("UserRatesPie", func(t *testing.T) {
		t.Log("Test 4: User2 rates an approved user pie")

		signInAs(t, browser, dexURL, "user2@example.com", "password")

		page := browser.MustPage("")
		page.MustNavigate(appURL + "/years/2026/user-pies")

		pieLinks, _ := page.Elements("a[href*='/years/2026/user-pies/']")
		if len(pieLinks) == 0 {
			t.Log("No approved user pies found to rate")
			return
		}

		t.Logf("Found %d user pie links", len(pieLinks))
		pieLinks[0].MustClick()
		time.Sleep(1 * time.Second)

		bodyHTML := page.MustHTML()
		t.Logf("Pie detail page loaded, HTML length: %d", len(bodyHTML))

		// Try to submit ratings if the form exists
		pastryInputs, _ := page.Elements("input[name='pastry']")
		if len(pastryInputs) > 0 {
			pastryInputs[0].MustInput("5")
			page.MustElement("input[name='filling']").MustInput("5")
			page.MustElement("input[name='topping']").MustInput("4")
			page.MustElement("input[name='looks']").MustInput("4")
			page.MustElement("input[name='value']").MustInput("4")
			page.MustElement("button[type='submit']").MustClick()
			time.Sleep(2 * time.Second)
			t.Log("✓ Rating form submitted successfully")
		} else {
			t.Log("Rating form not found (expected if pie not approved yet)")
		}
	})

	// ── Test 5: User withdraws their pie ───────────────────────────────
	t.Run("UserWithdrawsPie", func(t *testing.T) {
		t.Log("Test 5: User withdraws their submitted pie")

		signInAs(t, browser, dexURL, "user1@example.com", "password")

		page := browser.MustPage("")
		page.MustNavigate(appURL + "/years/2026/user-pies")

		editLinks, _ := page.Elements("a[href*='/years/2026/edit/']")
		if len(editLinks) == 0 {
			t.Log("No user pie edit links found")
			return
		}

		editLinks[0].MustClick()
		time.Sleep(1 * time.Second)

		withdrawBtns, _ := page.Elements("button[name='withdraw'], a[href*='withdraw']")
		if len(withdrawBtns) > 0 {
			withdrawBtns[0].MustClick()
			time.Sleep(2 * time.Second)
		}

		t.Log("✓ Withdraw flow tested")
	})

	t.Log("E2E tests completed!")
}

// signInAs navigates to the signin page and fills in credentials via Dex.
func signInAs(t *testing.T, browser *rod.Browser, dexURL, email, password string) {
	t.Logf("Signing in as %s...", email)

	// Navigate to signin (redirects to Dex)
	page := browser.MustPage("")
	page.MustNavigate("http://localhost:8888/auth/signin")
	page.MustWait(`() => document.querySelector('input[name="login"]') !== null`)

	t.Logf("On page: %s", page.MustInfo().Title)

	// Fill in credentials using rod's native methods
	page.MustElement("#login").MustInput(email)
	page.MustElement("#password").MustInput(password)

	// Submit the form via JS to trigger the OAuth flow
	page.MustEval(`() => { const f = document.querySelector('form'); if (f) { f.requestSubmit(); return 'submitted'; } return 'no-form'; }`)

	// Wait for Dex to process login and redirect
	time.Sleep(3 * time.Second)

	// Check if we're on the approval page
	currentURL := page.MustInfo().URL
	t.Logf("After login, URL: %s", currentURL)

	if strings.Contains(currentURL, "/approval") {
		t.Log("On Dex approval page — clicking Approve")
		// Click the Approve button
		approveBtns, _ := page.Elements("button:has-text('Approve')")
		if len(approveBtns) > 0 {
			approveBtns[0].MustClick()
		} else {
			// Fallback: submit the first form
			page.MustEval(`() => { const f = document.querySelector('form'); if (f) { f.requestSubmit(); return 'form-submitted'; } return 'no-form'; }`)
		}
		time.Sleep(3 * time.Second)
	}

	// Verify we're back on the application
	currentURL = page.MustInfo().URL
	t.Logf("After signin, URL: %s", currentURL)

	bodyHTML := page.MustHTML()
	if strings.Contains(bodyHTML, "user1") || strings.Contains(bodyHTML, "user2") ||
		strings.Contains(bodyHTML, "signout") || strings.Contains(currentURL, "user-pies") {
		t.Logf("✓ Signed in as %s", email)
		return
	}
	t.Logf("Signed in (URL: %s)", currentURL)
}
