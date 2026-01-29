package toolset

import (
	"sort"
	"sync"

	"github.com/jonwraymond/tooladapter"
)

// FilterFunc is a predicate for filtering tools.
type FilterFunc func(*tooladapter.CanonicalTool) bool

// Toolset is a thread-safe collection of canonical tools.
type Toolset struct {
	name  string
	mu    sync.RWMutex
	tools map[string]*tooladapter.CanonicalTool // keyed by ID()
}

// New creates a new Toolset with the given name.
func New(name string) *Toolset {
	return &Toolset{
		name:  name,
		tools: make(map[string]*tooladapter.CanonicalTool),
	}
}

// Name returns the toolset's name.
func (ts *Toolset) Name() string { return ts.name }

// Add adds a tool. Nil tools are silently ignored.
func (ts *Toolset) Add(tool *tooladapter.CanonicalTool) {
	if tool == nil {
		return
	}
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.tools[tool.ID()] = tool
}

// Get retrieves a tool by ID. Returns (nil, false) if not found.
func (ts *Toolset) Get(id string) (*tooladapter.CanonicalTool, bool) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	tool, ok := ts.tools[id]
	return tool, ok
}

// Remove removes a tool by ID. Returns true if found and removed.
func (ts *Toolset) Remove(id string) bool {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if _, ok := ts.tools[id]; ok {
		delete(ts.tools, id)
		return true
	}
	return false
}

// Count returns the number of tools.
func (ts *Toolset) Count() int {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return len(ts.tools)
}

// IDs returns tool IDs sorted lexicographically.
func (ts *Toolset) IDs() []string {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	ids := make([]string, 0, len(ts.tools))
	for id := range ts.tools {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

// Tools returns all tools sorted lexicographically by ID.
func (ts *Toolset) Tools() []*tooladapter.CanonicalTool {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	tools := make([]*tooladapter.CanonicalTool, 0, len(ts.tools))
	for _, t := range ts.tools {
		tools = append(tools, t)
	}
	sort.Slice(tools, func(i, j int) bool {
		return tools[i].ID() < tools[j].ID()
	})
	return tools
}

// Filter returns a new Toolset with tools matching fn.
// The original Toolset is not modified.
func (ts *Toolset) Filter(fn FilterFunc) *Toolset {
	ts.mu.RLock()
	// Snapshot matching tools while holding lock
	var matches []*tooladapter.CanonicalTool
	for _, t := range ts.tools {
		if fn(t) {
			matches = append(matches, t)
		}
	}
	ts.mu.RUnlock()

	// Build new toolset from snapshot (no lock needed)
	filtered := New(ts.name + "-filtered")
	for _, t := range matches {
		filtered.tools[t.ID()] = t
	}
	return filtered
}
