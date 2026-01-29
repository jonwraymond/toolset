# User Journey: Curated Toolset Exposure

This walkthrough shows how to build a safe MCP-only toolset and export it for use
by an MCP server or other consumer.

## Scenario

You have a large tool registry but want to expose **only safe, read-only tools**
from the `mcp` namespace. You also want a single export step that can later be
adapted to OpenAI or Anthropic formats.

## Step 1: Define canonical tools

```go
import "github.com/jonwraymond/tooladapter"

all := []*tooladapter.CanonicalTool{
    {
        Namespace:   "mcp",
        Name:        "search",
        Tags:        []string{"read", "safe"},
        Description: "Search resources",
        InputSchema: &tooladapter.JSONSchema{Type: "object"},
    },
    {
        Namespace:   "mcp",
        Name:        "execute",
        Tags:        []string{"write", "danger"},
        Description: "Execute a command",
        InputSchema: &tooladapter.JSONSchema{Type: "object"},
    },
}
```

## Step 2: Build a filtered toolset

```go
import "github.com/jonwraymond/toolset"

safe, err := toolset.NewBuilder("mcp-safe").
    FromTools(all).
    WithNamespace("mcp").
    WithTags([]string{"safe"}).
    ExcludeTools([]string{"mcp:execute"}).
    Build()
```

At this point, `safe` contains only the `mcp:search` tool.

## Step 3: Export for MCP usage

```go
import "github.com/jonwraymond/tooladapter/adapters"

exposure := toolset.NewExposure(safe, adapters.NewMCPAdapter())
exports, warnings, errs := exposure.ExportWithWarnings()
if len(errs) > 0 {
    // handle conversion errors (tool IDs included)
}
```

If the tools used unsupported schema features, you would see warnings here. If a
tool fails conversion, it will be omitted from `exports` and reported in `errs`.

## Flow Diagram

```mermaid
flowchart LR
  A[Canonical Tools] --> B[Builder]
  B --> C{Filters}
  C --> D[Policy]
  D --> E[Toolset]
  E --> F[Exposure via Adapter]
  F --> G[MCP / OpenAI / Anthropic]
```

## Notes

- Filters are AND-composed.
- Policies run last and can hard-deny tools.
- Exposure never mutates the underlying Toolset.
