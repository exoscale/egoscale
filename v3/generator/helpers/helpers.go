package helpers

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	abbr "github.com/BluntSporks/abbreviation"
)

var uppercaseAcronym = sync.Map{}

// ConfigureAcronym allows you to add additional words which will be considered acronyms
func ConfigureAcronym(key, val string) {
	uppercaseAcronym.Store(key, val)
}

type Overrides struct {
	Props map[string]string // Property name overrides (e.g., "instance-target" -> "instance")
	Refs  map[string]string // Reference path overrides (e.g., "#/components/schemas/instance-ref" -> "InstanceTarget")
}

var SchemaPropertyOverrides = map[string]*Overrides{
	"CreateInstance": {
		Props: nil,
		Refs: map[string]string{
			"#/components/schemas/anti-affinity-group-ref": "AntiAffinityGroup",
			"#/components/schemas/security-group-ref":      "SecurityGroup",
			"#/components/schemas/template-ref":            "Template",
			"#/components/schemas/ssh-key-ref":             "SSHKey",
			"#/components/schemas/deploy-target-ref":       "DeployTarget",
			"#/components/schemas/instance-type-ref":       "InstanceType",
		},
	},
	"ScaleInstance": {
		Props: nil,
		Refs: map[string]string{
			"#/components/schemas/instance-type-ref": "InstanceType",
		},
	},
	"ResetInstance": {
		Props: nil,
		Refs: map[string]string{
			"#/components/schemas/template-ref": "Template",
		},
	},
	"AttachInstanceToElasticIPRequest": {
		Props: map[string]string{
			"instance-target": "instance",
		},
		Refs: map[string]string{
			"#/components/schemas/instance-ref": "InstanceTarget",
		},
	},
	"DetachInstanceFromElasticIPRequest": {
		Props: map[string]string{
			"instance-target": "instance",
		},
		Refs: map[string]string{
			"#/components/schemas/instance-ref": "InstanceTarget",
		},
	},
	"CreateBlockStorageVolumeRequest": {
		Props: nil,
		Refs: map[string]string{
			"#/components/schemas/block-storage-snapshot-ref": "BlockStorageSnapshotTarget",
		},
	},
	"AttachBlockStorageVolumeToInstanceRequest": {
		Props: map[string]string{
			"block-storage-volume-ref": "block-storage-volume",
		},
		Refs: map[string]string{
			"#/components/schemas/instance-ref": "InstanceTarget",
		},
	},
	"CreateInstanceRequest": {
		Props: nil,
		Refs: map[string]string{
			"#/components/schemas/anti-affinity-group-ref": "AntiAffinityGroup",
			"#/components/schemas/security-group-ref":      "SecurityGroup",
			"#/components/schemas/template-ref":            "Template",
			"#/components/schemas/ssh-key-ref":             "SSHKey",
			"#/components/schemas/deploy-target-ref":       "DeployTarget",
			"#/components/schemas/instance-type-ref":       "InstanceType",
		},
	},
	"ResetInstanceRequest": {
		Props: nil,
		Refs: map[string]string{
			"#/components/schemas/template-ref": "Template",
		},
	},
	"ScaleInstanceRequest": {
		Props: nil,
		Refs: map[string]string{
			"#/components/schemas/instance-type-ref": "InstanceType",
		},
	},
	"BlockStorageSnapshot": {
		Props: nil,
		Refs: map[string]string{
			"#/components/schemas/block-storage-volume-ref": "BlockStorageVolumeTarget",
		},
	},
	"BlockStorageVolume": {
		Props: nil,
		Refs: map[string]string{
			"#/components/schemas/instance-ref":               "InstanceTarget",
			"#/components/schemas/block-storage-snapshot-ref": "BlockStorageSnapshotTarget",
		},
	},
	"CreateDBAASServiceMysqlRequest": {
		Props: nil,
		Refs: map[string]string{
			"#/components/schemas/dbaas-mysql-user-password": "string",
		},
	},
	"ResetDBAASMysqlUserPasswordRequest": {
		Props: nil,
		Refs: map[string]string{
			"#/components/schemas/dbaas-mysql-user-password": "DBAASUserPassword",
		},
	},
	"CreateInstancePoolRequest": {
		Props: nil,
		Refs: map[string]string{
			"#/components/schemas/anti-affinity-group-ref": "AntiAffinityGroup",
			"#/components/schemas/deploy-target-ref":       "DeployTarget",
			"#/components/schemas/elastic-ip-ref":          "ElasticIP",
			"#/components/schemas/instance-type-ref":       "InstanceType",
			"#/components/schemas/private-network-ref":     "PrivateNetwork",
			"#/components/schemas/security-group-ref":      "SecurityGroup",
			"#/components/schemas/template-ref":            "Template",
			"#/components/schemas/ssh-key-ref":             "SSHKey",
		},
	},
	"UpdateInstancePoolRequest": {
		Props: nil,
		Refs: map[string]string{
			"#/components/schemas/anti-affinity-group-ref": "AntiAffinityGroup",
			"#/components/schemas/deploy-target-ref":       "DeployTarget",
			"#/components/schemas/elastic-ip-ref":          "ElasticIP",
			"#/components/schemas/instance-type-ref":       "InstanceType",
			"#/components/schemas/private-network-ref":     "PrivateNetwork",
			"#/components/schemas/security-group-ref":      "SecurityGroup",
			"#/components/schemas/template-ref":            "Template",
			"#/components/schemas/ssh-key-ref":             "SSHKey",
		},
	},
	"CreateSKSNodepoolRequest": {
		Props: nil,
		Refs: map[string]string{
			"#/components/schemas/anti-affinity-group-ref": "AntiAffinityGroup",
			"#/components/schemas/deploy-target-ref":       "DeployTarget",
			"#/components/schemas/instance-type-ref":       "InstanceType",
			"#/components/schemas/private-network-ref":     "PrivateNetwork",
			"#/components/schemas/security-group-ref":      "SecurityGroup",
		},
	},
	"UpdateSKSNodepoolRequest": {
		Props: nil,
		Refs: map[string]string{
			"#/components/schemas/anti-affinity-group-ref": "AntiAffinityGroup",
			"#/components/schemas/deploy-target-ref":       "DeployTarget",
			"#/components/schemas/instance-type-ref":       "InstanceType",
			"#/components/schemas/private-network-ref":     "PrivateNetwork",
			"#/components/schemas/security-group-ref":      "SecurityGroup",
		},
	},
	"Instance": {
		Props: nil,
		Refs: map[string]string{
			"#/components/schemas/anti-affinity-group-ref": "AntiAffinityGroup",
			"#/components/schemas/deploy-target-ref":       "DeployTarget",
			"#/components/schemas/elastic-ip-ref":          "ElasticIP",
			"#/components/schemas/security-group-ref":      "SecurityGroup",
			"#/components/schemas/snapshot-ref":            "Snapshot",
		},
	},
	"InstancePool": {
		Props: nil,
		Refs: map[string]string{
			"#/components/schemas/anti-affinity-group-ref": "AntiAffinityGroup",
			"#/components/schemas/deploy-target-ref":       "DeployTarget",
			"#/components/schemas/elastic-ip-ref":          "ElasticIP",
			"#/components/schemas/instance-type-ref":       "InstanceType",
			"#/components/schemas/instance-ref":            "Instance",
			"#/components/schemas/private-network-ref":     "PrivateNetwork",
			"#/components/schemas/security-group-ref":      "SecurityGroup",
			"#/components/schemas/ssh-key-ref":             "SSHKey",
			"#/components/schemas/template-ref":            "Template",
		},
	},
	"SKSNodepool": {
		Props: nil,
		Refs: map[string]string{
			"#/components/schemas/anti-affinity-group-ref": "AntiAffinityGroup",
			"#/components/schemas/deploy-target-ref":       "DeployTarget",
			"#/components/schemas/instance-pool-ref":       "InstancePool",
			"#/components/schemas/instance-type-ref":       "InstanceType",
			"#/components/schemas/private-network-ref":     "PrivateNetwork",
			"#/components/schemas/security-group-ref":      "SecurityGroup",
			"#/components/schemas/template-ref":            "Template",
		},
	},
}

var SpecialAliases = []string{
	"type InstanceTarget = InstanceRef",
	"type BlockStorageSnapshotTarget = BlockStorageSnapshotRef",
	"type BlockStorageVolumeTarget = BlockStorageVolumeRef",
}

// RenderReference renders OpenAPI reference from path to go style.
func RenderReference(referencePath string, schemaName string) string {
	if overrides := SchemaPropertyOverrides[schemaName]; overrides != nil && overrides.Refs != nil {
		if override, ok := overrides.Refs[referencePath]; ok {
			return override
		}
	}
	return ToCamel(filepath.Base(referencePath))
}

// Header retruns header file for generated go source files.
func Header(packageName, version string) []byte {
	return []byte(fmt.Sprintf(`// Package %s provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/egoscale/v3/generator version %s DO NOT EDIT.
`,
		packageName, version,
	))
}

// ToLowerCamel converts a string to lowerCamelCase
func ToLowerCamel(s string) string {
	return toInitialCamel(splitWords(s, separatorRe), true)
}

// ToCamel converts a string to CamelCase
func ToCamel(s string) string {
	return toInitialCamel(splitWords(s, separatorRe), false)
}

// EnumValueName converts an OpenAPI enum value to a CamelCase Go
// identifier while keeping separator characters distinguishable, so
// values like "1g.24gb-me" and "1g.24gb+me" do not collide on
// "1g24gbMe". Separators become words: '-' -> "Minus", '+' -> "Plus",
// '.' -> "Dot", '/' -> "Slash". '_' and whitespace stay plain
// word boundaries.
//
// Use this only for enum-value identifiers; for operation IDs, query
// parameters, and property names, prefer ToCamel/ToLowerCamel to avoid
// renaming the public client API.
func EnumValueName(s string) string {
	return toInitialCamel(splitEnumWords(s), false)
}

// separatorRe matches characters treated as word boundaries by ToCamel
// and ToLowerCamel (operation IDs, property names, query params, ...).
var separatorRe = `[-_./+\s]`

func toInitialCamel(words []string, lower bool) string {
	for i, w := range words {
		if w == "" {
			continue
		}

		_, ok := abbr.Acronyms[strings.ToUpper(w)]
		if ok {
			words[i] = strings.ToUpper(w)
		}

		a, hasAcronym := uppercaseAcronym.Load(w)
		if hasAcronym {
			words[i] = a.(string)
		}

		if i == 0 && lower {
			words[i] = strings.ToLower(words[i])
			continue
		}

		bytes := []byte(words[i])
		v := bytes[0]
		vIsLow := v >= 'a' && v <= 'z'
		if vIsLow {
			bytes[0] += 'A'
			bytes[0] -= 'a'
		}

		words[i] = string(bytes)
	}

	return strings.Join(words, "")
}

// splitWords splits s on any character in separators, drops empty
// fragments, and trims leading/trailing separators.
func splitWords(s, separators string) []string {
	if s == "" {
		return nil
	}
	trimRe := regexp.MustCompile(fmt.Sprintf("^%s+|%s+$", separators, separators))
	s = trimRe.ReplaceAllString(s, "")

	re := regexp.MustCompile(separators)
	parts := re.Split(s, -1)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// enumSeparatorMap maps each distinguishable separator to a wrapped
// word inside enum values: '-' -> "_Minus_", '+' -> "_Plus_", ... The
// underscores act as word boundaries so splitEnumWords can hand each
// fragment to the existing CamelCase pipeline and get a properly
// capitalised identifier (e.g. "read-replica" -> "ReadMinusReplica").
var enumSeparatorMap = map[string]string{
	"-": "_Minus_",
	"+": "_Plus_",
	".": "_Dot_",
	"/": "_Slash_",
}

// enumSeparatorRe matches any of the separator characters
// distinguishable inside enum values.
var enumSeparatorRe = regexp.MustCompile(`[-+./]`)

// enumWordBoundaryRe matches runs of underscores or whitespace, used
// after enum-separator expansion to recover word boundaries.
var enumWordBoundaryRe = regexp.MustCompile(`[\s_]+`)

// enumTrimRe matches leading or trailing whitespace/underscores/separator
// characters so that "---trim---" collapses to "trim" before expansion.
var enumTrimRe = regexp.MustCompile(`^[\s_\-+./]+|[\s_\-+./]+$`)

func splitEnumWords(s string) []string {
	if s == "" {
		return nil
	}

	s = enumTrimRe.ReplaceAllString(s, "")

	expanded := enumSeparatorRe.ReplaceAllStringFunc(s, func(m string) string {
		return enumSeparatorMap[m]
	})

	parts := enumWordBoundaryRe.Split(expanded, -1)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// RenderDoc returns proper go doc comment from
// an OpenAPI spec field documentation.
func RenderDoc(doc string) string {
	if doc == "null" {
		return ""
	}

	docs := strings.Split(doc, "\n")
	r := []string{}
	for i, d := range docs {
		if d == "" {
			docs = append(docs[:i], docs[i+1:]...)
			continue
		}
		r = append(r, "// "+strings.TrimSpace(d))
	}

	if len(r) == 0 {
		return ""
	}

	return strings.Join(r, "\n")
}
