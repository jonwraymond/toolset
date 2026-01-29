package toolset

import "github.com/jonwraymond/tooladapter"

// Policy decides whether a tool is allowed.
type Policy interface {
	Allow(tool *tooladapter.CanonicalTool) bool
}

// PolicyFunc adapts a function to the Policy interface.
type PolicyFunc func(*tooladapter.CanonicalTool) bool

// Allow implements Policy.
func (f PolicyFunc) Allow(t *tooladapter.CanonicalTool) bool {
	return f(t)
}

// AllowAll returns a policy that allows all tools.
func AllowAll() Policy {
	return PolicyFunc(func(t *tooladapter.CanonicalTool) bool {
		return true
	})
}

// DenyAll returns a policy that denies all tools.
func DenyAll() Policy {
	return PolicyFunc(func(t *tooladapter.CanonicalTool) bool {
		return false
	})
}

// AllowNamespaces returns a policy allowing only listed namespaces.
func AllowNamespaces(ns ...string) Policy {
	set := make(map[string]bool, len(ns))
	for _, n := range ns {
		set[n] = true
	}
	return PolicyFunc(func(t *tooladapter.CanonicalTool) bool {
		if t == nil {
			return false
		}
		return set[t.Namespace]
	})
}

// DenyTags returns a policy denying tools with any of the tags.
func DenyTags(tags ...string) Policy {
	set := make(map[string]bool, len(tags))
	for _, tag := range tags {
		set[tag] = true
	}
	return PolicyFunc(func(t *tooladapter.CanonicalTool) bool {
		if t == nil {
			return false
		}
		for _, tag := range t.Tags {
			if set[tag] {
				return false
			}
		}
		return true
	})
}

// AllowScopes returns a policy allowing tools requiring only allowed scopes.
func AllowScopes(allowed ...string) Policy {
	set := make(map[string]bool, len(allowed))
	for _, scope := range allowed {
		set[scope] = true
	}
	return PolicyFunc(func(t *tooladapter.CanonicalTool) bool {
		if t == nil {
			return false
		}
		// Tools with no required scopes are allowed
		for _, scope := range t.RequiredScopes {
			if !set[scope] {
				return false
			}
		}
		return true
	})
}
