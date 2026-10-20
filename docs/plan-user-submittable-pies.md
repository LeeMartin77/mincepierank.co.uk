# Plan: User Submittable Pies (submission + moderation + promotion)

**Branch:** `feat/user-submittable-pies`
**Status:** Plan for review — all design questions resolved, ready to build.

---

## 1. Where this feature went

The feature is **not** new — it shipped in the SvelteKit era and was deleted wholesale in the Go rewrite.

| When | What |
|---|---|
| `a659172` (2023-09-23) | **"added: Simply implemented user pie submission"** — the original feature: `src/routes/create/`, `src/lib/storage/userPies.ts`, `userPieRankings.ts`, Cassandra tables `user_pie_yearly` / `user_pie_ranking_yearly` |
| `92225fc` / `f997e25` | Feature flag `customPiesEnabled`, then partially reverted |
| `c67e474` / `67994ee` | Owner edit + delete |
| `0799ea2` | Admin "claim anonymous rankings" (only admin tool that ever touched user content) |
| `06fece4` (2024-04-08) | **"init: Hello world docker container"** — the SvelteKit app is deleted, user pie feature dies with it |
| `c0bd1bf` / `23900da` | Cassandra→Postgres port copies the *tables* and the *data* (the extractor still lists `user_pie_yearly` / `user_pie_ranking_yearly`) |
| `8a074cd` / `32741a1` / `673b76e` | Go structs, generated CRUD and sqlc models are generated for the user pie tables |
| `6c60c18` (2024-07-20) | Go menu: "My Custom Pies" link gets commented out (sign-in work) |
| `dfdc0d1` (2026-08-28) | Bootstrap rewrite — the commented link is deleted outright |

### What survived vs. what was lost

**Survived (dead code, no callers):**
- `user_pie_yearly` + `user_pie_ranking_yearly` tables — `internal/storage/sqlc/schema/base.sql:49-80`, `internal/storage/migrations.go:59-89`
- `types.UserPieYearly` / `types.UserPieRankingYearly` — `internal/storage/types/tables.go:69-101`
- Generated CRUD — `internal/storage/generated/generate.userpieyearly.gen.go`, `generate.userpierankingyearly.gen.go`
- sqlc models — `internal/storage/sqlcgen/models.go:73-99`
- Cassandra→Postgres data extraction — `scripts/cassandra-extractor/index.js:21-29`

**Lost:**
- Every route, handler, template and query. `internal/storage/sqlc/queries/` has **no** `user_pies.sql` and `sqlcgen.Querier` has **no** user-pie methods.

**Never existed in any era — genuinely new work:**
- Any moderation, approval, dedupe or promotion. `git log --all -S "moderat" / "approv" / "dedupe" / "promote"` returns **zero** commits.
- The only vestigial hook is `user_pie_yearly.clean`, hard-coded `false` at submit time in the old code and never read by anything, ever. It was clearly the original author's intended moderation flag.

**Conclusion:** the data model survived; the feature did not. We rebuild it natively in Go, and add the moderation pipeline that never existed.

---

## 2. Decisions taken (from review Q&A)

| Question | Decision |
|---|---|
| Approval outcome | **Two paths:** approve as a *user pie* (published, own page, own rankings) **or** promote to a *brand pie* (enters the main rankings) |
| Dedupe | **Suggest candidates, admin decides.** No silent merging. |
| Sign-in to submit | **Required.** `owner_userid` = email (same convention as the lost feature; matches `user_pie_yearly.owner_userid`) |
| Rankings | **Both.** User pies are rankable on their own page via `user_pie_ranking_yearly`; brand pies keep `maker_pie_ranking_yearly` |
| Submitter identity | **Hidden from non-admins.** Never rendered on public pages. Admins see it on the review screen. (Matches the lost feature, which stripped `owner_userid` from the public payload.) |
| Brand selection | **Known-brand dropdown + free-text fallback.** Picking a known brand stores `makerid`; choosing "other" falls back to the free-text `maker` field. |
| Feature flag | **Disabled by default**, enabled via a `config` key (`userPiesEnabled`) using the same pattern as `activeYear` / `readonly`. |
| Approved user pies | **Separate "Community Pies" list** at `/years/{year}/user-pies`. Never mixed into the main "All Pies" list or the brand leaderboards. |
| After promotion | **The user pie page stays live** with a prominent banner linking to the brand pie it became. |
| Owner editing | **Only before approval** (`pending`). Once approved, the admin owns the record. |
| Rejection feedback | **A single `moderation_note`** shown to the owner. No comment thread. |
| After rejection | **A fresh submission is required.** `rejected` is terminal — there is no path back to `pending`. |
| Community list filtering | **Plain pagination only.** No brand filter. |
| State machine | **Lives in Go, not in the database.** The DB is a data store; all transition rules, guards and permissions are Go code, unit-testable without Postgres. See §3.3. |
| Settled states | **`promoted`, `duplicate` and `rejected` are terminal.** To remove a submission in one of those states, the admin deletes it — there are no further status transitions. |

---

## 3. Target design

### 3.1 Terminology

| Concept | Table | Notes |
|---|---|---|
| Brand | `maker` | id (slug), name, logo, website |
| Brand pie | `maker_pie_yearly` | admin-curated, `validated` flag, `oid` uuid, categories via `maker_pie_categories` |
| **Submission / user pie** | `user_pie_yearly` | **new**: `status`, `oid`, moderation audit columns |
| User pie ranking | `user_pie_ranking_yearly` | `(year, pieid, userid)` — note **no** `makerid` |

### 3.2 State machine — overview

```
  submit ──────► pending ──────► approved ──────► promoted
                    │              │                  │
                    │              │                  └─► brand pie created;
                    │              │                      rankings migrated;
                    │              │                      user pie page stays
                    │              │                      live with a banner
                    │              └─► published user pie: public page,
                    │                  own ranking form
                    │
                    ├──── admin marks duplicate ────► duplicate
                    │
                    └──── admin rejects (note) ────► rejected

  promoted, duplicate and rejected are TERMINAL: to remove a submission
  that has reached one, the admin deletes it — there are no further
  status transitions.
```

- `pending` — submitted, awaiting review. Visible to owner + admins only.
- `approved` — published as a user pie. Public page, publicly rankable, appears in the Community Pies list.
- `promoted` — became a brand pie. Page stays live with a banner + link; its own ranking form is disabled (ratings migrated to the brand pie). **Terminal** — to remove it, delete it.
- `duplicate` — collapsed into another submission. Hidden publicly; `duplicate_of` points at the survivor. **Terminal** — to remove it, delete it.
- `rejected` — removed by an admin. Hidden publicly, soft-kept for audit; a hard delete is also available. **Terminal** — there is no `rejected → pending` path; the owner makes a fresh submission.

**The transition table in §3.3 is the authoritative definition.** The diagram above is illustrative.

### 3.3 State machine — the transitions live in Go, not in the database

> **Design constraint (agreed in review):** the database is a data store. It holds rows; it does not decide what may happen to a submission. All status-transition rules, guards and permissions are expressed in Go, where they can be unit-tested without Postgres.

**New file: `internal/storage/user_pie_status.go`** (package `storage`) — deliberately free of any `pgx` import, so the whole state machine is testable as a pure unit:

```go
type UserPieStatus string

const (
    UserPieStatusPending   UserPieStatus = "pending"
    UserPieStatusApproved  UserPieStatus = "approved"
    UserPieStatusPromoted  UserPieStatus = "promoted"
    UserPieStatusDuplicate UserPieStatus = "duplicate"
    UserPieStatusRejected  UserPieStatus = "rejected"
)

// The transition table. Anything not listed here is refused.
//
// promoted, duplicate and rejected are terminal: to remove a submission
// that has reached one of them, the admin deletes it.
var userPieTransitions = map[UserPieStatus][]UserPieStatus{
    UserPieStatusPending: {
        UserPieStatusApproved, UserPieStatusPromoted,
        UserPieStatusDuplicate, UserPieStatusRejected,
    },
    UserPieStatusApproved: {
        UserPieStatusPromoted, UserPieStatusDuplicate, UserPieStatusRejected,
    },
    UserPieStatusPromoted:  {}, // terminal - removal is a delete
    UserPieStatusDuplicate: {}, // terminal - removal is a delete
    UserPieStatusRejected:  {}, // terminal - a fresh submission is required
}

func (s UserPieStatus) CanTransitionTo(next UserPieStatus) bool
func (s UserPieStatus) IsPublic() bool        // approved, promoted  -> publicly visible
func (s UserPieStatus) IsOwnerEditable() bool // pending only       -> owner may edit/withdraw
func (s UserPieStatus) IsRankable() bool      // approved only      -> accepts user rankings
func (s UserPieStatus) IsTerminal() bool      // promoted, duplicate, rejected
```

**Transition table**

| From | To | Actor | Side effects |
|---|---|---|---|
| — | `pending` | owner, on submit | row created, `submitted_at` set |
| `pending` | `approved` | admin | becomes publicly visible + rankable |
| `pending` | `promoted` | admin | brand pie created, image copied, rankings migrated, `promoted_maker_pie_oid` set |
| `pending` | `duplicate` | admin | `duplicate_of` = survivor oid |
| `pending` | `rejected` | admin | `moderation_note` set; **terminal** |
| `approved` | `promoted` | admin | as above |
| `approved` | `duplicate` | admin | `duplicate_of` = survivor oid |
| `approved` | `rejected` | admin | `moderation_note` set |
| `promoted` | — | — | **terminal.** Removal is a delete: `DELETE /admin/userpies/{oid}`; the brand pie it became is removed separately via `/admin/makerpies/{oid}` |
| `duplicate` | — | — | **terminal.** Removal is a delete: `DELETE /admin/userpies/{oid}`. There is no "un-mark as duplicate" |
| `rejected` | — | — | **terminal.** A fresh submission is required |

**Terminal states and removal.**

`promoted`, `duplicate` and `rejected` are **terminal** — they have no outgoing transitions. To remove a submission that has reached one of them, the admin **deletes it**:

| State | How it is removed |
|---|---|
| `pending` / `approved` | Either `→ rejected` (a moderation decision, with a note to the owner) **or** a hard delete. The admin chooses. |
| `promoted` | Hard delete of the submission, plus deletion of the brand pie it became via the existing `/admin/makerpies/{oid}` route. The review page offers both together. |
| `duplicate` | Hard delete. There is no "un-mark as duplicate" — the surviving submission is the one that matters. |
| `rejected` | Hard delete, once the admin no longer wants the audit record. |

This keeps the moderation vocabulary (`approve` / `promote` / `mark duplicate` / `reject`) distinct from housekeeping (`delete`), and it means the `status` column only ever records a decision an admin actually made.

Deliberately **not** included: an `approved → pending` "send back for changes" path, or any outgoing transition from `promoted`, `duplicate` or `rejected`. Rejection plus a fresh submission covers the need for changes, and keeping `pending` as the only owner-editable state keeps the ownership rule simple.

**How the rules are enforced at write time**

`internal/storage/user_pies.go` exposes one guarded write path:

```go
// SetUserPieStatus moves a submission between states.
// It loads the current status, asks the state machine whether the move is
// legal, and then does a compare-and-set UPDATE so a concurrent change is
// caught rather than silently overwritten.
func (o *OperationWrapper) SetUserPieStatus(
    ctx context.Context, oid uuid.UUID, next UserPieStatus,
    moderatedBy string, note string,
) error
```

The `UPDATE ... WHERE oid = $1 AND status = $2` compare-and-set is the only place the old status appears in SQL, and it is a concurrency guard, not a rule engine. All the *deciding* is done in Go before the statement is built.

**What the database is (and is not) responsible for**

| Concern | Where it lives |
|---|---|
| Which transitions are legal | Go — `userPieTransitions` |
| Who may perform a transition (admin vs owner) | Go — the handler, using `IsOwnerEditable` / `IsPublic` |
| Whether a pie may be ranked | Go — `IsRankable` |
| Whether a pie is publicly visible | Go — `IsPublic`, applied in the public queries |
| Storing the current status | Postgres — one `text` column |
| Preventing two admins racing on the same submission | Postgres — compare-and-set on `status` |
| Keeping the column within its domain | Postgres — optional `CHECK (status IN (...))`, a value constraint only |

A `CHECK` constraint on the allowed *values* is the one piece of SQL-level enforcement worth having, and it encodes no transition rules. Everything else is Go.

**Testing**

- `internal/storage/user_pie_status_test.go` — pure table-driven unit test of the transition table and the predicates. No Postgres, no testcontainers, runs in milliseconds.
- `internal/storage/user_pies_test.go` — testcontainers test covering the guarded write path, the compare-and-set race, and the promotion transaction.

### 3.4 Schema changes

`user_pie_yearly` gains:

```sql
ALTER TABLE user_pie_yearly ADD COLUMN oid uuid DEFAULT gen_random_uuid();
ALTER TABLE user_pie_yearly ADD COLUMN status text NOT NULL DEFAULT 'pending';
ALTER TABLE user_pie_yearly ADD COLUMN submitted_at timestamp NOT NULL DEFAULT NOW();
ALTER TABLE user_pie_yearly ADD COLUMN moderated_by text;         -- admin email
ALTER TABLE user_pie_yearly ADD COLUMN moderated_at timestamp;
ALTER TABLE user_pie_yearly ADD COLUMN moderation_note text;      -- reason, shown to owner
ALTER TABLE user_pie_yearly ADD COLUMN duplicate_of uuid;          -- oid of survivor
ALTER TABLE user_pie_yearly ADD COLUMN promoted_maker_pie_oid uuid;-- oid of brand pie
ALTER TABLE user_pie_yearly ADD COLUMN makerid text;               -- optional brand slug

CREATE UNIQUE INDEX user_pie_yearly_oid_idx  ON user_pie_yearly (oid);
CREATE INDEX        user_pie_yearly_status_idx ON user_pie_yearly (status, year);
CREATE INDEX        user_pie_yearly_owner_idx  ON user_pie_yearly (year, owner_userid);
CREATE INDEX        user_pie_yearly_maker_idx  ON user_pie_yearly (year, makerid);
CREATE INDEX        user_pie_ranking_pie_idx   ON user_pie_ranking_yearly (year, pieid);
```

Notes:
- `oid` mirrors the existing `maker_pie_yearly.oid` pattern → clean admin URLs (`/admin/userpies/{oid}`) and a safe target for `duplicate_of` / `promoted_maker_pie_oid`.
- `clean` is left in place (the Cassandra extractor still round-trips the table) but is **no longer used** — `status` supersedes it. We won't drop it in this change.
- `makerid` is optional: the submitter picks a brand from the known `maker` list when they can; `maker` (free text) remains the fallback display value.

### 3.5 Slug / identity

- `user_pie_yearly.id` = server-generated slug from `displayname` (kebab-case), **uniquified across the whole year** (not per owner) so `user_pie_ranking_yearly.pieid` is unambiguous within a year.
- The slug is generated once at creation and **not** regenerated on edit (so URLs stay stable).
- `oid` is the surrogate used by admin URLs and cross-references.

---

## 4. Routes

### 4.1 Public / user-facing

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/years/{year}/submit` | Submission form (sign-in required) |
| `POST` | `/years/{year}/submit` | Create submission (multipart, image optional) |
| `GET` | `/years/{year}/user-pies` | **"Community Pies"** — public list of approved user pies for the year. Kept entirely separate from "All Pies" and the brand leaderboards. |
| `GET` | `/years/{year}/user-pies/{id}` | User pie page (public if `approved`/`promoted`; owner+admin if `pending`) |
| `PUT` | `/years/{year}/user-pies/{id}/rating` | Rate an approved user pie (mirrors the brand-pie rating route) |
| `GET` | `/profile/pies/{year}` | "My Pies" — the signed-in user's submissions for a year |
| `GET` | `/profile/pies/{year}/{id}` | Edit form (owner, `pending` only) |
| `PUT` | `/profile/pies/{year}/{id}` | Save edit (owner, `pending` only) |
| `DELETE` | `/profile/pies/{year}/{id}` | Withdraw submission (owner, `pending` only) |

### 4.2 Admin

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/admin/userpies` | **The queue.** Filter by status / year / free-text search |
| `GET` | `/admin/userpies/duplicates` | Duplicate candidate groups |
| `GET` | `/admin/userpies/{oid}` | Review one submission + all actions |
| `POST` | `/admin/userpies/{oid}/approve` | Approve as a user pie (publish) |
| `POST` | `/admin/userpies/{oid}/promote` | Promote to a brand pie |
| `POST` | `/admin/userpies/{oid}/reject` | Remove (soft), with a note |
| `POST` | `/admin/userpies/{oid}/duplicate` | Mark as duplicate of `{survivor_oid}` |
| `DELETE` | `/admin/userpies/{oid}` | Hard delete |

---

## 5. Admin workflows in detail

### 5.1 Check submissions (the queue)

`GET /admin/userpies` renders:

- Status tabs / filter chips: **Pending** (default) · Approved · Promoted · Duplicate · Rejected · All
- Year filter (defaults to the active year) and a free-text search over displayname / maker / web_link / owner email
- Table: thumbnail, display name, brand, location, submitter, submitted-at, status badge, **Review** link
- Pending count badge, so the admin knows when there's work waiting

`GET /admin/userpies/{oid}` shows the full submission (large image, every field, submitter, timestamps) plus:

- **Possible existing brand pies** — matches against `maker_pie_yearly` (same brand + similar name, or same web link). This is the cross-queue dedupe: catching "someone submitted Greggs Classic Mince Pie" when we already stock it.
- **Possible duplicate submissions** — other submissions that look like this one.
- Actions: Approve / Promote / Mark duplicate / Reject / Delete / Edit fields.
  - **Approve, Promote, Mark duplicate and Reject** are only offered while the submission is live (`pending` or `approved`) — they are status transitions, and the state machine refuses them from a terminal state.
  - **Delete** is always offered. It is the only way out of `promoted`, `duplicate` and `rejected`.
  - For a `promoted` submission the delete action also offers to remove the brand pie it became.
  - The buttons are driven by the Go predicates, so the UI cannot offer an action the state machine would reject.

### 5.2 Dedupe

**Detection — "suggest candidates, admin decides".** Two independent signals, both computed in SQL without needing any Postgres extension:

1. **Exact link match** — two submissions share a non-empty `web_link` (normalised: lowercase, strip trailing `/`).
2. **Normalised name + brand** — `lower(regexp_replace(displayname,'[^a-z0-9]+','','g'))` plus the brand key (`makerid`, else `maker`), grouped; groups of 2+ are candidates.
   - A looser variant groups on the normalised name alone (different brand) — shown as a weaker signal.

Output is a set of candidate groups; the admin picks a survivor per group. Marking a submission a duplicate sets `status='duplicate'` and `duplicate_of=<survivor oid>`. The survivor stays editable by the admin, so the best image and details can be consolidated onto it. `duplicate` is **terminal** (§3.3): the marked submission is neither editable nor transitionable, so once the admin is happy with the survivor it is **hard-deleted**. There is no "un-mark as duplicate".

No automatic merging, ever.

### 5.3 Promote to a brand pie

`POST /admin/userpies/{oid}/promote` form fields:

| Field | Notes |
|---|---|
| `makerid` | Dropdown of `maker` rows + "add new brand" (name/logo/website) |
| `pieid` | Slug for the brand pie; pre-filled from the submission's slug; **collision-checked** against `maker_pie_yearly (year, makerid, pieid)` |
| `displayname` | Pre-filled from the submission |
| `fresh`, `web_link`, `pack_count`, `pack_price_in_pence` | Pre-filled from the submission |
| `categories[]` | Multi-select, same widget as `/admin/makerpies/{oid}` |
| `validated` | Checkbox, default ticked (admin-created pies are validated by definition) |

Executed **in a single transaction**:

1. Load the submission; must be `pending` or `approved`.
2. Resolve/validate `makerid` against the `maker` table.
3. **Conflict check** — if `(year, makerid, pieid)` already exists, abort with a link to the existing pie (this is the dedupe guard on the brand-pie side).
4. `INSERT INTO maker_pie_yearly (...)` with `validated = true` → returns `oid`.
5. `SetPieCategoriesTx` for the new pie.
6. **Copy the image** from `{ImgprssrDir}/{year}/user-pies/{id}/{file}` → `{ImgprssrDir}/{year}/{makerid}/{pieid}/{file}`; the brand pie's `image_file` becomes `/{year}/{makerid}/{pieid}/{file}` (matches the existing brand-pie path convention).
7. **Migrate rankings**: `INSERT INTO maker_pie_ranking_yearly SELECT ... FROM user_pie_ranking_yearly WHERE year=$1 AND pieid=$2 ON CONFLICT (year,makerid,pieid,userid) DO UPDATE ...` — brand-pie ratings win on conflict.
8. Update the submission: `status='promoted'`, `promoted_maker_pie_oid=<new oid>`, `moderated_by/at`.
9. Commit → redirect to `/admin/makerpies/{oid}`.

### 5.4 Remove / reject

Two different things, deliberately kept distinct:

- **Reject** — a *moderation decision*, available only while a submission is still live (`pending` or `approved`). Sets `status='rejected'`, `moderated_by/at`, `moderation_note`. The row stays for audit; hidden from all public queries.
  - `rejected` is **terminal** (§3.3). The owner is shown the `moderation_note` and must make a **fresh submission** — there is no "un-reject" and no edit path back to `pending`.
- **Delete** — *removal*, and the only way out of a settled state. `DELETE FROM user_pie_yearly WHERE oid=$1` plus best-effort removal of the image directory. Also deletes the user's own rankings for that pie.
  - Available from any status.
  - For a `promoted` submission this removes the user pie page; removing the **brand pie** it became is a separate action via the existing `/admin/makerpies/{oid}` delete. The review page offers both together as "Delete this submission and the brand pie it became".

---

## 6. User-facing feature detail

### 6.1 Submission form

`GET /years/{year}/submit` — year must be the active year (the same rule the brand-pie rating form uses: `eq .ActiveYear .Pie.Year`).

Fields (mirroring the old lost feature, plus a couple of improvements):

| Field | Required | Notes |
|---|---|---|
| `displayname` | ✅ | "What's it called?" |
| `makerid` | | Brand picker — dropdown of known `maker` rows, **plus an "Other / not listed" option** that reveals the free-text field below |
| `maker` | ✅ (if `makerid` is "other") | "Who made it?" — free text. Always submitted alongside `makerid` so the display name is never blank. If a known brand is picked, `maker` is prefilled with that brand's name. |
| `image` | | optional, ≤ 20 MB (the old limit) |
| `location` | | "Where did you find it?" |
| `labels` | | comma-separated, e.g. `premium,rum,puff pastry` |
| `web_link` | | validated as a URL, blanked if invalid |
| `pack_count` | | integer ≥ 1 |
| `pack_price_in_pence` | | integer ≥ 0 |
| `fresh` | | checkbox |

The brand picker is a `<select>` listing the known brands (from `maker`), with a final option "Other / not listed". Selecting "Other" reveals a free-text input (progressive disclosure — htmx toggle, or plain CSS/JS). If a known brand is selected, `makerid` is stored and `maker` is filled from that brand's name; if "Other" is selected, `makerid` is stored as `NULL` and `maker` holds the free text.

This keeps the promotion flow a dropdown rather than free text for the common case, while still allowing brands we don't stock yet.

Behaviour:
- **Sign-in required** → anonymous visitors get a "Sign in to submit" prompt with a sign-in link (the old feature did exactly this).
- `readonly=true` in config → submissions blocked (matches both the old feature and the current brand-pie rating behaviour).
- `userPiesEnabled` is **false by default** → the route 404s and the menu links are hidden until an admin enables it. See §6.6.
- On success → redirect to `/years/{year}/user-pies/{id}` with a "submitted, awaiting review" notice.

### 6.2 "Community Pies" list

`GET /years/{year}/user-pies` is a **separate list**, deliberately not merged into "All Pies" (`/years/{year}/all-pies`) or the brand pages:

- Lists only `approved` (and `promoted`) user pies for the year, ordered by ranking average.
- Reuses the existing `page:pielist` template + `fragment:components:pielinkcard` so it looks consistent with the rest of the site.
- The card's "From {Brand}" line is omitted for user pies with no known brand (free-text `maker` is shown instead) — the card component already guards on `{{ if .Maker }}`.
- Ranking summary comes from `user_pie_ranking_yearly`, so a community pie with no votes shows 0 / 0 votes.
- Reachable from the main nav as **"Community Pies"** and linked from the "Submit a Pie" page.
- **Plain pagination only** — the existing `page:pielist` "Load More" pattern, with `page`/`limit` query params handled by the existing `getLimitAndIdx` helper. **No brand filter** (and no category filter): a user pie has no fixed brand or category taxonomy, so filtering by them would be misleading.
- Hidden when `userPiesEnabled` is false.

### 6.3 "My Pies"

`GET /profile/pies/{year}` lists the signed-in user's submissions for a year with their status badge. Pending rows get **Edit** / **Withdraw** buttons; approved/promoted rows link to the public page. Rejected rows show the admin's `moderation_note` with a **"Submit a new pie"** link (rejection is terminal — there is nothing to edit).

### 6.4 Public user pie page

`GET /years/{year}/user-pies/{id}`:

- Title, image, brand, location, pack count, price, cost-per-pie, labels
- **No submitter identity is rendered.** The `owner_userid` column is never passed to the template for a non-admin viewer. The old feature explicitly stripped it from the public payload (`const { owner_userid, ...pie } = res.value`); we keep that, and there is no masked-handle fallback. Admins see the submitter email on `/admin/userpies/{oid}` and in the admin queue table.
- Ranking summary + star-rating form, reusing `fragment:components:userranking` and `fragment:components:rankingsummary`
- If `pending`: owner/admin-only, "awaiting review" notice, no rating form
- If `promoted`: **the page stays live** with a prominent banner — "This is now the official *{Brand} {Name}*" — linking to the brand pie at `/years/{year}/brands/{makerid}/{pieid}`. The own-ranking form is disabled (the ratings were migrated to the brand pie) and the brand pie's summary is shown instead.
- If `duplicate`: hidden from the public; the owner sees "this was marked as a duplicate" and a link to the surviving submission. Terminal — it is not editable, and the admin removes it with a delete.
- If `rejected`: hidden from the public; the owner sees the admin's single `moderation_note` and a **"Submit a new pie"** link. They cannot edit or restore this one — rejection is terminal.
- Owner sees Edit/Withdraw (pending only); admins see a "Review" link

### 6.5 Ranking an approved user pie

`PUT /years/{year}/user-pies/{id}/rating` mirrors `RateYearBrandPie`:
- Reuses the existing, already-unit-tested `parseRatings` helper (`internal/website/brand_pie_page.go:55`)
- Upserts into `user_pie_ranking_yearly (year, pieid, userid, ...)`
- Only allowed when `status='approved'` and `year == activeYear` and not `readonly`
- `user_pie_ranking_yearly` is currently `int` columns; ratings are 1–5

### 6.6 The `userPiesEnabled` config key

The feature is **off until an admin turns it on**, using the same `config` key/value mechanism as `activeYear` and `readonly`.

| | |
|---|---|
| Key | `userPiesEnabled` |
| Values | `'true'` / `'false'` |
| Default | **Disabled** — the key is absent from the `config` table until an admin creates it |

- `internal/storage/config.go` gains `GetUserPiesEnabled(c) (bool, error)`, mirroring `GetReadonly` (`internal/storage/config.go:26-40`): read `config.userPiesEnabled`, `strconv.ParseBool` the value, and **fail closed** (return `false`) if the key is missing or unparseable.
  - Note the difference from `GetActiveYear` (`config.go:11-24`), which *errors* on a missing key. A missing `userPiesEnabled` must **not** be an error — it means "not enabled yet".
- The admin enables it from the existing **Admin → Config** page (`/admin/config`), which already does generic key/value create/update/delete via htmx. Adding the key is a two-field form entry: key `userPiesEnabled`, value `true`. No new admin UI needed for the flag itself.
- When disabled:
  - `GET/POST /years/{year}/submit` → 404
  - `/years/{year}/user-pies` and `/years/{year}/user-pies/{id}` → 404 (there is nothing public to show)
  - `/profile/pies/{year}` → 404
  - The "Submit a Pie" / "My Pies" menu links are hidden
  - **The admin queue at `/admin/userpies` stays reachable** — an admin may want to review and triage what was submitted while the flag was on, and existing submissions must not become unreachable.
- `readonly=true` continues to block submissions independently of this flag.

---

## 7. Implementation notes / gotchas

1. **`fragment:components:userranking` hardcodes the brand rating URL** (`hx-put="/years/{{.Pie.Year}}/brands/{{.Pie.MakerId}}/{{.Pie.Id}}/rating"`). It needs a `RatingUrl` key passed in the `dict` so the same fragment can serve user pies. Small, safe change to one shared fragment + one call site.
2. **The `userranking` fragment's `eq .ActiveYear .Pie.Year` guard** works for user pies unchanged.
3. **Image writes.** The existing admin code writes with `fs.ModeAppend` (`internal/website/admin.go:270`, `:355`, `admin_bulk.go:149`) — that's a permission bit, not a write mode, and it appends to an existing file. New code will use `0o644` + `os.Create`. I will **not** silently change the existing admin paths in this feature (separate concern), but it's worth a follow-up ticket.
4. **sqlc vs. raw SQL.** The codebase mixes both: sqlc for simple CRUD, raw SQL in `internal/storage/*.go` for joins/aggregates (see `GetMakerPieYearlyRankingSummary` in `core.go:60`). I'll follow that: sqlc for the new `user_pies.sql` / `user_pie_ranking.sql` CRUD, raw SQL for summaries, the admin queue query, the duplicate-candidate query and the ranking migration.
5. **Generation.** `task generate` runs `sqlc generate` + `go generate ./...` (storage generator) + `oapi-codegen` for the API. `api/website.yaml` must be updated **before** the handlers, because `oapi-codegen` generates the `ServerInterface` that the handlers must satisfy.
6. **`user_pie_yearly.id` uniqueness.** The generator's `UserPieYearlyCreate` uses `ON CONFLICT DO NOTHING` on `(year, id, owner_userid)`. Because we uniquify slugs across the whole year, that PK is effectively `(year, id)` for our purposes — but the *generated* struct's PK stays as-is.
7. **Struct tags.** The storage generator reads `col:"..."` tags, so new columns must be added to `types.UserPieYearly` in `internal/storage/types/tables.go` and the generated file regenerated, or the CRUD will silently skip them.
8. **Anonymity vs. ownership.** `auth.go` sets `c.Keys["userid"]` to the email when signed in, else an anonymous UUID cookie. Submissions must check `signedin == true` and use the email as `owner_userid`; anonymous users must not be able to submit.
9. **One way in to `status`.** Because the rules live in Go, they are only enforced if every write goes through the guarded path. **No handler, query or script may `UPDATE user_pie_yearly SET status = ...` directly** — everything goes through `SetUserPieStatus`, which consults the transition table first. This is the one convention to watch in review.
10. **Public queries filter on status in Go, not in SQL.** `IsPublic()` decides visibility, but the query still needs to narrow the result set. The narrowing is a `WHERE status IN (...)` with the list derived from the Go predicate, so the two cannot drift. Same for `IsRankable` on the ranking write path.
11. **Terminal states are removed, not transitioned.** `promoted`, `duplicate` and `rejected` have no outgoing transitions, so the admin "removes" a settled submission with `DELETE /admin/userpies/{oid}` rather than a status change. The delete route must therefore accept a submission in **any** status, and for a `promoted` submission it should also offer to delete the brand pie it became.

---

## 8. Work breakdown

### Phase 0 — Data layer + state machine
- [ ] `internal/storage/user_pie_status.go` — **new**: the `UserPieStatus` type, the transition table and the predicates. **No `pgx` import** — this file is the testable heart of the feature (§3.3)
- [ ] `internal/storage/user_pie_status_test.go` — **new**: pure table-driven unit test of the transition table and predicates. No Postgres, no testcontainers
- [ ] `internal/storage/migrations.go` — the `ALTER TABLE` / `CREATE INDEX` statements
- [ ] `internal/storage/sqlc/schema/base.sql` — mirror the new DDL
- [ ] `internal/storage/types/tables.go` — extend `UserPieYearly` with the new columns
- [ ] `internal/storage/sqlc/queries/user_pies.sql` — **new**: CRUD + admin queue + slug uniqueness
- [ ] `internal/storage/sqlc/queries/user_pie_ranking.sql` — **new**: rating CRUD + summary
- [ ] `internal/storage/user_pies.go` — **new**: `Operations` methods (queue, candidates, summaries, promotion, ranking migration, guarded `SetUserPieStatus`)
- [ ] `internal/storage/core.go` — extend the `Operations` interface
- [ ] `task generate` — regenerate sqlc + storage generator output
- [ ] `internal/storage/user_pies_test.go` — testcontainers test for the guarded write path and the promotion transaction

### Phase 1 — Submission (public)
- [ ] `api/website.yaml` — submit / user-pies / rating / profile routes
- [ ] `internal/website/user_pies.go` — **new**: submit, list, detail, ranking handlers
- [ ] `internal/website/user_pie_profile.go` — **new**: "My Pies" + edit/withdraw
- [ ] `assets/templates/html/pages/submit.gohtml` — **new**
- [ ] `assets/templates/html/pages/userpie.gohtml` — **new**
- [ ] `assets/templates/html/pages/userpielist.gohtml` — **new**
- [ ] `assets/templates/html/pages/profilepies.gohtml` — **new**
- [ ] `assets/templates/html/fragments/components/userranking.gohtml` — parameterise the rating URL
- [ ] Image upload handling + slug generation + validation
- [ ] Brand picker (known-brand dropdown + free-text "other" fallback)

### Phase 2 — Admin queue
- [ ] `api/website.yaml` — admin userpie routes
- [ ] `internal/website/admin_user_pies.go` — **new**: queue, review, approve, reject, delete
- [ ] `assets/templates/html/pages/admin/userpies.gohtml` — **new**: the queue
- [ ] `assets/templates/html/pages/admin/userpie.gohtml` — **new**: review + actions
- [ ] `assets/templates/html/pages/admin/index.gohtml` — add a "Submissions" link
- [ ] All admin transitions go through `SetUserPieStatus` — no handler writes `status` directly

### Phase 3 — Dedupe
- [ ] Duplicate-candidate query (link match + normalised name+brand + normalised name)
- [ ] `GET /admin/userpies/duplicates` page
- [ ] "Possible existing brand pies" lookup on the review page
- [ ] `POST /admin/userpies/{oid}/duplicate` action

### Phase 4 — Promotion
- [ ] `POST /admin/userpies/{oid}/promote` handler (transactional)
- [ ] Promotion form in the review page (brand dropdown, slug with collision check, categories, validated)
- [ ] Image copy `user-pies/{id}` → `{makerid}/{pieid}`
- [ ] Ranking migration `user_pie_ranking_yearly` → `maker_pie_ranking_yearly`
- [ ] Promoted banner + brand-pie link on the user pie page

### Phase 5 — Polish
- [ ] Menu links ("Community Pies", "Submit a Pie", "My Pies") — the last two signed-in only, all hidden when `userPiesEnabled` is false
- [ ] `internal/storage/config.go` — `GetUserPiesEnabled`, **fails closed** when the key is absent
- [ ] `userPiesEnabled` config key + gating on all user-pie routes
- [ ] `cmd/template-preview` scenarios for the new templates
- [ ] Tests for the new handlers (`internal/website/user_pies_test.go`)
- [ ] README / docs update

---

## 9. Decision log

### 9.1 Resolved (review feedback folded in)

| # | Question | Decision |
|---|---|---|
| 1 | Submitter identity on public pages | **Hidden from non-admins.** Never rendered publicly, no masked-handle fallback. Admins only. |
| 2 | Brand selection | **Known-brand dropdown + free-text "Other" fallback.** Known brand → `makerid` set; "Other" → `makerid` NULL + free-text `maker`. |
| 3 | Feature flag default | **Disabled by default.** `userPiesEnabled` config key, same mechanism as `activeYear` / `readonly`. |
| 4 | Where approved user pies appear | **Separate "Community Pies" list** at `/years/{year}/user-pies`. Not in "All Pies", not in brand leaderboards. |
| 5 | After promotion | **Page stays live** with a prominent banner linking to the brand pie. |
| 6 | Owner editing after approval | **No.** Owner edits are only possible while `pending`. Once approved, the admin owns the record. |
| 7 | Moderation feedback | **A single `moderation_note`** shown to the owner. No comment thread. |
| 8 | After rejection | **A fresh submission is required.** `rejected` is terminal — no path back to `pending`. |
| 9 | Community list filtering | **Plain pagination only.** No brand filter, no category filter. |
| 10 | Where the state machine lives | **In Go, not in the database.** The DB is a data store; all transition rules, guards and permissions are Go code, unit-testable without Postgres. See §3.3. |
| 11 | After promotion or duplicate marking | **Terminal.** To remove a submission in either state, the admin deletes it — and for a promoted pie, deletes the brand pie it became via the existing admin route. |
| 12 | Legacy rows carried over from Cassandra | **No special handling.** They arrive as `pending` and are triaged in the queue like any other submission. |

### 9.2 Deliberate omissions

Things that were considered and are **not** in scope, so nobody re-litigates them later:

- **No `approved → pending` "send back for changes".** Rejection plus a fresh submission covers the need, and keeping `pending` as the only owner-editable state keeps the ownership rule simple.
- **No comment thread on a submission.** A single `moderation_note` is enough for a site this size.
- **No brand or category filter on the community list.** A user pie has no fixed brand or category taxonomy, so filtering by them would be misleading.
- **No soft-delete of a promoted brand pie.** `promoted` is terminal. Removing a promoted pie means deleting the brand pie through the existing `/admin/makerpies/{oid}` route, not changing the submission's status.
- **No Postgres `pg_trgm` dependency.** Duplicate detection uses plain normalisation (lowercase, strip non-alphanumerics) and exact link matching — no extension required, so no superuser `CREATE EXTENSION` in production.
- **No automatic dedupe merging.** Suggestions only; the admin always picks the survivor.
- **No outgoing transitions from `promoted`, `duplicate` or `rejected`.** All three are terminal; removal is a delete. A mis-marked duplicate is deleted and the owner resubmits.

---

## 10. Risks

| Risk | Mitigation |
|---|---|
| Image upload abuse (disk fill) | 20 MB cap (as the old feature), content-type sniffing, and the existing `imgprssr` path convention |
| Slug collisions across owners | Slug uniquifier runs across the whole year, not per owner |
| Ranking data lost on promotion | Ranking migration runs inside the same transaction as the brand-pie insert |
| Promotion creates a duplicate brand pie | Collision check on `(year, makerid, pieid)` + the "possible existing brand pies" lookup on the review page |
| Migration breaks the Cassandra extractor | New columns are additive; the extractor copies whole tables, so it round-trips unchanged |
| `user_pie_yearly` rows carried over from Cassandra have `status = 'pending'` (the column default) | They arrive in the admin queue as ordinary pending submissions and are triaged by hand like any other. No special handling needed. |
| Feature flag is off while submissions already exist | The admin queue stays reachable regardless of the flag, so nothing is stranded. Public pages 404 while off, which is the intended behaviour. |
| State machine rules get bypassed by a stray SQL write | The rules live in Go, so the DB cannot enforce them. Mitigation is convention + review: `SetUserPieStatus` is the only sanctioned writer of `status` (§7, note 9). A `CHECK` constraint on the *values* catches typos, not rule violations. |
| A rule change needs a DB migration | By design it does not — the transition table is a Go map, changed in a normal PR with a unit test. |

---

## 11. Where the lost code still lives (for reference)

```bash
# The original SvelteKit implementation, still in git history
git show 06fece4^:src/routes/create/+page.server.ts
git show 06fece4^:src/routes/create/+page.svelte
git show 06fece4^:src/lib/storage/userPies.ts
git show 06fece4^:src/lib/storage/userPieRankings.ts
git show 06fece4^:src/lib/storage/types.ts
git show 06fece4^:src/routes/years/[year]/user-pies/[id]/+page.server.ts
git show 06fece4^:src/routes/years/[year]/user-pies/[id]/+page.svelte
git show 06fece4^:src/components/UserPieRankingInterface.svelte
git show 06fece4^:src/routes/profile/pies/[year]/+page.server.ts
git show 06fece4^:src/lib/claimAnonymous.ts
```
