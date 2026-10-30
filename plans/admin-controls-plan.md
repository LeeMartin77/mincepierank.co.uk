# Admin area gap-closing plan

---

## Current state (verified)

| Entity | Read | Admin write | Notes |
|---|---|---|---|
| `config` | yes | **full CRUD** (`/admin/config`) | includes `activeYear`, `userPiesEnabled` |
| `maker_pie_yearly` | yes | **full CRUD + bulk zip** (`/admin/makerpies/*`) | image upload, category assign |
| `user_pie_yearly` | yes | **moderation** (`/admin/userpies/*`) | approve/reject/dup/promote/delete |
| `maker` (brands) | `GetMakersForYear`, `GetMaker` only | **none** | generated `MakerCreate/Update/Delete` in `internal/storage/generated/generate.maker.gen.go` exist but are never called |
| `categories` | yes | **assign only** | seeded by one-off migration (`migrations.go:116`); no create/rename/delete |
| `maker_pie_ranking_yearly` / `user_pie_ranking_yearly` | per-user only | **none** | no delete query exists for maker-pie ratings; `notes` column shown nowhere |

---

## Phase 1 — Brand (maker) management ★ do first — **IMPLEMENTED**

Implemented as above with these notes: OpenAPI-first routes in
`api/website.yaml` + regen; ids validated path-safe
(`^[a-z0-9][a-z0-9-]{0,63}$`) and auto-derived from the name using the
migration's slug convention; storage wrappers `CreateMaker`/`UpdateMaker`/
`DeleteMaker`/`GetAllMakers`/`GetMakerPieCount` on `OperationWrapper`
(`ON CONFLICT DO NOTHING` inserts mean existence is checked first, returning
`ErrBrandAlreadyExists`/`ErrBrandNotFound`); delete refused while pies
reference the brand; promote requires an existing brand (the userpie review
form's action URLs were also broken — missing `{year}/{ownerid}` — and are
fixed); pie create rejects unknown `makerid`; shared `saveUploadedFile`
helper extracted. Templates: `admin/brands.gohtml`, `admin/brand.gohtml`,
admin-index link, brand `<select>` on pies + promote. Tests: storage CRUD
(testcontainers), id derivation, handler suite with fake store, template render.

The biggest gap: brands are read-only app-wide; every fix below prevents
silently-orphaned pies (brand pages `INNER JOIN maker`, see `makers.go:22`).

### 1.1 Storage
- Add a `makers.go`-level set of operations on `internal/storage/core.go`
  wrapping the already-generated `generated.MakerCreate/MakerUpdate/MakerDelete`
  (currently dead code), plus `GetAllMakers` (ordered by name) for the admin list,
  and `GetMakerIdsWithPies` if the picker needs it.
- `logo` is a plain text column — treat it as an image path under
  `ImgprssrDir`, same convention as pie images (`/years/{y}/{brand}/{pie}/{file}`).
  Plan to store brand logos at `/brands/{id}/{filename}`.

### 1.2 API spec + codegen
- Add to `api/website.yaml` (admin tag):
  - `GET /admin/brands` — list all brands (+ pie counts)
  - `GET /admin/brands/new` + `POST /admin/brands` — create form/submit
  - `GET /admin/brands/{id}` + `PUT` + `DELETE` — edit/delete
- `task generate` regenerates `internal/website/generated/website-server.gen.go`.

### 1.3 Handlers
- New `internal/website/admin_brands.go` following the `admin.go` pattern
  (`validateAdmin` guard, `templater.PageData` with `IsAdminPage: true`).
- Logo upload reuses the pie-image upload pattern (`c.FormFile`, `os.MkdirAll`,
  write under `ImgprssrDir`) — factor a small shared helper since the pattern
  is repeated three times already (`admin.go:269-307`, `admin.go:355-395`,
  `admin_bulk.go:88-153`).
- On delete: refuse (or warn) if pies still reference the brand id.

### 1.4 Templates
- `assets/templates/html/pages/admin/brands.gohtml` (list + create form) and
  `brand.gohtml` (edit), linked from `admin/index.gohtml`.
- Replace free-text `makerid` in `admin/pies.gohtml:17` (create form) and the
  promote form (`admin_user_pies.go` / `userpie.gohtml`) with a `<select>` of
  existing brands.
- Server-side: reject create/promote with an `makerid` that has no `maker` row.

## Phase 2 — Category management

- Storage: `GetCategoryById`, `CreateCategory(slug,label)`, `UpdateCategory`,
  `DeleteCategory` (delete guarded or cascading `maker_pie_categories`; check
  pie counts first as with brands).
- Routes `GET/POST/PUT/DELETE /admin/categories*`; templates `categories.gohtml`
  + `category.gohtml`.
- Fix the **silent-drop** bug: `SetPieCategoriesTx` (`maker_pies.go:544`) drops
  unknown slugs without error. Return an error listing unknown slugs so the UI
  can surface it (and so the promote form stops quietly discarding categories).
- Slug generation from label (reuse migration convention: lowercase, spaces→`-`).

## Phase 3 — Ratings visibility + moderation

- sqlc queries (new `-- name:` entries in `maker_pie_ranking.sql`,
  `user_pie_ranking.sql`):
  - `ListMakerPieRankings` / `ListUserPieRankings` (by pie, paginated, incl. `userid` + `notes`)
  - `DeleteMakerPieRanking` (does not exist at all today — users can't un-rate either)
  - `DeleteUserPieRanking` (per-rating; `DeleteUserPieRankingsByUser` already exists)
- Admin routes: `GET /admin/rankings?...` (pie picker + rating list),
  `DELETE /admin/rankings/{pie}/{userid}`.
- Decide `notes` fate (currently collected, displayed nowhere):
  - **recommended:** admin-visible only on the review screen + delete control,
    keep off public pages; or
  - surface publicly with moderation controls.
- Out of scope: rating-rate limits/anti-spam heuristics — note for later.

## Phase 4 — Season lifecycle + CSRF + cleanups

- **CSRF:** no middleware anywhere (verified — no `csrf` hits in `internal/`).
  Add per-session token (session already in redis, see `auth.go` + `UserIdMiddleware`)
  and validate on all admin POST/PUT/DELETE. Public rating PUT + submit forms
  benefit too.
- **Spec/handler drift:** `api/website.yaml` declares `imageFile`/`webLink`/
  `packCount` form fields but handlers read `img`/`web_link`/`pack_count` —
  align spec to reality.
- Admin index counts (admin.go:54-60) and dashboard could fold in Phase 3/4 data.

---

## Sequencing rationale

1 → 2 → 3 → 4. Phase 2's slug fix is low-risk and unblocks the category
checkbox on promote; Phase 1 is first because every other phase's forms either
reference brands (promote) or reuse its upload helper; CSRF is last because it
touches every form and is easier once the final set exists.

## Testing

- Follow existing patterns: `internal/website/user_pies_e2e_test.go`
  (playwright-style) for happy paths, `*_test.go` handler tests for validation
  rejections.
- `cmd/template-preview` (`task preview`) for the new admin templates' markup
  pass before wiring storage.
- Dev data: `data/2024/bulk/bulk.zip` per README for a populated stack.
