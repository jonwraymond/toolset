# Migration Guide: toolset to toolcompose/set

This guide covers migrating from `github.com/jonwraymond/toolset` to `github.com/jonwraymond/toolcompose/set`.

## Import Path Changes

| Old Import | New Import |
|------------|------------|
| `github.com/jonwraymond/toolset` | `github.com/jonwraymond/toolcompose/set` |

## Step-by-Step Migration

### 1. Update go.mod

Remove the old dependency and add the new one:

```bash
go get github.com/jonwraymond/toolcompose
go mod tidy
```

### 2. Update Import Statements

Find and replace all imports in your codebase:

```go
// Before
import "github.com/jonwraymond/toolset"

// After
import "github.com/jonwraymond/toolcompose/set"
```

### 3. Update Package References

Update all package references from `toolset` to `set`:

```go
// Before
ts, err := toolset.NewBuilder("my-tools").
    FromTools(allTools).
    WithNamespace("github").
    Build()

filtered := ts.Filter(toolset.CategoryFilter("search"))
exposure := toolset.NewExposure(filtered, mcpAdapter)

// After
ts, err := set.NewBuilder("my-tools").
    FromTools(allTools).
    WithNamespace("github").
    Build()

filtered := ts.Filter(set.CategoryFilter("search"))
exposure := set.NewExposure(filtered, mcpAdapter)
```

## API Compatibility

The `toolcompose/set` package maintains API compatibility with `toolset`. The following types and functions are available:

### Types

- `set.Toolset` (was `toolset.Toolset`)
- `set.Builder` (was `toolset.Builder`)
- `set.Exposure` (was `toolset.Exposure`)
- `set.Filter` (was `toolset.Filter`)
- `set.Policy` (was `toolset.Policy`)

### Functions

- `set.NewBuilder()` (was `toolset.NewBuilder()`)
- `set.NewExposure()` (was `toolset.NewExposure()`)
- `set.CategoryFilter()` (was `toolset.CategoryFilter()`)
- `set.NamespaceFilter()` (was `toolset.NamespaceFilter()`)
- `set.TagFilter()` (was `toolset.TagFilter()`)
- `set.IDFilter()` (was `toolset.IDFilter()`)

## Automated Migration

For large codebases, use `sed` or `gofmt` to automate the migration:

```bash
# Find all Go files with the old import
grep -r "github.com/jonwraymond/toolset" --include="*.go" .

# Replace imports (macOS)
find . -name "*.go" -exec sed -i '' 's|github.com/jonwraymond/toolset|github.com/jonwraymond/toolcompose/set|g' {} +

# Replace imports (Linux)
find . -name "*.go" -exec sed -i 's|github.com/jonwraymond/toolset|github.com/jonwraymond/toolcompose/set|g' {} +

# Replace package references
find . -name "*.go" -exec sed -i '' 's|toolset\.|set.|g' {} +
```

After running the automated migration, verify with:

```bash
go build ./...
go test ./...
```

## Questions?

If you encounter issues during migration, please open an issue in the [toolcompose repository](https://github.com/jonwraymond/toolcompose/issues).
