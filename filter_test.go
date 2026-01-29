package toolset

import (
	"testing"

	"github.com/jonwraymond/tooladapter"
)

func TestNamespaceFilter(t *testing.T) {
	t.Run("single namespace matches exact", func(t *testing.T) {
		filter := NamespaceFilter("github")
		tool := makeTool("github", "list-repos", nil)
		if !filter(tool) {
			t.Error("NamespaceFilter should match tool in namespace")
		}
	})

	t.Run("single namespace rejects other", func(t *testing.T) {
		filter := NamespaceFilter("github")
		tool := makeTool("slack", "send-message", nil)
		if filter(tool) {
			t.Error("NamespaceFilter should not match tool in different namespace")
		}
	})

	t.Run("multiple namespaces matches any", func(t *testing.T) {
		filter := NamespaceFilter("github", "gitlab")
		tool1 := makeTool("github", "list-repos", nil)
		tool2 := makeTool("gitlab", "list-repos", nil)
		tool3 := makeTool("bitbucket", "list-repos", nil)

		if !filter(tool1) {
			t.Error("Should match github")
		}
		if !filter(tool2) {
			t.Error("Should match gitlab")
		}
		if filter(tool3) {
			t.Error("Should not match bitbucket")
		}
	})

	t.Run("empty namespace list matches nothing", func(t *testing.T) {
		filter := NamespaceFilter()
		tool := makeTool("github", "list-repos", nil)
		if filter(tool) {
			t.Error("Empty namespace list should match nothing")
		}
	})

	t.Run("nil tool returns false", func(t *testing.T) {
		filter := NamespaceFilter("github")
		if filter(nil) {
			t.Error("Nil tool should return false")
		}
	})
}

func TestTagsAny(t *testing.T) {
	t.Run("tool with any tag matches", func(t *testing.T) {
		filter := TagsAny("read", "write")
		tool := makeTool("ns", "foo", []string{"read"})
		if !filter(tool) {
			t.Error("TagsAny should match tool with one of the tags")
		}
	})

	t.Run("tool with none of tags does not match", func(t *testing.T) {
		filter := TagsAny("read", "write")
		tool := makeTool("ns", "foo", []string{"delete", "admin"})
		if filter(tool) {
			t.Error("TagsAny should not match tool without any of the tags")
		}
	})

	t.Run("empty tag list matches nothing", func(t *testing.T) {
		filter := TagsAny()
		tool := makeTool("ns", "foo", []string{"read"})
		if filter(tool) {
			t.Error("Empty tag list should match nothing")
		}
	})

	t.Run("tool with no tags never matches", func(t *testing.T) {
		filter := TagsAny("read", "write")
		tool := makeTool("ns", "foo", nil)
		if filter(tool) {
			t.Error("Tool with no tags should not match TagsAny")
		}
	})

	t.Run("nil tool returns false", func(t *testing.T) {
		filter := TagsAny("read")
		if filter(nil) {
			t.Error("Nil tool should return false")
		}
	})
}

func TestTagsAll(t *testing.T) {
	t.Run("tool with all tags matches", func(t *testing.T) {
		filter := TagsAll("read", "public")
		tool := makeTool("ns", "foo", []string{"read", "public", "fast"})
		if !filter(tool) {
			t.Error("TagsAll should match tool with all required tags")
		}
	})

	t.Run("tool missing one tag does not match", func(t *testing.T) {
		filter := TagsAll("read", "write")
		tool := makeTool("ns", "foo", []string{"read"})
		if filter(tool) {
			t.Error("TagsAll should not match tool missing a required tag")
		}
	})

	t.Run("empty tag list matches everything", func(t *testing.T) {
		filter := TagsAll()
		tool := makeTool("ns", "foo", []string{"read"})
		if !filter(tool) {
			t.Error("Empty tag list should match everything (vacuously true)")
		}
	})

	t.Run("tool with superset of tags matches", func(t *testing.T) {
		filter := TagsAll("read")
		tool := makeTool("ns", "foo", []string{"read", "write", "admin"})
		if !filter(tool) {
			t.Error("TagsAll should match tool with superset of required tags")
		}
	})

	t.Run("nil tool returns false", func(t *testing.T) {
		filter := TagsAll("read")
		if filter(nil) {
			t.Error("Nil tool should return false")
		}
	})
}

func TestTagsNone(t *testing.T) {
	t.Run("tool with none of forbidden tags matches", func(t *testing.T) {
		filter := TagsNone("dangerous", "deprecated")
		tool := makeTool("ns", "foo", []string{"safe", "public"})
		if !filter(tool) {
			t.Error("TagsNone should match tool without forbidden tags")
		}
	})

	t.Run("tool with any forbidden tag does not match", func(t *testing.T) {
		filter := TagsNone("dangerous", "deprecated")
		tool := makeTool("ns", "foo", []string{"safe", "dangerous"})
		if filter(tool) {
			t.Error("TagsNone should not match tool with forbidden tag")
		}
	})

	t.Run("empty forbidden list matches everything", func(t *testing.T) {
		filter := TagsNone()
		tool := makeTool("ns", "foo", []string{"anything"})
		if !filter(tool) {
			t.Error("Empty forbidden list should match everything")
		}
	})

	t.Run("nil tool returns false", func(t *testing.T) {
		filter := TagsNone("dangerous")
		if filter(nil) {
			t.Error("Nil tool should return false")
		}
	})
}

func TestCategoryFilter(t *testing.T) {
	t.Run("matches exact category", func(t *testing.T) {
		filter := CategoryFilter("search")
		tool := &tooladapter.CanonicalTool{
			Name:        "find",
			Category:    "search",
			InputSchema: &tooladapter.JSONSchema{Type: "object"},
		}
		if !filter(tool) {
			t.Error("CategoryFilter should match tool in category")
		}
	})

	t.Run("multiple categories matches any", func(t *testing.T) {
		filter := CategoryFilter("search", "query")
		tool1 := &tooladapter.CanonicalTool{
			Name:        "find",
			Category:    "search",
			InputSchema: &tooladapter.JSONSchema{Type: "object"},
		}
		tool2 := &tooladapter.CanonicalTool{
			Name:        "lookup",
			Category:    "query",
			InputSchema: &tooladapter.JSONSchema{Type: "object"},
		}
		tool3 := &tooladapter.CanonicalTool{
			Name:        "delete",
			Category:    "mutation",
			InputSchema: &tooladapter.JSONSchema{Type: "object"},
		}

		if !filter(tool1) {
			t.Error("Should match search category")
		}
		if !filter(tool2) {
			t.Error("Should match query category")
		}
		if filter(tool3) {
			t.Error("Should not match mutation category")
		}
	})

	t.Run("empty matches nothing", func(t *testing.T) {
		filter := CategoryFilter()
		tool := &tooladapter.CanonicalTool{
			Name:        "find",
			Category:    "search",
			InputSchema: &tooladapter.JSONSchema{Type: "object"},
		}
		if filter(tool) {
			t.Error("Empty category list should match nothing")
		}
	})

	t.Run("nil tool returns false", func(t *testing.T) {
		filter := CategoryFilter("search")
		if filter(nil) {
			t.Error("Nil tool should return false")
		}
	})
}

func TestAllowIDs(t *testing.T) {
	t.Run("only listed IDs pass", func(t *testing.T) {
		filter := AllowIDs("ns:foo", "bar")
		tool1 := makeTool("ns", "foo", nil)
		tool2 := makeTool("", "bar", nil)
		tool3 := makeTool("ns", "baz", nil)

		if !filter(tool1) {
			t.Error("AllowIDs should pass ns:foo")
		}
		if !filter(tool2) {
			t.Error("AllowIDs should pass bar")
		}
		if filter(tool3) {
			t.Error("AllowIDs should not pass ns:baz")
		}
	})

	t.Run("unlisted IDs fail", func(t *testing.T) {
		filter := AllowIDs("allowed")
		tool := makeTool("", "denied", nil)
		if filter(tool) {
			t.Error("AllowIDs should not pass unlisted ID")
		}
	})

	t.Run("handles namespace:name format", func(t *testing.T) {
		filter := AllowIDs("github:list-repos")
		tool := makeTool("github", "list-repos", nil)
		if !filter(tool) {
			t.Error("AllowIDs should handle namespace:name format")
		}
	})

	t.Run("empty list allows nothing", func(t *testing.T) {
		filter := AllowIDs()
		tool := makeTool("", "any", nil)
		if filter(tool) {
			t.Error("Empty AllowIDs should allow nothing")
		}
	})

	t.Run("nil tool returns false", func(t *testing.T) {
		filter := AllowIDs("foo")
		if filter(nil) {
			t.Error("Nil tool should return false")
		}
	})
}

func TestDenyIDs(t *testing.T) {
	t.Run("listed IDs fail", func(t *testing.T) {
		filter := DenyIDs("ns:dangerous", "deprecated")
		tool1 := makeTool("ns", "dangerous", nil)
		tool2 := makeTool("", "deprecated", nil)

		if filter(tool1) {
			t.Error("DenyIDs should fail ns:dangerous")
		}
		if filter(tool2) {
			t.Error("DenyIDs should fail deprecated")
		}
	})

	t.Run("unlisted IDs pass", func(t *testing.T) {
		filter := DenyIDs("denied")
		tool := makeTool("", "allowed", nil)
		if !filter(tool) {
			t.Error("DenyIDs should pass unlisted ID")
		}
	})

	t.Run("empty list denies nothing", func(t *testing.T) {
		filter := DenyIDs()
		tool := makeTool("", "any", nil)
		if !filter(tool) {
			t.Error("Empty DenyIDs should deny nothing (all pass)")
		}
	})

	t.Run("nil tool returns false", func(t *testing.T) {
		filter := DenyIDs("foo")
		if filter(nil) {
			t.Error("Nil tool should return false")
		}
	})
}
