package toolset

import (
	"errors"
	"testing"

	"github.com/jonwraymond/tooladapter"
)

// mockAdapter implements tooladapter.Adapter for testing
type mockAdapter struct {
	name             string
	supportedFeatures map[tooladapter.SchemaFeature]bool
	fromCanonicalErr error
}

func (m *mockAdapter) Name() string {
	return m.name
}

func (m *mockAdapter) ToCanonical(raw any) (*tooladapter.CanonicalTool, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAdapter) FromCanonical(tool *tooladapter.CanonicalTool) (any, error) {
	if m.fromCanonicalErr != nil {
		return nil, m.fromCanonicalErr
	}
	// Return a simple map representation
	return map[string]any{
		"name":        tool.Name,
		"namespace":   tool.Namespace,
		"description": tool.Description,
	}, nil
}

func (m *mockAdapter) SupportsFeature(f tooladapter.SchemaFeature) bool {
	if m.supportedFeatures == nil {
		return true // default: supports all
	}
	return m.supportedFeatures[f]
}

func TestExposure_Export(t *testing.T) {
	t.Run("export returns slice of protocol tools", func(t *testing.T) {
		ts := New("test")
		ts.Add(makeTool("ns", "a", nil))
		ts.Add(makeTool("ns", "b", nil))

		adapter := &mockAdapter{name: "mock"}
		exp := NewExposure(ts, adapter)

		result, err := exp.Export()
		if err != nil {
			t.Fatalf("Export() error = %v", err)
		}
		if len(result) != 2 {
			t.Errorf("len(Export()) = %d, want 2", len(result))
		}
	})

	t.Run("each tool converted via adapter.FromCanonical", func(t *testing.T) {
		ts := New("test")
		ts.Add(&tooladapter.CanonicalTool{
			Name:        "foo",
			Namespace:   "ns",
			Description: "test tool",
			InputSchema: &tooladapter.JSONSchema{Type: "object"},
		})

		adapter := &mockAdapter{name: "mock"}
		exp := NewExposure(ts, adapter)

		result, err := exp.Export()
		if err != nil {
			t.Fatalf("Export() error = %v", err)
		}

		// Check the converted result has expected fields
		converted := result[0].(map[string]any)
		if converted["name"] != "foo" {
			t.Errorf("converted name = %v, want 'foo'", converted["name"])
		}
	})

	t.Run("order matches Tools order", func(t *testing.T) {
		ts := New("test")
		ts.Add(makeTool("", "zebra", nil))
		ts.Add(makeTool("", "apple", nil))

		adapter := &mockAdapter{name: "mock"}
		exp := NewExposure(ts, adapter)

		result, _ := exp.Export()

		// Tools() returns sorted by ID, so apple comes first
		first := result[0].(map[string]any)
		second := result[1].(map[string]any)
		if first["name"] != "apple" || second["name"] != "zebra" {
			t.Errorf("Export order doesn't match Tools order: got %v, %v", first["name"], second["name"])
		}
	})
}

func TestExposure_ExportWithWarnings(t *testing.T) {
	t.Run("returns warnings for unsupported features", func(t *testing.T) {
		ts := New("test")
		ts.Add(&tooladapter.CanonicalTool{
			Name:      "tool-with-pattern",
			Namespace: "ns",
			InputSchema: &tooladapter.JSONSchema{
				Type:    "object",
				Pattern: "^[a-z]+$", // uses pattern feature
			},
		})

		adapter := &mockAdapter{
			name: "mock",
			supportedFeatures: map[tooladapter.SchemaFeature]bool{
				tooladapter.FeaturePattern: false, // doesn't support pattern
			},
		}
		exp := NewExposure(ts, adapter)

		_, warnings := exp.ExportWithWarnings()
		if len(warnings) == 0 {
			t.Error("Expected warning for unsupported pattern feature")
		}

		found := false
		for _, w := range warnings {
			if w.Feature == tooladapter.FeaturePattern {
				found = true
				break
			}
		}
		if !found {
			t.Error("Warning should include pattern feature")
		}
	})

	t.Run("warnings include feature name and adapter names", func(t *testing.T) {
		ts := New("test")
		ts.Add(&tooladapter.CanonicalTool{
			Name:         "tool",
			Namespace:    "ns",
			SourceFormat: "mcp",
			InputSchema: &tooladapter.JSONSchema{
				Type:   "object",
				Format: "email",
			},
		})

		adapter := &mockAdapter{
			name: "openai",
			supportedFeatures: map[tooladapter.SchemaFeature]bool{
				tooladapter.FeatureFormat: false,
			},
		}
		exp := NewExposure(ts, adapter)

		_, warnings := exp.ExportWithWarnings()
		if len(warnings) == 0 {
			t.Fatal("Expected warning")
		}

		w := warnings[0]
		if w.Feature != tooladapter.FeatureFormat {
			t.Errorf("Feature = %v, want FeatureFormat", w.Feature)
		}
		if w.ToAdapter != "openai" {
			t.Errorf("ToAdapter = %q, want 'openai'", w.ToAdapter)
		}
	})
}

func TestExposure_NilAdapter(t *testing.T) {
	t.Run("export returns error for nil adapter", func(t *testing.T) {
		ts := New("test")
		ts.Add(makeTool("ns", "a", nil))

		exp := NewExposure(ts, nil)

		_, err := exp.Export()
		if err == nil {
			t.Error("Export() should return error for nil adapter")
		}
	})
}

func TestExposure_EmptyToolset(t *testing.T) {
	t.Run("returns empty slice no warnings no error", func(t *testing.T) {
		ts := New("empty")
		adapter := &mockAdapter{name: "mock"}
		exp := NewExposure(ts, adapter)

		result, err := exp.Export()
		if err != nil {
			t.Fatalf("Export() error = %v", err)
		}
		if len(result) != 0 {
			t.Errorf("len(Export()) = %d, want 0", len(result))
		}

		resultW, warnings := exp.ExportWithWarnings()
		if len(resultW) != 0 {
			t.Errorf("len(ExportWithWarnings result) = %d, want 0", len(resultW))
		}
		if len(warnings) != 0 {
			t.Errorf("len(warnings) = %d, want 0", len(warnings))
		}
	})
}

func TestExposure_NestedSchema(t *testing.T) {
	t.Run("tool with nested Properties/Items/Defs", func(t *testing.T) {
		ts := New("test")
		ts.Add(&tooladapter.CanonicalTool{
			Name:      "nested",
			Namespace: "ns",
			InputSchema: &tooladapter.JSONSchema{
				Type: "object",
				Properties: map[string]*tooladapter.JSONSchema{
					"items": {
						Type: "array",
						Items: &tooladapter.JSONSchema{
							Type:    "string",
							Pattern: "^[a-z]+$", // nested pattern
						},
					},
				},
			},
		})

		adapter := &mockAdapter{
			name: "mock",
			supportedFeatures: map[tooladapter.SchemaFeature]bool{
				tooladapter.FeaturePattern: false,
			},
		}
		exp := NewExposure(ts, adapter)

		_, warnings := exp.ExportWithWarnings()
		if len(warnings) == 0 {
			t.Error("Should detect pattern in nested Items")
		}
	})

	t.Run("detectFeatures finds features at all depths", func(t *testing.T) {
		ts := New("test")
		ts.Add(&tooladapter.CanonicalTool{
			Name:      "deep",
			Namespace: "ns",
			InputSchema: &tooladapter.JSONSchema{
				Type: "object",
				Defs: map[string]*tooladapter.JSONSchema{
					"inner": {
						Type:   "string",
						Format: "email", // in defs
					},
				},
			},
		})

		adapter := &mockAdapter{
			name: "mock",
			supportedFeatures: map[tooladapter.SchemaFeature]bool{
				tooladapter.FeatureFormat: false,
				tooladapter.FeatureDefs:   false,
			},
		}
		exp := NewExposure(ts, adapter)

		_, warnings := exp.ExportWithWarnings()

		foundFormat := false
		foundDefs := false
		for _, w := range warnings {
			if w.Feature == tooladapter.FeatureFormat {
				foundFormat = true
			}
			if w.Feature == tooladapter.FeatureDefs {
				foundDefs = true
			}
		}
		if !foundFormat {
			t.Error("Should detect format in nested Defs")
		}
		if !foundDefs {
			t.Error("Should detect Defs usage")
		}
	})
}

func TestExposure_Combinators(t *testing.T) {
	t.Run("tool with anyOf/oneOf/allOf/not", func(t *testing.T) {
		ts := New("test")
		ts.Add(&tooladapter.CanonicalTool{
			Name:      "combinators",
			Namespace: "ns",
			InputSchema: &tooladapter.JSONSchema{
				Type: "object",
				AnyOf: []*tooladapter.JSONSchema{
					{Type: "string"},
					{Type: "number"},
				},
				OneOf: []*tooladapter.JSONSchema{
					{Type: "boolean"},
				},
				AllOf: []*tooladapter.JSONSchema{
					{Type: "object"},
				},
				Not: &tooladapter.JSONSchema{Type: "null"},
			},
		})

		adapter := &mockAdapter{
			name: "mock",
			supportedFeatures: map[tooladapter.SchemaFeature]bool{
				tooladapter.FeatureAnyOf: false,
				tooladapter.FeatureOneOf: false,
				tooladapter.FeatureAllOf: false,
				tooladapter.FeatureNot:   false,
			},
		}
		exp := NewExposure(ts, adapter)

		_, warnings := exp.ExportWithWarnings()

		features := make(map[tooladapter.SchemaFeature]bool)
		for _, w := range warnings {
			features[w.Feature] = true
		}

		if !features[tooladapter.FeatureAnyOf] {
			t.Error("Should warn about anyOf")
		}
		if !features[tooladapter.FeatureOneOf] {
			t.Error("Should warn about oneOf")
		}
		if !features[tooladapter.FeatureAllOf] {
			t.Error("Should warn about allOf")
		}
		if !features[tooladapter.FeatureNot] {
			t.Error("Should warn about not")
		}
	})

	t.Run("detectFeatures walks combinator branches", func(t *testing.T) {
		ts := New("test")
		ts.Add(&tooladapter.CanonicalTool{
			Name:      "nested-combinator",
			Namespace: "ns",
			InputSchema: &tooladapter.JSONSchema{
				Type: "object",
				AnyOf: []*tooladapter.JSONSchema{
					{
						Type:    "string",
						Pattern: "^test$", // pattern inside anyOf branch
					},
				},
			},
		})

		adapter := &mockAdapter{
			name: "mock",
			supportedFeatures: map[tooladapter.SchemaFeature]bool{
				tooladapter.FeaturePattern: false,
				tooladapter.FeatureAnyOf:   true, // supports anyOf but not pattern
			},
		}
		exp := NewExposure(ts, adapter)

		_, warnings := exp.ExportWithWarnings()

		found := false
		for _, w := range warnings {
			if w.Feature == tooladapter.FeaturePattern {
				found = true
				break
			}
		}
		if !found {
			t.Error("Should detect pattern inside anyOf branch")
		}
	})
}

func TestExposure_RefInNested(t *testing.T) {
	t.Run("$ref inside Properties", func(t *testing.T) {
		ts := New("test")
		ts.Add(&tooladapter.CanonicalTool{
			Name:      "with-ref",
			Namespace: "ns",
			InputSchema: &tooladapter.JSONSchema{
				Type: "object",
				Properties: map[string]*tooladapter.JSONSchema{
					"user": {
						Ref: "#/$defs/User",
					},
				},
			},
		})

		adapter := &mockAdapter{
			name: "mock",
			supportedFeatures: map[tooladapter.SchemaFeature]bool{
				tooladapter.FeatureRef: false,
			},
		}
		exp := NewExposure(ts, adapter)

		_, warnings := exp.ExportWithWarnings()

		found := false
		for _, w := range warnings {
			if w.Feature == tooladapter.FeatureRef {
				found = true
				break
			}
		}
		if !found {
			t.Error("Should warn about $ref in nested property")
		}
	})
}
