# template-preview

Renders the site's HTML templates with fixed sample data, so you can look at
(and screenshot) UI changes without bringing up Postgres, Redis, Dex or
imgprssr.

It depends only on the standard library and `internal/templatefuncs`, and
reads templates straight from `assets/templates/html`, so it works anywhere a
Go toolchain does — including CI or sandboxes with no Docker and no module proxy.

## Usage

Run from the repo root.

### Serve (live)

```
task preview
# or
go run ./cmd/template-preview
```

Open <http://127.0.0.1:8099/> for an index of scenarios. Templates are
re-parsed on every request, so edit a `.gohtml` file and refresh.

`/static/` serves `assets/static` (Bootstrap, htmx) and `/img/` serves
`data/images` in place of imgprssr.

### Write static HTML

```
go run ./cmd/template-preview -out tmp/preview
```

Writes one `<scenario>.html` per scenario and exits non-zero if any template
fails to execute — handy as a quick smoke check. The files reference
`/static/...` and `/img/...`, so view them through a server that maps those
paths (or use serve mode).

### Before/after comparison for a PR

Point `-templates` at another checkout to render the same data through the old
templates:

```
git worktree add ../mpr-main origin/main
go run ./cmd/template-preview -addr 127.0.0.1:8098 -templates ../mpr-main/assets/templates/html
go run ./cmd/template-preview -addr 127.0.0.1:8099
```

### Screenshots

Any headless browser works against serve mode, e.g. with Playwright:

```
npx playwright screenshot --full-page --viewport-size=1100,900 \
  http://127.0.0.1:8099/preview/singlepie-ranked ranked.png
```

## Flags

| Flag         | Default                 | Purpose                                         |
|--------------|-------------------------|-------------------------------------------------|
| `-templates` | `assets/templates/html` | Template directory to parse                     |
| `-static`    | `assets/static`         | Served at `/static/`                            |
| `-images`    | `data/images`           | Served at `/img/` (stand-in for imgprssr)       |
| `-addr`      | `127.0.0.1:8099`        | Listen address in serve mode                    |
| `-out`       | *(empty)*               | Write HTML files to this dir and exit           |

## Scenarios

| Name                 | What it shows                                                   |
|----------------------|-----------------------------------------------------------------|
| `singlepie-ranked`   | Signed in, active year, existing ranking — editable form        |
| `singlepie-unranked` | Signed in, active year, not yet ranked — empty form             |
| `singlepie-errors`   | Rating rejected with validation errors                          |
| `singlepie-readonly` | Site in read-only mode — ranking shown, no form                 |
| `singlepie-pastyear` | Pie from a previous year — ranking shown, no form               |

To add one, append to `scenarios` in `main.go`. Page data is a plain
`map[string]interface{}` mirroring what the handler in `internal/website`
passes, so copy the keys from there.

## Limitations

- **Not the real server.** Forms render but submitting them (`hx-put` etc.)
  goes nowhere; use this for layout and visual state, not behaviour.
- **Shared template functions.** Both this tool and `internal/templater` take
  their functions from `internal/templatefuncs`. Put new stdlib-only template
  functions there and both pick them up. `containscat` needs storage types, so
  it lives in the templater and is stubbed (always false) here.
- **Sample data only.** Structs here are minimal stand-ins with the field
  names the templates read; they don't touch `internal/storage`.
