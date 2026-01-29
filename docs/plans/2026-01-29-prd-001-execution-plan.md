# PRD-001 Execution Plan — toolset (Strict TDD)

**Status:** Ready
**Date:** 2026-01-29

This plan is a strict TDD execution guide for PRD-001. Each task must follow:
1) Red (write failing test)
2) Red verification (run test)
3) Green (minimal implementation)
4) Green verification (run test)
5) Commit (one commit per task)

---

## Task 0 — Repo scaffolding

**Goal:** Establish Go module and baseline docs without behavior.

- Create `go.mod` (Go 1.24+) with dependency on `tooladapter` and `toolindex`.
- Add `doc.go` package documentation.
- Add minimal `README.md` with scope + usage note.

**Commit:** `chore(toolset): initialize module skeleton`

---

## Task 1 — Toolset core

**Goal:** Implement Toolset type and basic operations.

**Tests:**
- New toolset has name and empty list.
- Add/Get/Remove.
- Count, IDs, Tools.
- Filter returns a new toolset with subset.

**Implementation:**
- `toolset.go` with Toolset struct and methods.
- Thread-safe storage with RWMutex.

**Commit:** `feat(toolset): add Toolset core`

---

## Task 2 — Builder pattern

**Goal:** Fluent builder for constructing toolsets.

**Tests:**
- Build from registry.
- WithNamespace(s), WithTags, WithTools, ExcludeTools.
- Chained filters combine as expected.

**Implementation:**
- `builder.go` with Builder, registry interface, and Build().

**Commit:** `feat(toolset): add Builder`

---

## Task 3 — Filter helpers

**Goal:** Provide reusable filter functions.

**Tests:**
- Namespace filter.
- Tags filter (any/none/all semantics as documented).
- Category filter.
- ID allow/deny filters.

**Implementation:**
- `filter.go` with FilterFunc type and helper constructors.

**Commit:** `feat(toolset): add filter helpers`

---

## Task 4 — Policy interface

**Goal:** Access control policies for tool inclusion.

**Tests:**
- Policy allow/deny behavior.
- Builder applies policy and excludes denied tools.

**Implementation:**
- `policy.go` with Policy interface and default allow policy.

**Commit:** `feat(toolset): add policy support`

---

## Task 5 — Exposure helpers

**Goal:** Export toolsets to target protocol formats.

**Tests:**
- Exposure via MCP/OpenAI/Anthropic adapters.
- Feature loss warnings surfaced.

**Implementation:**
- `exposure.go` with exposure helper types and adapter wiring.

**Commit:** `feat(toolset): add exposure helpers`

---

## Task 6 — Documentation & quality gates

- `README.md` with usage examples.
- Run `go test ./...`.
- Run `golangci-lint run` if configured.

**Commit:** `docs(toolset): add usage docs`

---

## Verification Checklist

- [ ] `go test ./...` passes
- [ ] coverage > 80%
- [ ] no lints (if configured)
- [ ] README + doc.go updated

---

## Commit Order

1) chore(toolset): initialize module skeleton
2) feat(toolset): add Toolset core
3) feat(toolset): add Builder
4) feat(toolset): add filter helpers
5) feat(toolset): add policy support
6) feat(toolset): add exposure helpers
7) docs(toolset): add usage docs
