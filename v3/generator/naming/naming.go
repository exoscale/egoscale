// Package naming converts OpenAPI names and descriptions to Go identifiers and doc comments.
package naming

import (
	"regexp"
	"strings"

	abbr "github.com/BluntSporks/abbreviation"
)

const separators = `[-_./+\s]`

var (
	trimRe  = regexp.MustCompile("^" + separators + "+|" + separators + "+$")
	splitRe = regexp.MustCompile(separators)
)

// Namer converts names to Go identifiers. Words found in the generic
// acronyms list of github.com/BluntSporks/abbreviation are uppercased,
// then its own acronyms (lowercase word -> Go spelling) take precedence.
type Namer struct {
	acronyms map[string]string
}

// New returns a Namer using the given additional acronyms.
func New(acronyms map[string]string) *Namer {
	return &Namer{acronyms: acronyms}
}

// Camel converts a string to CamelCase.
func (n *Namer) Camel(s string) string {
	return n.toInitialCamel(s, false)
}

// LowerCamel converts a string to lowerCamelCase.
func (n *Namer) LowerCamel(s string) string {
	return n.toInitialCamel(s, true)
}

// toInitialCamel got inspiration from https://github.com/iancoleman/strcase
// with improvement on acronym conversion.
func (n *Namer) toInitialCamel(s string, lower bool) string {
	words := splitRe.Split(trimRe.ReplaceAllString(s, ""), -1)
	words = slicesDeleteEmpty(words)

	for i, w := range words {
		if _, ok := abbr.Acronyms[strings.ToUpper(w)]; ok {
			w = strings.ToUpper(w)
		}
		if a, ok := n.acronyms[words[i]]; ok {
			w = a
		}

		if i == 0 && lower {
			words[i] = strings.ToLower(w)
			continue
		}

		if c := w[0]; c >= 'a' && c <= 'z' {
			w = string(c-'a'+'A') + w[1:]
		}
		words[i] = w
	}

	return strings.Join(words, "")
}

func slicesDeleteEmpty(words []string) []string {
	result := words[:0]
	for _, w := range words {
		if w != "" {
			result = append(result, w)
		}
	}

	return result
}

// Doc returns a Go doc comment from an OpenAPI description, "" if empty.
func Doc(doc string) string {
	if doc == "null" {
		return ""
	}

	var lines []string
	for _, l := range strings.Split(doc, "\n") {
		if l == "" {
			continue
		}
		lines = append(lines, "// "+strings.TrimSpace(l))
	}

	return strings.Join(lines, "\n")
}
