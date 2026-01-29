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
		return nil, nil, []error{errors.New("adapter is nil")}
	}

	tools := e.toolset.Tools()
	result := make([]any, 0, len(tools))
	var warnings []tooladapter.FeatureLossWarning
	var errs []error

	for _, t := range tools {
		sourceName := t.SourceFormat
		if sourceName == "" {
			sourceName = "canonical"
		}

		warnings = append(warnings, detectSchemaFeatureLoss(t.InputSchema, sourceName, e.adapter)...)
		warnings = append(warnings, detectSchemaFeatureLoss(t.OutputSchema, sourceName, e.adapter)...)

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

// detectSchemaFeatureLoss checks which features in a schema are not supported.
func detectSchemaFeatureLoss(schema *tooladapter.JSONSchema, sourceName string, adapter tooladapter.Adapter) []tooladapter.FeatureLossWarning {
	if schema == nil {
		return nil
	}

	featureUsage := map[tooladapter.SchemaFeature]bool{
		tooladapter.FeatureRef:                  schema.Ref != "",
		tooladapter.FeatureDefs:                 len(schema.Defs) > 0,
		tooladapter.FeatureAnyOf:                len(schema.AnyOf) > 0,
		tooladapter.FeatureOneOf:                len(schema.OneOf) > 0,
		tooladapter.FeatureAllOf:                len(schema.AllOf) > 0,
		tooladapter.FeatureNot:                  schema.Not != nil,
		tooladapter.FeaturePattern:              schema.Pattern != "",
		tooladapter.FeatureFormat:               schema.Format != "",
		tooladapter.FeatureAdditionalProperties: schema.AdditionalProperties != nil,
		tooladapter.FeatureMinimum:              schema.Minimum != nil,
		tooladapter.FeatureMaximum:              schema.Maximum != nil,
		tooladapter.FeatureMinLength:            schema.MinLength != nil,
		tooladapter.FeatureMaxLength:            schema.MaxLength != nil,
		tooladapter.FeatureEnum:                 len(schema.Enum) > 0,
		tooladapter.FeatureConst:                schema.Const != nil,
		tooladapter.FeatureDefault:              schema.Default != nil,
	}

	var warnings []tooladapter.FeatureLossWarning
	for feature, used := range featureUsage {
		if used && !adapter.SupportsFeature(feature) {
			warnings = append(warnings, tooladapter.FeatureLossWarning{
				Feature:     feature,
				FromAdapter: sourceName,
				ToAdapter:   adapter.Name(),
			})
		}
	}

	if schema.Properties != nil {
		for _, prop := range schema.Properties {
			warnings = append(warnings, detectSchemaFeatureLoss(prop, sourceName, adapter)...)
		}
	}
	if schema.Items != nil {
		warnings = append(warnings, detectSchemaFeatureLoss(schema.Items, sourceName, adapter)...)
	}
	if schema.Defs != nil {
		for _, def := range schema.Defs {
			warnings = append(warnings, detectSchemaFeatureLoss(def, sourceName, adapter)...)
		}
	}
	for _, s := range schema.AnyOf {
		warnings = append(warnings, detectSchemaFeatureLoss(s, sourceName, adapter)...)
	}
	for _, s := range schema.OneOf {
		warnings = append(warnings, detectSchemaFeatureLoss(s, sourceName, adapter)...)
	}
	for _, s := range schema.AllOf {
		warnings = append(warnings, detectSchemaFeatureLoss(s, sourceName, adapter)...)
	}
	if schema.Not != nil {
		warnings = append(warnings, detectSchemaFeatureLoss(schema.Not, sourceName, adapter)...)
	}

	return warnings
}
