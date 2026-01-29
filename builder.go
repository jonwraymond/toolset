package toolset

import (
	"errors"

	"github.com/jonwraymond/tooladapter"
)

// Registry provides tools for the builder.
type Registry interface {
	Tools() []*tooladapter.CanonicalTool
}

// Builder constructs Toolsets with filtering.
type Builder struct {
	name      string
	source    []*tooladapter.CanonicalTool
	sourceSet bool // tracks whether FromTools was called (even with nil)
	registry  Registry
	filters   []FilterFunc
	policy    Policy
}

// NewBuilder creates a new Builder with the given toolset name.
func NewBuilder(name string) *Builder {
	return &Builder{name: name}
}

// FromTools sets tools as the source.
func (b *Builder) FromTools(tools []*tooladapter.CanonicalTool) *Builder {
	b.source = tools
	b.sourceSet = true
	return b
}

// FromRegistry sets a registry as the source.
func (b *Builder) FromRegistry(r Registry) *Builder {
	b.registry = r
	return b
}

// WithNamespace filters to a single namespace.
func (b *Builder) WithNamespace(ns string) *Builder {
	b.filters = append(b.filters, NamespaceFilter(ns))
	return b
}

// WithNamespaces filters to multiple namespaces.
func (b *Builder) WithNamespaces(ns []string) *Builder {
	b.filters = append(b.filters, NamespaceFilter(ns...))
	return b
}

// WithTags filters to tools with ALL specified tags.
func (b *Builder) WithTags(tags []string) *Builder {
	b.filters = append(b.filters, TagsAll(tags...))
	return b
}

// WithCategories filters to tools with ANY category.
func (b *Builder) WithCategories(categories []string) *Builder {
	b.filters = append(b.filters, CategoryFilter(categories...))
	return b
}

// WithTools includes only listed tool IDs.
func (b *Builder) WithTools(ids []string) *Builder {
	b.filters = append(b.filters, AllowIDs(ids...))
	return b
}

// ExcludeTools excludes listed tool IDs.
func (b *Builder) ExcludeTools(ids []string) *Builder {
	b.filters = append(b.filters, DenyIDs(ids...))
	return b
}

// WithFilter adds a custom filter.
func (b *Builder) WithFilter(fn FilterFunc) *Builder {
	b.filters = append(b.filters, fn)
	return b
}

// WithPolicy sets the access control policy (applied after filters).
func (b *Builder) WithPolicy(p Policy) *Builder {
	b.policy = p
	return b
}

// Build creates the Toolset.
func (b *Builder) Build() (*Toolset, error) {
	// Gather source tools
	var tools []*tooladapter.CanonicalTool
	if b.registry != nil {
		tools = b.registry.Tools()
	} else if b.source != nil || b.sourceSet {
		tools = b.source
	} else {
		return nil, errors.New("no source: call FromTools or FromRegistry")
	}

	// Apply filters (AND composition)
	for _, filter := range b.filters {
		var filtered []*tooladapter.CanonicalTool
		for _, t := range tools {
			if filter(t) {
				filtered = append(filtered, t)
			}
		}
		tools = filtered
	}

	// Apply policy (last)
	if b.policy != nil {
		var allowed []*tooladapter.CanonicalTool
		for _, t := range tools {
			if b.policy.Allow(t) {
				allowed = append(allowed, t)
			}
		}
		tools = allowed
	}

	// Build toolset
	ts := New(b.name)
	for _, t := range tools {
		ts.Add(t)
	}
	return ts, nil
}
