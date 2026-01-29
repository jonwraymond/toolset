// Package toolset provides composable tool collection building.
//
// Toolset enables curated, filtered, and access-controlled tool surfaces
// from multiple sources. It is pure data composition with no I/O, execution,
// or network dependencies.
//
// Core concepts:
//   - Toolset: thread-safe collection of canonical tools
//   - Builder: fluent API for constructing toolsets
//   - FilterFunc: predicates for filtering tools
//   - Policy: access control decisions
//   - Exposure: export to MCP/OpenAI/Anthropic via tooladapter
package toolset
