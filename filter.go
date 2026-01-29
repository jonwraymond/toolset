package toolset

import "github.com/jonwraymond/tooladapter"

// NamespaceFilter returns a filter matching tools in any of the namespaces.
func NamespaceFilter(namespaces ...string) FilterFunc {
	set := make(map[string]bool, len(namespaces))
	for _, ns := range namespaces {
		set[ns] = true
	}
	return func(t *tooladapter.CanonicalTool) bool {
		if t == nil {
			return false
		}
		return set[t.Namespace]
	}
}

// TagsAny returns a filter matching tools with ANY of the tags.
func TagsAny(tags ...string) FilterFunc {
	set := make(map[string]bool, len(tags))
	for _, tag := range tags {
		set[tag] = true
	}
	return func(t *tooladapter.CanonicalTool) bool {
		if t == nil {
			return false
		}
		for _, tag := range t.Tags {
			if set[tag] {
				return true
			}
		}
		return false
	}
}

// TagsAll returns a filter matching tools with ALL of the tags.
func TagsAll(tags ...string) FilterFunc {
	return func(t *tooladapter.CanonicalTool) bool {
		if t == nil {
			return false
		}
		if len(tags) == 0 {
			return true // vacuously true
		}
		toolTags := make(map[string]bool, len(t.Tags))
		for _, tag := range t.Tags {
			toolTags[tag] = true
		}
		for _, required := range tags {
			if !toolTags[required] {
				return false
			}
		}
		return true
	}
}

// TagsNone returns a filter matching tools with NONE of the tags.
func TagsNone(tags ...string) FilterFunc {
	set := make(map[string]bool, len(tags))
	for _, tag := range tags {
		set[tag] = true
	}
	return func(t *tooladapter.CanonicalTool) bool {
		if t == nil {
			return false
		}
		for _, tag := range t.Tags {
			if set[tag] {
				return false
			}
		}
		return true
	}
}

// CategoryFilter returns a filter matching tools in any of the categories.
func CategoryFilter(categories ...string) FilterFunc {
	set := make(map[string]bool, len(categories))
	for _, cat := range categories {
		set[cat] = true
	}
	return func(t *tooladapter.CanonicalTool) bool {
		if t == nil {
			return false
		}
		return set[t.Category]
	}
}

// AllowIDs returns a filter matching only the listed tool IDs.
func AllowIDs(ids ...string) FilterFunc {
	set := make(map[string]bool, len(ids))
	for _, id := range ids {
		set[id] = true
	}
	return func(t *tooladapter.CanonicalTool) bool {
		if t == nil {
			return false
		}
		return set[t.ID()]
	}
}

// DenyIDs returns a filter excluding the listed tool IDs.
func DenyIDs(ids ...string) FilterFunc {
	set := make(map[string]bool, len(ids))
	for _, id := range ids {
		set[id] = true
	}
	return func(t *tooladapter.CanonicalTool) bool {
		if t == nil {
			return false
		}
		return !set[t.ID()]
	}
}
