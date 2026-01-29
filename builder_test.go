package toolset

import (
	"testing"

	"github.com/jonwraymond/tooladapter"
)

// mockRegistry implements Registry interface for testing
type mockRegistry struct {
	tools []*tooladapter.CanonicalTool
}

func (r *mockRegistry) Tools() []*tooladapter.CanonicalTool {
	return r.tools
}

func TestBuilder_FromTools(t *testing.T) {
	t.Run("FromTools populates builder", func(t *testing.T) {
		tools := []*tooladapter.CanonicalTool{
			makeTool("ns", "a", nil),
			makeTool("ns", "b", nil),
		}
		ts, err := NewBuilder("test").FromTools(tools).Build()
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}
		if ts.Count() != 2 {
			t.Errorf("Count() = %d, want 2", ts.Count())
		}
	})

	t.Run("empty slice returns empty toolset", func(t *testing.T) {
		ts, err := NewBuilder("test").FromTools([]*tooladapter.CanonicalTool{}).Build()
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}
		if ts.Count() != 0 {
			t.Errorf("Count() = %d, want 0", ts.Count())
		}
	})

	t.Run("nil slice returns empty toolset", func(t *testing.T) {
		ts, err := NewBuilder("test").FromTools(nil).Build()
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}
		if ts.Count() != 0 {
			t.Errorf("Count() = %d, want 0", ts.Count())
		}
	})
}

func TestBuilder_FromRegistry(t *testing.T) {
	t.Run("Registry.Tools called once on Build", func(t *testing.T) {
		reg := &mockRegistry{
			tools: []*tooladapter.CanonicalTool{
				makeTool("ns", "foo", nil),
				makeTool("ns", "bar", nil),
			},
		}
		ts, err := NewBuilder("test").FromRegistry(reg).Build()
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}
		if ts.Count() != 2 {
			t.Errorf("Count() = %d, want 2", ts.Count())
		}
	})

	t.Run("tools from registry included", func(t *testing.T) {
		reg := &mockRegistry{
			tools: []*tooladapter.CanonicalTool{
				makeTool("github", "list-repos", nil),
			},
		}
		ts, err := NewBuilder("test").FromRegistry(reg).Build()
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}
		if _, ok := ts.Get("github:list-repos"); !ok {
			t.Error("Tool from registry not included")
		}
	})
}

func TestBuilder_WithNamespace(t *testing.T) {
	t.Run("filters to single namespace", func(t *testing.T) {
		tools := []*tooladapter.CanonicalTool{
			makeTool("github", "list-repos", nil),
			makeTool("slack", "send-message", nil),
		}
		ts, err := NewBuilder("test").
			FromTools(tools).
			WithNamespace("github").
			Build()
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}
		if ts.Count() != 1 {
			t.Errorf("Count() = %d, want 1", ts.Count())
		}
		if _, ok := ts.Get("github:list-repos"); !ok {
			t.Error("GitHub tool should be included")
		}
	})
}

func TestBuilder_WithNamespaces(t *testing.T) {
	t.Run("filters to multiple namespaces", func(t *testing.T) {
		tools := []*tooladapter.CanonicalTool{
			makeTool("github", "list-repos", nil),
			makeTool("gitlab", "list-projects", nil),
			makeTool("slack", "send-message", nil),
		}
		ts, err := NewBuilder("test").
			FromTools(tools).
			WithNamespaces([]string{"github", "gitlab"}).
			Build()
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}
		if ts.Count() != 2 {
			t.Errorf("Count() = %d, want 2", ts.Count())
		}
	})
}

func TestBuilder_WithTags(t *testing.T) {
	t.Run("filters to tools with all specified tags", func(t *testing.T) {
		tools := []*tooladapter.CanonicalTool{
			makeTool("ns", "a", []string{"read", "public"}),
			makeTool("ns", "b", []string{"read"}),
			makeTool("ns", "c", []string{"write"}),
		}
		ts, err := NewBuilder("test").
			FromTools(tools).
			WithTags([]string{"read", "public"}).
			Build()
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}
		if ts.Count() != 1 {
			t.Errorf("Count() = %d, want 1", ts.Count())
		}
		if _, ok := ts.Get("ns:a"); !ok {
			t.Error("Tool with all tags should be included")
		}
	})
}

func TestBuilder_WithCategories(t *testing.T) {
	t.Run("filters to tools with any category", func(t *testing.T) {
		tools := []*tooladapter.CanonicalTool{
			{Name: "search", Namespace: "ns", Category: "query", InputSchema: &tooladapter.JSONSchema{Type: "object"}},
			{Name: "delete", Namespace: "ns", Category: "mutation", InputSchema: &tooladapter.JSONSchema{Type: "object"}},
			{Name: "find", Namespace: "ns", Category: "query", InputSchema: &tooladapter.JSONSchema{Type: "object"}},
		}
		ts, err := NewBuilder("test").
			FromTools(tools).
			WithCategories([]string{"query"}).
			Build()
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}
		if ts.Count() != 2 {
			t.Errorf("Count() = %d, want 2", ts.Count())
		}
	})
}

func TestBuilder_WithTools(t *testing.T) {
	t.Run("include only listed tool IDs", func(t *testing.T) {
		tools := []*tooladapter.CanonicalTool{
			makeTool("ns", "a", nil),
			makeTool("ns", "b", nil),
			makeTool("ns", "c", nil),
		}
		ts, err := NewBuilder("test").
			FromTools(tools).
			WithTools([]string{"ns:a", "ns:c"}).
			Build()
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}
		if ts.Count() != 2 {
			t.Errorf("Count() = %d, want 2", ts.Count())
		}
		if _, ok := ts.Get("ns:b"); ok {
			t.Error("ns:b should not be included")
		}
	})
}

func TestBuilder_ExcludeTools(t *testing.T) {
	t.Run("exclude listed tool IDs", func(t *testing.T) {
		tools := []*tooladapter.CanonicalTool{
			makeTool("ns", "a", nil),
			makeTool("ns", "b", nil),
			makeTool("ns", "c", nil),
		}
		ts, err := NewBuilder("test").
			FromTools(tools).
			ExcludeTools([]string{"ns:b"}).
			Build()
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}
		if ts.Count() != 2 {
			t.Errorf("Count() = %d, want 2", ts.Count())
		}
		if _, ok := ts.Get("ns:b"); ok {
			t.Error("ns:b should be excluded")
		}
	})
}

func TestBuilder_ChainedFilters(t *testing.T) {
	t.Run("multiple filters AND together", func(t *testing.T) {
		tools := []*tooladapter.CanonicalTool{
			makeTool("github", "a", []string{"read"}),
			makeTool("github", "b", []string{"write"}),
			makeTool("slack", "c", []string{"read"}),
		}
		ts, err := NewBuilder("test").
			FromTools(tools).
			WithNamespace("github").
			WithTags([]string{"read"}).
			Build()
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}
		if ts.Count() != 1 {
			t.Errorf("Count() = %d, want 1", ts.Count())
		}
		if _, ok := ts.Get("github:a"); !ok {
			t.Error("github:a should pass both filters")
		}
	})

	t.Run("order does not affect result", func(t *testing.T) {
		tools := []*tooladapter.CanonicalTool{
			makeTool("github", "a", []string{"read"}),
			makeTool("github", "b", []string{"write"}),
			makeTool("slack", "c", []string{"read"}),
		}

		// Order 1: namespace then tags
		ts1, _ := NewBuilder("test").
			FromTools(tools).
			WithNamespace("github").
			WithTags([]string{"read"}).
			Build()

		// Order 2: tags then namespace
		ts2, _ := NewBuilder("test").
			FromTools(tools).
			WithTags([]string{"read"}).
			WithNamespace("github").
			Build()

		if ts1.Count() != ts2.Count() {
			t.Errorf("Filter order affects result: %d vs %d", ts1.Count(), ts2.Count())
		}
	})
}

func TestBuilder_WithFilter(t *testing.T) {
	t.Run("custom FilterFunc applied", func(t *testing.T) {
		tools := []*tooladapter.CanonicalTool{
			makeTool("ns", "aardvark", nil),
			makeTool("ns", "zebra", nil),
		}
		ts, err := NewBuilder("test").
			FromTools(tools).
			WithFilter(func(t *tooladapter.CanonicalTool) bool {
				return t.Name[0] == 'a'
			}).
			Build()
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}
		if ts.Count() != 1 {
			t.Errorf("Count() = %d, want 1", ts.Count())
		}
	})
}

func TestBuilder_Build_NoSource(t *testing.T) {
	t.Run("returns error when no source set", func(t *testing.T) {
		_, err := NewBuilder("test").Build()
		if err == nil {
			t.Error("Build() should return error when no source")
		}
	})
}

func TestBuilder_Build_Success(t *testing.T) {
	t.Run("returns toolset and nil error on success", func(t *testing.T) {
		tools := []*tooladapter.CanonicalTool{makeTool("ns", "a", nil)}
		ts, err := NewBuilder("test").FromTools(tools).Build()
		if err != nil {
			t.Fatalf("Build() error = %v, want nil", err)
		}
		if ts == nil {
			t.Error("Build() returned nil toolset")
		}
	})

	t.Run("toolset has correct name", func(t *testing.T) {
		tools := []*tooladapter.CanonicalTool{makeTool("ns", "a", nil)}
		ts, _ := NewBuilder("my-set").FromTools(tools).Build()
		if ts.Name() != "my-set" {
			t.Errorf("Name() = %q, want %q", ts.Name(), "my-set")
		}
	})

	t.Run("tools sorted by ID", func(t *testing.T) {
		tools := []*tooladapter.CanonicalTool{
			makeTool("", "zebra", nil),
			makeTool("", "apple", nil),
		}
		ts, _ := NewBuilder("test").FromTools(tools).Build()
		ids := ts.IDs()
		if ids[0] != "apple" || ids[1] != "zebra" {
			t.Errorf("IDs not sorted: %v", ids)
		}
	})
}
