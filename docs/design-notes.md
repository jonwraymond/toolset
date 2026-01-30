# Design Notes

This document captures design decisions, filter semantics, and integration
constraints for the toolset library.

## Toolset Core

### Thread safety

`Toolset` must be safe for concurrent reads/writes:
- Internally uses `sync.RWMutex`.
- `Add`, `Remove` take write locks.
- `Get`, `Tools`, `IDs`, `Count` use read locks.

### Deterministic ordering

`Tools()` and `IDs()` must return tools in stable order:
- Sort lexicographically by canonical tool ID (`namespace:name`).
- Determinism matters for:
  - pagination in downstream layers
  - stable exposure output
  - reproducible tests

## Filtering Semantics

Filters are predicates applied to canonical tools. In the Builder, filters are
**AND-composed** in the order added. The recommended filter helpers:

- **Namespace filter**: include only tools with namespace in allowed set.
- **TagsAny**: tool matches if it contains *any* tag in the set.
- **TagsAll**: tool matches only if it contains *all* tags in the set.
- **Category filter**: exact match on `Category` field.
- **Allow IDs / Deny IDs**: explicit allow/deny lists for tool IDs.

### Policy order

Policies apply **after** all filters. This guarantees that a policy decision can
hard-deny any tool even if it passed filters.

## Policy Interface

A policy is a simple allow/deny decision:

```go
type Policy interface {
    Allow(tool *tooladapter.CanonicalTool) bool
}
```

Use cases:
- Restrict tools by namespace or tag
- Enforce security categories
- Restrict by `RequiredScopes`

### Contract

- **Concurrency:** policies must be safe for concurrent use after construction.
- **Nil handling:** `Allow(nil)` must return `false`.
- **Determinism:** results must be stable for identical inputs.
- **Ownership:** implementations must not mutate tool data.

## Registry Interface

Registries provide the tool source for builders:

```go
type Registry interface {
    Tools() []*tooladapter.CanonicalTool
}
```

### Contract

- **Concurrency:** `Tools()` may be called concurrently.
- **Ownership:** returned slice is caller-owned; tools are read-only snapshots.
- **Determinism:** ordering should be stable for identical registry state.
- **Nil handling:** returning `nil` is treated as empty.

## Exposure Semantics

Exposure uses `tooladapter.Adapter` to export toolsets:

- Conversion is **read-only**; no mutation to canonical tools.
- Feature loss warnings are aggregated from the adapter.
- Conversion errors are surfaced via `ExportWithWarnings()` as a `[]error`.
- Exposure returns `[]any` for protocol-specific tool shapes.

## Integration with toolindex

`toolset` can optionally ingest tools from `toolindex` by:
- Listing tools from index
- Converting `toolmodel.Tool` to `tooladapter.CanonicalTool` via MCP adapter

This keeps toolset independent of execution (`toolrun`) or docs (`tooldocs`).

## Non-goals

- Runtime execution or transport wiring
- Persistence or caching
- Multi-tenant enforcement

## Error Strategy

- Builder returns errors for missing sources or invalid inputs.
- Filters and policies are pure functions and should not panic.
- `ExportWithWarnings()` returns conversion errors (tool ID included in the error).

## Performance

- Filters are O(n) over tool count.
- Repeated `Tools()` calls should avoid re-sorting if possible (optional cache).
