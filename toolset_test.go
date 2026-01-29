package toolset

import (
	"sync"
	"testing"

	"github.com/jonwraymond/tooladapter"
)

// Helper to create a test tool
func makeTool(namespace, name string, tags []string) *tooladapter.CanonicalTool {
	return &tooladapter.CanonicalTool{
		Namespace:   namespace,
		Name:        name,
		Tags:        tags,
		InputSchema: &tooladapter.JSONSchema{Type: "object"},
	}
}

func TestToolset_New(t *testing.T) {
	t.Run("creates toolset with name", func(t *testing.T) {
		ts := New("test-set")
		if ts.Name() != "test-set" {
			t.Errorf("Name() = %q, want %q", ts.Name(), "test-set")
		}
	})

	t.Run("new toolset has count 0", func(t *testing.T) {
		ts := New("empty")
		if ts.Count() != 0 {
			t.Errorf("Count() = %d, want 0", ts.Count())
		}
	})

	t.Run("new toolset Tools returns empty slice not nil", func(t *testing.T) {
		ts := New("empty")
		tools := ts.Tools()
		if tools == nil {
			t.Error("Tools() returned nil, want empty slice")
		}
		if len(tools) != 0 {
			t.Errorf("len(Tools()) = %d, want 0", len(tools))
		}
	})

	t.Run("new toolset IDs returns empty slice not nil", func(t *testing.T) {
		ts := New("empty")
		ids := ts.IDs()
		if ids == nil {
			t.Error("IDs() returned nil, want empty slice")
		}
		if len(ids) != 0 {
			t.Errorf("len(IDs()) = %d, want 0", len(ids))
		}
	})
}

func TestToolset_Add(t *testing.T) {
	t.Run("add tool then get returns it", func(t *testing.T) {
		ts := New("test")
		tool := makeTool("ns", "foo", nil)
		ts.Add(tool)

		got, ok := ts.Get("ns:foo")
		if !ok {
			t.Fatal("Get returned false, want true")
		}
		if got != tool {
			t.Error("Get returned different tool instance")
		}
	})

	t.Run("add nil is no-op", func(t *testing.T) {
		ts := New("test")
		ts.Add(nil) // should not panic
		if ts.Count() != 0 {
			t.Errorf("Count() = %d after Add(nil), want 0", ts.Count())
		}
	})

	t.Run("add duplicate ID replaces existing", func(t *testing.T) {
		ts := New("test")
		tool1 := makeTool("ns", "foo", []string{"v1"})
		tool2 := makeTool("ns", "foo", []string{"v2"})

		ts.Add(tool1)
		ts.Add(tool2)

		if ts.Count() != 1 {
			t.Errorf("Count() = %d, want 1", ts.Count())
		}

		got, _ := ts.Get("ns:foo")
		if got.Tags[0] != "v2" {
			t.Errorf("Tool not replaced: tags = %v, want [v2]", got.Tags)
		}
	})
}

func TestToolset_Get(t *testing.T) {
	t.Run("get existing returns tool and true", func(t *testing.T) {
		ts := New("test")
		tool := makeTool("ns", "bar", nil)
		ts.Add(tool)

		got, ok := ts.Get("ns:bar")
		if !ok {
			t.Error("Get returned false, want true")
		}
		if got != tool {
			t.Error("Get returned wrong tool")
		}
	})

	t.Run("get missing returns nil and false", func(t *testing.T) {
		ts := New("test")
		got, ok := ts.Get("missing")
		if ok {
			t.Error("Get returned true for missing tool")
		}
		if got != nil {
			t.Error("Get returned non-nil for missing tool")
		}
	})

	t.Run("get by namespace:name format works", func(t *testing.T) {
		ts := New("test")
		tool := makeTool("github", "list-repos", nil)
		ts.Add(tool)

		got, ok := ts.Get("github:list-repos")
		if !ok || got != tool {
			t.Error("Get by namespace:name failed")
		}
	})

	t.Run("get by name only works for unnamespaced tools", func(t *testing.T) {
		ts := New("test")
		tool := makeTool("", "simple", nil)
		ts.Add(tool)

		got, ok := ts.Get("simple")
		if !ok || got != tool {
			t.Error("Get by name only failed for unnamespaced tool")
		}
	})
}

func TestToolset_Remove(t *testing.T) {
	t.Run("remove existing returns true", func(t *testing.T) {
		ts := New("test")
		ts.Add(makeTool("ns", "foo", nil))

		if !ts.Remove("ns:foo") {
			t.Error("Remove returned false for existing tool")
		}
	})

	t.Run("remove missing returns false", func(t *testing.T) {
		ts := New("test")
		if ts.Remove("missing") {
			t.Error("Remove returned true for missing tool")
		}
	})

	t.Run("removed tool no longer in IDs or Tools", func(t *testing.T) {
		ts := New("test")
		ts.Add(makeTool("ns", "foo", nil))
		ts.Remove("ns:foo")

		if ts.Count() != 0 {
			t.Errorf("Count() = %d after remove, want 0", ts.Count())
		}
		if _, ok := ts.Get("ns:foo"); ok {
			t.Error("Get returned true for removed tool")
		}
	})
}

func TestToolset_Count(t *testing.T) {
	t.Run("empty toolset has count 0", func(t *testing.T) {
		ts := New("test")
		if ts.Count() != 0 {
			t.Errorf("Count() = %d, want 0", ts.Count())
		}
	})

	t.Run("count increments on add", func(t *testing.T) {
		ts := New("test")
		ts.Add(makeTool("", "a", nil))
		ts.Add(makeTool("", "b", nil))
		if ts.Count() != 2 {
			t.Errorf("Count() = %d, want 2", ts.Count())
		}
	})

	t.Run("count decrements on remove", func(t *testing.T) {
		ts := New("test")
		ts.Add(makeTool("", "a", nil))
		ts.Add(makeTool("", "b", nil))
		ts.Remove("a")
		if ts.Count() != 1 {
			t.Errorf("Count() = %d, want 1", ts.Count())
		}
	})
}

func TestToolset_IDs(t *testing.T) {
	t.Run("returns sorted slice", func(t *testing.T) {
		ts := New("test")
		// Add in non-alphabetical order
		ts.Add(makeTool("", "zebra", nil))
		ts.Add(makeTool("", "apple", nil))
		ts.Add(makeTool("ns", "middle", nil))
		ts.Add(makeTool("", "banana", nil))
		ts.Add(makeTool("aaa", "first", nil))

		ids := ts.IDs()
		expected := []string{"aaa:first", "apple", "banana", "ns:middle", "zebra"}

		if len(ids) != len(expected) {
			t.Fatalf("len(IDs()) = %d, want %d", len(ids), len(expected))
		}
		for i, id := range ids {
			if id != expected[i] {
				t.Errorf("IDs()[%d] = %q, want %q", i, id, expected[i])
			}
		}
	})

	t.Run("multiple calls return same order", func(t *testing.T) {
		ts := New("test")
		ts.Add(makeTool("", "c", nil))
		ts.Add(makeTool("", "a", nil))
		ts.Add(makeTool("", "b", nil))

		ids1 := ts.IDs()
		ids2 := ts.IDs()

		for i := range ids1 {
			if ids1[i] != ids2[i] {
				t.Errorf("IDs() not deterministic: got %v and %v", ids1, ids2)
				break
			}
		}
	})
}

func TestToolset_Tools(t *testing.T) {
	t.Run("returns sorted slice by ID", func(t *testing.T) {
		ts := New("test")
		ts.Add(makeTool("", "zebra", nil))
		ts.Add(makeTool("", "apple", nil))
		ts.Add(makeTool("ns", "middle", nil))

		tools := ts.Tools()
		if len(tools) != 3 {
			t.Fatalf("len(Tools()) = %d, want 3", len(tools))
		}

		// Check order matches sorted IDs
		expected := []string{"apple", "ns:middle", "zebra"}
		for i, tool := range tools {
			if tool.ID() != expected[i] {
				t.Errorf("Tools()[%d].ID() = %q, want %q", i, tool.ID(), expected[i])
			}
		}
	})

	t.Run("does not alias internal storage", func(t *testing.T) {
		ts := New("test")
		ts.Add(makeTool("", "a", nil))
		ts.Add(makeTool("", "b", nil))

		tools := ts.Tools()
		tools[0] = nil // modify returned slice

		// Original should be unchanged
		got, ok := ts.Get("a")
		if !ok || got == nil {
			t.Error("Modifying returned slice affected internal storage")
		}
	})

	t.Run("deterministic across calls", func(t *testing.T) {
		ts := New("test")
		ts.Add(makeTool("", "b", nil))
		ts.Add(makeTool("", "a", nil))

		tools1 := ts.Tools()
		tools2 := ts.Tools()

		for i := range tools1 {
			if tools1[i].ID() != tools2[i].ID() {
				t.Error("Tools() not deterministic")
				break
			}
		}
	})
}

func TestToolset_Filter(t *testing.T) {
	t.Run("filter returns new toolset instance", func(t *testing.T) {
		ts := New("original")
		ts.Add(makeTool("", "a", nil))

		filtered := ts.Filter(func(*tooladapter.CanonicalTool) bool { return true })

		if filtered == ts {
			t.Error("Filter returned same instance, want new instance")
		}
	})

	t.Run("original unchanged after filter", func(t *testing.T) {
		ts := New("original")
		ts.Add(makeTool("", "a", nil))
		ts.Add(makeTool("", "b", nil))

		_ = ts.Filter(func(t *tooladapter.CanonicalTool) bool {
			return t.Name == "a"
		})

		if ts.Count() != 2 {
			t.Errorf("Original count = %d after filter, want 2", ts.Count())
		}
	})

	t.Run("filter with always-true returns all tools", func(t *testing.T) {
		ts := New("test")
		ts.Add(makeTool("", "a", nil))
		ts.Add(makeTool("", "b", nil))

		filtered := ts.Filter(func(*tooladapter.CanonicalTool) bool { return true })

		if filtered.Count() != 2 {
			t.Errorf("Filtered count = %d, want 2", filtered.Count())
		}
	})

	t.Run("filter with always-false returns empty toolset", func(t *testing.T) {
		ts := New("test")
		ts.Add(makeTool("", "a", nil))
		ts.Add(makeTool("", "b", nil))

		filtered := ts.Filter(func(*tooladapter.CanonicalTool) bool { return false })

		if filtered.Count() != 0 {
			t.Errorf("Filtered count = %d, want 0", filtered.Count())
		}
	})

	t.Run("filtered toolset has derived name", func(t *testing.T) {
		ts := New("original")
		filtered := ts.Filter(func(*tooladapter.CanonicalTool) bool { return true })

		if filtered.Name() != "original-filtered" {
			t.Errorf("Filtered name = %q, want %q", filtered.Name(), "original-filtered")
		}
	})
}

func TestToolset_Concurrency(t *testing.T) {
	ts := New("concurrent")

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(3)

		go func(n int) {
			defer wg.Done()
			tool := makeTool("ns", string(rune('a'+n%26)), nil)
			ts.Add(tool)
		}(i)

		go func(n int) {
			defer wg.Done()
			id := "ns:" + string(rune('a'+n%26))
			ts.Get(id)
		}(i)

		go func(n int) {
			defer wg.Done()
			id := "ns:" + string(rune('a'+n%26))
			ts.Remove(id)
		}(i)
	}

	// Also test Tools and IDs concurrently
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			_ = ts.Tools()
		}()
		go func() {
			defer wg.Done()
			_ = ts.IDs()
		}()
	}

	wg.Wait()
	// If we get here without panics or race detector complaints, test passes
}
