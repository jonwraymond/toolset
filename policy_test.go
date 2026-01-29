package toolset

import (
	"testing"

	"github.com/jonwraymond/tooladapter"
)

func TestPolicy_Interface(t *testing.T) {
	t.Run("Policy.Allow signature", func(t *testing.T) {
		var p Policy
		p = AllowAll() // must implement Policy
		tool := makeTool("ns", "foo", nil)
		_ = p.Allow(tool) // must accept *CanonicalTool and return bool
	})
}

func TestAllowAllPolicy(t *testing.T) {
	t.Run("returns true for all tools", func(t *testing.T) {
		p := AllowAll()
		tool := makeTool("ns", "foo", nil)
		if !p.Allow(tool) {
			t.Error("AllowAll should return true")
		}
	})

	t.Run("returns true for nil", func(t *testing.T) {
		p := AllowAll()
		if !p.Allow(nil) {
			t.Error("AllowAll should return true for nil")
		}
	})
}

func TestDenyAllPolicy(t *testing.T) {
	t.Run("returns false for all tools", func(t *testing.T) {
		p := DenyAll()
		tool := makeTool("ns", "foo", nil)
		if p.Allow(tool) {
			t.Error("DenyAll should return false")
		}
	})
}

func TestAllowNamespaces(t *testing.T) {
	t.Run("allows listed namespaces", func(t *testing.T) {
		p := AllowNamespaces("github", "gitlab")
		tool1 := makeTool("github", "list-repos", nil)
		tool2 := makeTool("gitlab", "list-repos", nil)

		if !p.Allow(tool1) {
			t.Error("Should allow github namespace")
		}
		if !p.Allow(tool2) {
			t.Error("Should allow gitlab namespace")
		}
	})

	t.Run("denies unlisted namespaces", func(t *testing.T) {
		p := AllowNamespaces("github")
		tool := makeTool("slack", "send-message", nil)
		if p.Allow(tool) {
			t.Error("Should deny slack namespace")
		}
	})

	t.Run("denies nil tool", func(t *testing.T) {
		p := AllowNamespaces("github")
		if p.Allow(nil) {
			t.Error("Should deny nil tool")
		}
	})
}

func TestDenyTags(t *testing.T) {
	t.Run("deny tools with any forbidden tag", func(t *testing.T) {
		p := DenyTags("dangerous", "deprecated")
		tool := makeTool("ns", "foo", []string{"safe", "dangerous"})
		if p.Allow(tool) {
			t.Error("Should deny tool with forbidden tag")
		}
	})

	t.Run("allow tools without forbidden tags", func(t *testing.T) {
		p := DenyTags("dangerous", "deprecated")
		tool := makeTool("ns", "foo", []string{"safe", "public"})
		if !p.Allow(tool) {
			t.Error("Should allow tool without forbidden tags")
		}
	})

	t.Run("denies nil tool", func(t *testing.T) {
		p := DenyTags("dangerous")
		if p.Allow(nil) {
			t.Error("Should deny nil tool")
		}
	})
}

func TestAllowScopes(t *testing.T) {
	t.Run("allow tools requiring only allowed scopes", func(t *testing.T) {
		p := AllowScopes("read:user", "read:repo")
		tool := &tooladapter.CanonicalTool{
			Name:           "list-repos",
			Namespace:      "github",
			RequiredScopes: []string{"read:repo"},
			InputSchema:    &tooladapter.JSONSchema{Type: "object"},
		}
		if !p.Allow(tool) {
			t.Error("Should allow tool with subset of allowed scopes")
		}
	})

	t.Run("deny tools requiring scopes not in allowed set", func(t *testing.T) {
		p := AllowScopes("read:user")
		tool := &tooladapter.CanonicalTool{
			Name:           "delete-repo",
			Namespace:      "github",
			RequiredScopes: []string{"write:repo"},
			InputSchema:    &tooladapter.JSONSchema{Type: "object"},
		}
		if p.Allow(tool) {
			t.Error("Should deny tool requiring unauthorized scope")
		}
	})

	t.Run("tools with no required scopes are allowed", func(t *testing.T) {
		p := AllowScopes("read:user")
		tool := &tooladapter.CanonicalTool{
			Name:           "public-info",
			Namespace:      "github",
			RequiredScopes: nil,
			InputSchema:    &tooladapter.JSONSchema{Type: "object"},
		}
		if !p.Allow(tool) {
			t.Error("Should allow tool with no required scopes")
		}
	})

	t.Run("denies nil tool", func(t *testing.T) {
		p := AllowScopes("read")
		if p.Allow(nil) {
			t.Error("Should deny nil tool")
		}
	})
}

func TestPolicyFunc(t *testing.T) {
	t.Run("PolicyFunc adapts function to Policy interface", func(t *testing.T) {
		var p Policy = PolicyFunc(func(t *tooladapter.CanonicalTool) bool {
			return t != nil && t.Name == "allowed"
		})

		tool1 := makeTool("ns", "allowed", nil)
		tool2 := makeTool("ns", "denied", nil)

		if !p.Allow(tool1) {
			t.Error("PolicyFunc should allow 'allowed'")
		}
		if p.Allow(tool2) {
			t.Error("PolicyFunc should deny 'denied'")
		}
	})
}

func TestBuilder_PolicyAfterFilters(t *testing.T) {
	t.Run("policy applied after all filters", func(t *testing.T) {
		tools := []*tooladapter.CanonicalTool{
			makeTool("github", "public", []string{"safe"}),
			makeTool("github", "private", []string{"internal"}),
			makeTool("slack", "send", []string{"safe"}),
		}

		// Filter to github namespace, then policy denies "internal" tag
		ts, err := NewBuilder("test").
			FromTools(tools).
			WithNamespace("github").
			WithPolicy(DenyTags("internal")).
			Build()

		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}

		// Should only have github:public (filtered to github, then policy denied internal)
		if ts.Count() != 1 {
			t.Errorf("Count() = %d, want 1", ts.Count())
		}
		if _, ok := ts.Get("github:public"); !ok {
			t.Error("github:public should be included")
		}
	})

	t.Run("policy can deny tools that passed filters", func(t *testing.T) {
		tools := []*tooladapter.CanonicalTool{
			makeTool("github", "list-repos", nil),
		}

		ts, err := NewBuilder("test").
			FromTools(tools).
			WithPolicy(DenyAll()).
			Build()

		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}
		if ts.Count() != 0 {
			t.Errorf("DenyAll policy should result in empty toolset")
		}
	})
}
