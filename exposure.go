package toolset

import (
	"errors"

	"github.com/jonwraymond/tooladapter"
)

// Exposure exports a Toolset to protocol-specific formats.
type Exposure struct {
	toolset *Toolset
	adapter tooladapter.Adapter
}

// NewExposure creates an Exposure for the given toolset and adapter.
func NewExposure(ts *Toolset, adapter tooladapter.Adapter) *Exposure {
	return &Exposure{toolset: ts, adapter: adapter}
}

// Export converts all tools to the adapter's format.
func (e *Exposure) Export() ([]any, error) {
	if e.adapter == nil {
		return nil, errors.New("adapter is nil")
	}
	tools := e.toolset.Tools()
	result := make([]any, 0, len(tools))
	for _, t := range tools {
		converted, err := e.adapter.FromCanonical(t)
		if err != nil {
			return nil, err
		}
		result = append(result, converted)
	}
	return result, nil
}

// ExportWithWarnings converts tools and returns feature loss warnings and conversion errors.
// Unlike Export, this method continues on conversion errors and collects them for reporting.
// Callers should check the errors slice to detect tools that failed to convert.
func (e *Exposure) ExportWithWarnings() ([]any, []tooladapter.FeatureLossWarning, []error) {
	if e.adapter == nil {
		return nil, nil, nil
	}

	tools := e.toolset.Tools()
	result := make([]any, 0, len(tools))
	var warnings []tooladapter.FeatureLossWarning
	var errs []error

	for _, t := range tools {
		// Detect features used by this tool
		features := detectFeatures(t)

		// Check for unsupported features
		for _, f := range features {
			if !e.adapter.SupportsFeature(f) {
				warnings = append(warnings, tooladapter.FeatureLossWarning{
					Feature:     f,
					FromAdapter: t.SourceFormat,
					ToAdapter:   e.adapter.Name(),
				})
			}
		}

		// Convert
		converted, err := e.adapter.FromCanonical(t)
		if err != nil {
			errs = append(errs, &ConversionError{
				ToolID: t.ID(),
				Cause:  err,
			})
			continue
		}
		result = append(result, converted)
	}
	return result, warnings, errs
}

// ConversionError represents a tool that failed to convert.
type ConversionError struct {
	ToolID string
	Cause  error
}

func (e *ConversionError) Error() string {
	return "failed to convert tool " + e.ToolID + ": " + e.Cause.Error()
}

func (e *ConversionError) Unwrap() error {
	return e.Cause
}

// detectFeatures recursively walks the schema to find all used features.
func detectFeatures(t *tooladapter.CanonicalTool) []tooladapter.SchemaFeature {
	var features []tooladapter.SchemaFeature
	seen := make(map[tooladapter.SchemaFeature]bool)

	var walk func(s *tooladapter.JSONSchema)
	walk = func(s *tooladapter.JSONSchema) {
		if s == nil {
			return
		}

		// Check each feature field
		if s.Ref != "" && !seen[tooladapter.FeatureRef] {
			features = append(features, tooladapter.FeatureRef)
			seen[tooladapter.FeatureRef] = true
		}
		if len(s.Defs) > 0 && !seen[tooladapter.FeatureDefs] {
			features = append(features, tooladapter.FeatureDefs)
			seen[tooladapter.FeatureDefs] = true
		}
		if len(s.AnyOf) > 0 && !seen[tooladapter.FeatureAnyOf] {
			features = append(features, tooladapter.FeatureAnyOf)
			seen[tooladapter.FeatureAnyOf] = true
		}
		if len(s.OneOf) > 0 && !seen[tooladapter.FeatureOneOf] {
			features = append(features, tooladapter.FeatureOneOf)
			seen[tooladapter.FeatureOneOf] = true
		}
		if len(s.AllOf) > 0 && !seen[tooladapter.FeatureAllOf] {
			features = append(features, tooladapter.FeatureAllOf)
			seen[tooladapter.FeatureAllOf] = true
		}
		if s.Not != nil && !seen[tooladapter.FeatureNot] {
			features = append(features, tooladapter.FeatureNot)
			seen[tooladapter.FeatureNot] = true
		}
		if s.Pattern != "" && !seen[tooladapter.FeaturePattern] {
			features = append(features, tooladapter.FeaturePattern)
			seen[tooladapter.FeaturePattern] = true
		}
		if s.Format != "" && !seen[tooladapter.FeatureFormat] {
			features = append(features, tooladapter.FeatureFormat)
			seen[tooladapter.FeatureFormat] = true
		}
		if s.AdditionalProperties != nil && !seen[tooladapter.FeatureAdditionalProperties] {
			features = append(features, tooladapter.FeatureAdditionalProperties)
			seen[tooladapter.FeatureAdditionalProperties] = true
		}
		if s.Minimum != nil && !seen[tooladapter.FeatureMinimum] {
			features = append(features, tooladapter.FeatureMinimum)
			seen[tooladapter.FeatureMinimum] = true
		}
		if s.Maximum != nil && !seen[tooladapter.FeatureMaximum] {
			features = append(features, tooladapter.FeatureMaximum)
			seen[tooladapter.FeatureMaximum] = true
		}
		if s.MinLength != nil && !seen[tooladapter.FeatureMinLength] {
			features = append(features, tooladapter.FeatureMinLength)
			seen[tooladapter.FeatureMinLength] = true
		}
		if s.MaxLength != nil && !seen[tooladapter.FeatureMaxLength] {
			features = append(features, tooladapter.FeatureMaxLength)
			seen[tooladapter.FeatureMaxLength] = true
		}
		if len(s.Enum) > 0 && !seen[tooladapter.FeatureEnum] {
			features = append(features, tooladapter.FeatureEnum)
			seen[tooladapter.FeatureEnum] = true
		}
		if s.Const != nil && !seen[tooladapter.FeatureConst] {
			features = append(features, tooladapter.FeatureConst)
			seen[tooladapter.FeatureConst] = true
		}
		if s.Default != nil && !seen[tooladapter.FeatureDefault] {
			features = append(features, tooladapter.FeatureDefault)
			seen[tooladapter.FeatureDefault] = true
		}

		// Recurse into nested schemas
		for _, prop := range s.Properties {
			walk(prop)
		}
		walk(s.Items)
		for _, def := range s.Defs {
			walk(def)
		}
		for _, schema := range s.AnyOf {
			walk(schema)
		}
		for _, schema := range s.OneOf {
			walk(schema)
		}
		for _, schema := range s.AllOf {
			walk(schema)
		}
		walk(s.Not)
	}

	walk(t.InputSchema)
	walk(t.OutputSchema)
	return features
}
