//go:build tools

// Package tools pins build-time code generators as real module dependencies.
//
// Both generators below are only ever run as commands, never imported by
// application code, so `go mod tidy` sees no dependency for them and drops
// them from go.mod. Generation then fails until someone re-fetches a binary
// by hand, with no record of which version produced the checked-in files.
//
//   - oapi-codegen runs via the //go:generate directives in
//     internal/website/run.go.
//   - sqlc runs via the `generate` task in Taskfile.yaml.
//
// The build tag above keeps this file out of normal builds, vet and tests while
// `go mod tidy` still counts the imports, so the versions stay pinned in go.mod.
//
// To upgrade a generator, change its version in go.mod, then run
// `task generate` and commit any regenerated files.
package tools

import (
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
)
