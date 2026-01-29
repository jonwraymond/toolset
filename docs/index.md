# toolset

Composable tool collection library for building curated tool surfaces.

## Overview

toolset provides a thread-safe `Toolset` type and a fluent `Builder` for composing
curated tool collections from multiple sources. It is pure data composition:
**no execution, no I/O, no network**. It integrates with `tooladapter` for
multi-format export (MCP/OpenAI/Anthropic), and can optionally ingest tools from
`toolindex` via adapters.

## Design Goals

1. **Deterministic output**: tool listing and exposure are stable across runs.
2. **Non-destructive filtering**: filters never mutate the original toolset.
3. **Composability**: small, reusable filters and policies.
4. **Protocol-agnostic exposure**: export via `tooladapter` adapters.
5. **Minimal dependencies**: no runtime execution or transport coupling.

## Position in the Stack

```
toolmodel --> tooladapter --> toolset --> metatools-mcp
```

- `tooladapter` normalizes tools into canonical form.
- `toolset` builds curated collections and exports to different formats.
- `metatools-mcp` may later consume toolset to expose filtered subsets.

## Core Types

| Type | Purpose |
|------|---------|
| `Toolset` | Thread-safe collection of canonical tools |
| `Builder` | Fluent builder for filters/policies |
| `FilterFunc` | Reusable filter predicate |
| `Policy` | Hard allow/deny decisions |
| `Exposure` | Export to MCP/OpenAI/Anthropic via adapters |

## Quick Start

```go
package main

import (
    "github.com/jonwraymond/toolset"
    "github.com/jonwraymond/tooladapter"
)

func main() {
    ts := toolset.New("safe-core")

    ts.Add(&tooladapter.CanonicalTool{
        Namespace: "mcp",
        Name:      "search",
        Tags:      []string{"read", "safe"},
        InputSchema: &tooladapter.JSONSchema{Type: "object"},
    })

    ts.Add(&tooladapter.CanonicalTool{
        Namespace: "mcp",
        Name:      "exec",
        Tags:      []string{"write", "danger"},
        InputSchema: &tooladapter.JSONSchema{Type: "object"},
    })

    safe := ts.Filter(func(t *tooladapter.CanonicalTool) bool {
        return toolset.TagsAll([]string{"safe"})(t)
    })

    _ = safe
}
```

## Builder Example

```go
builder := toolset.NewBuilder("mcp-safe")

safeSet, err := builder.
    FromTools(allTools).
    WithNamespace("mcp").
    WithTags([]string{"safe"}).
    ExcludeTools([]string{"mcp:exec"}).
    Build()
```

## Exposure Example (via tooladapter)

```go
exposure := toolset.NewExposure(safeSet, adapters.NewMCPAdapter())
exports, warnings := exposure.ExportWithWarnings()
```

## Versioning

toolset follows semantic versioning aligned with the stack. The source of truth
is `ai-tools-stack/go.mod`, and `VERSIONS.md` is synchronized across repos.

See `VERSIONS.md` for current versions.

## Next Steps

- [Design Notes](design-notes.md)
- [User Journey](user-journey.md)
