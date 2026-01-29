# toolset

Composable tool collection library for building curated, filtered, and access-controlled tool surfaces.

## Overview

`toolset` sits between `tooladapter` (protocol normalization) and `metatools-mcp` (server exposure). It provides:

- **Toolset**: Thread-safe collection of canonical tools with filtering
- **Builder**: Fluent API for constructing toolsets from registries
- **Filters**: Reusable predicates (namespace, tags, categories, IDs)
- **Policies**: Access control layer applied after filtering
- **Exposure**: Export to protocol-specific formats via adapters

## Installation

```bash
go get github.com/jonwraymond/toolset
```

## Quick Start

```go
package main

import (
    "github.com/jonwraymond/toolset"
    "github.com/jonwraymond/tooladapter"
)

func main() {
    // Build a toolset from tools
    ts, err := toolset.NewBuilder("my-tools").
        FromTools(allTools).
        WithNamespace("github").
        WithTags([]string{"read-only"}).
        ExcludeTools([]string{"github:delete-repo"}).
        Build()
    if err != nil {
        panic(err)
    }

    // Filter further
    filtered := ts.Filter(toolset.CategoryFilter("search"))

    // Export to MCP format
    exposure := toolset.NewExposure(filtered, mcpAdapter)
    mcpTools, warnings := exposure.ExportWithWarnings()
}
```

## Stack Position

```
toolmodel → tooladapter → toolset → metatools-mcp
```

## What This Is

- Pure data composition library
- Thread-safe tool collections
- Declarative filtering and policies

## What This Is NOT

- Tool execution (see `toolrun`)
- Network I/O or external calls
- Schema validation (see `toolmodel`)

## License

See repository root.
