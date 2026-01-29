package toolset

import "github.com/jonwraymond/tooladapter"

// Policy decides whether a tool is allowed.
type Policy interface {
	Allow(tool *tooladapter.CanonicalTool) bool
}
