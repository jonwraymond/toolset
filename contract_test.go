package toolset

import (
	"testing"

	"github.com/jonwraymond/tooladapter"
)

func TestPolicyContract_Builtins(t *testing.T) {
	tool := &tooladapter.CanonicalTool{
		Namespace:      "alpha",
		Name:           "echo",
		Tags:           []string{"safe"},
		RequiredScopes: []string{"read"},
	}

	if AllowAll().Allow(nil) {
		t.Fatalf("AllowAll should deny nil tool")
	}
	if !AllowAll().Allow(tool) {
		t.Fatalf("AllowAll should allow non-nil tool")
	}
	if DenyAll().Allow(tool) {
		t.Fatalf("DenyAll should deny tool")
	}

	ns := AllowNamespaces("alpha")
	if !ns.Allow(tool) {
		t.Fatalf("AllowNamespaces should allow matching namespace")
	}
	if ns.Allow(&tooladapter.CanonicalTool{Namespace: "beta", Name: "echo"}) {
		t.Fatalf("AllowNamespaces should deny non-matching namespace")
	}

	denyTags := DenyTags("danger")
	if !denyTags.Allow(tool) {
		t.Fatalf("DenyTags should allow tools without denied tags")
	}
	tool.Tags = []string{"danger"}
	if denyTags.Allow(tool) {
		t.Fatalf("DenyTags should deny tools with denied tags")
	}

	allowScopes := AllowScopes("read", "write")
	tool.Tags = []string{"safe"}
	tool.RequiredScopes = []string{"read"}
	if !allowScopes.Allow(tool) {
		t.Fatalf("AllowScopes should allow tools with allowed scopes")
	}
	tool.RequiredScopes = []string{"admin"}
	if allowScopes.Allow(tool) {
		t.Fatalf("AllowScopes should deny tools with disallowed scopes")
	}
}

func TestPolicyContract_NilFunc(t *testing.T) {
	var fn PolicyFunc
	if fn.Allow(&tooladapter.CanonicalTool{Name: "noop"}) {
		t.Fatalf("nil PolicyFunc should deny")
	}
}

type stubRegistry struct {
	tools []*tooladapter.CanonicalTool
}

func (s stubRegistry) Tools() []*tooladapter.CanonicalTool {
	return s.tools
}

func TestRegistryContract_BuilderUsesSnapshot(t *testing.T) {
	tools := []*tooladapter.CanonicalTool{
		{Name: "one", InputSchema: &tooladapter.JSONSchema{Type: "object"}},
	}
	reg := stubRegistry{tools: tools}

	ts, err := NewBuilder("snapshot").FromRegistry(reg).Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	if ts.Count() != 1 {
		t.Fatalf("expected 1 tool, got %d", ts.Count())
	}

	// Mutate the original slice to ensure Toolset retains its own snapshot.
	tools[0] = &tooladapter.CanonicalTool{Name: "changed", InputSchema: &tooladapter.JSONSchema{Type: "object"}}
	if ts.Count() != 1 {
		t.Fatalf("toolset should not be affected by registry slice mutation")
	}
}

func TestRegistryContract_NilTools(t *testing.T) {
	ts, err := NewBuilder("nil-registry").FromRegistry(stubRegistry{tools: nil}).Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	if ts.Count() != 0 {
		t.Fatalf("expected empty toolset for nil registry slice, got %d", ts.Count())
	}
}
