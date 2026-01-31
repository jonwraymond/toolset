# toolset

> **DEPRECATED**: This package has been migrated to [`github.com/jonwraymond/toolcompose/set`](https://github.com/jonwraymond/toolcompose).
>
> Please update your imports. See [MIGRATION.md](./MIGRATION.md) for details.

---

## Migration

This repository is no longer maintained. All functionality has been moved to the `set` package within `toolcompose`, which provides a unified approach to tool composition.

### New Location

```bash
go get github.com/jonwraymond/toolcompose
```

Then import:

```go
import "github.com/jonwraymond/toolcompose/set"
```

### Why the Change?

The `toolcompose` module consolidates tool composition functionality into a single, cohesive package:

- **Unified API**: Tool composition, filtering, and policies in one module
- **Better maintainability**: Single repository for related functionality
- **Cleaner dependencies**: Reduced module fragmentation

### Timeline

- **Now**: This repository accepts no new features
- **Future**: This repository will be archived

For migration instructions, see [MIGRATION.md](./MIGRATION.md).
