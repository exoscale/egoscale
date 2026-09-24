// Package config holds the generator settings that are not part of the OpenAPI spec.
package config

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

// Config is the generator configuration, loaded from config.yaml.
type Config struct {
	// Acronyms are Exoscale specific acronyms (lowercase word -> Go spelling),
	// applied on top of the generic list of github.com/BluntSporks/abbreviation.
	Acronyms map[string]string `yaml:"acronyms"`
	// IgnoredSchemas are component schemas not rendered.
	IgnoredSchemas []string `yaml:"ignored-schemas"`
	// SkipAuthOperations are operation IDs called without credentials.
	SkipAuthOperations []string `yaml:"skip-auth-operations"`
	// SchemaOverrides are per Go type name property and reference overrides.
	SchemaOverrides map[string]Overrides `yaml:"schema-overrides"`
	// Aliases are type aliases kept for backward compatibility.
	Aliases []Alias `yaml:"aliases"`
}

// Overrides of a schema rendering.
type Overrides struct {
	// Props renames properties (e.g. "instance-target" -> "instance").
	Props map[string]string `yaml:"props"`
	// Refs replaces referenced types (e.g. "#/components/schemas/instance-ref" -> "InstanceTarget").
	Refs map[string]string `yaml:"refs"`
}

// Alias is a Go type alias: type Name = Target.
type Alias struct {
	Name   string `yaml:"name"`
	Target string `yaml:"target"`
}

// Parse parses a YAML configuration, rejecting unknown fields.
func Parse(buf []byte) (Config, error) {
	var c Config
	dec := yaml.NewDecoder(bytes.NewReader(buf))
	dec.KnownFields(true)
	if err := dec.Decode(&c); err != nil {
		return Config{}, fmt.Errorf("config: %w", err)
	}

	return c, nil
}

// PropName returns the overridden name of a schema property, or propName.
func (c Config) PropName(schemaName, propName string) string {
	if name, ok := c.SchemaOverrides[schemaName].Props[propName]; ok {
		return name
	}

	return propName
}

// RefType returns the overridden Go type of a reference within a schema, if any.
func (c Config) RefType(schemaName, ref string) (string, bool) {
	typ, ok := c.SchemaOverrides[schemaName].Refs[ref]
	return typ, ok
}
