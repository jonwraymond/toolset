# PRD-001: toolset Composition Library Implementation

> **For agents:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Create a composable tool collection library that enables curated, filtered, and access-controlled tool sets from multiple sources.

**Architecture:** Fluent builder pattern for constructing toolsets from registries, with filtering by namespace, tags, categories, and custom policies. Supports multiple exposure formats (MCP, OpenAI, Anthropic) via tooladapter.

**Tech Stack:** Go 1.24+, tooladapter dependency, toolindex dependency

**Priority:** P1

---

## Overview

`toolset` provides composable tool collections for creating curated API surfaces. It enables filtering, access control, and multi-format exposure of tools.

**Reference:** `metatools-mcp/docs/proposals/protocol-agnostic-tools.md`

---

## Scope

### In scope
- Toolset core type with add/get/remove/filter operations.
- Builder pattern for assembling toolsets from registries or explicit lists.
- Filter helpers: namespace, tags, categories, explicit allow/deny lists.
- Policy interface for access control decisions.
- Exposure helpers for exporting to MCP/OpenAI/Anthropic via tooladapter adapters.
- Unit tests for all exported behavior.

### Out of scope (future)
- Persistence, versioning, or live sync.
- Runtime execution or tool invocation.
- Multi-tenant enforcement (handled in later PRD).

---

## Directory Structure

```
toolset/
├── toolset.go
├── toolset_test.go
├── builder.go
├── builder_test.go
├── filter.go
├── filter_test.go
├── policy.go
├── policy_test.go
├── exposure.go
├── exposure_test.go
├── doc.go
├── go.mod
└── go.sum
```

---

## Requirements

### R1 — Toolset operations
- Thread-safe add/get/remove.
- Stable `IDs()` output.
- Filter returns a new toolset without mutating original.

### R2 — Builder
- Build toolsets from:
  - Registry (adapter or index-based)
  - Explicit tool list
- Fluent filters: namespace(s), tags, categories, include/exclude list.

### R3 — Policy
- Policy interface that can deny tools based on tool metadata.
- Builder can apply policy during build.

### R4 — Exposure
- Export to MCP/OpenAI/Anthropic via adapters.
- Return warnings for feature loss.

### R5 — Tests
- TDD: failing tests first.
- Coverage target >80% for the module.

---

## Acceptance Criteria

- Toolset core + builder + filters + policy + exposure implemented.
- MCP/OpenAI/Anthropic exposure works via tooladapter.
- All tests pass with >80% coverage.
- Documentation (`doc.go`, README) present and accurate.

---

## Dependencies

- `github.com/jonwraymond/tooladapter` (protocol adapters)
- `github.com/jonwraymond/toolindex` (optional registry integration)

---

## Notes

- Avoid UTCP terminology. MCP terminology only.
- Do not leak adapter-specific types outside exposure functions.
