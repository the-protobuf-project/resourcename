package resourcename

import (
	"fmt"
	"regexp"
	"strings"
)

// Template provides simple placeholder-based resource name handling.
type Template struct {
	template     string
	regex        *regexp.Regexp
	placeholders []string
}

// New constructs a Template from a template string like
// "//domain.com/Thing/{id}/{sub}".
// bubuild(creates a template from a template string. Kept private as requested.
func build(template string) (*Template, error) {
	placeholderRegex := regexp.MustCompile(`\{([^{}]+)\}`)
	matches := placeholderRegex.FindAllStringSubmatch(template, -1)

	placeholders := make([]string, 0, len(matches))
	for _, m := range matches {
		placeholders = append(placeholders, m[1])
	}

	pattern := regexp.QuoteMeta(template)
	for _, p := range placeholders {
		pattern = strings.Replace(pattern, regexp.QuoteMeta("{"+p+"}"), `([^/]+)`, 1)
	}
	pattern = "^" + pattern + "$"

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid template pattern: %w", err)
	}

	return &Template{template: template, regex: re, placeholders: placeholders}, nil
}

// Generate builds a resource string from values. All placeholders must be provided.
func (t *Template) Generate(values map[string]string) (string, error) {
	res := t.template
	for _, ph := range t.placeholders {
		v, ok := values[ph]
		if !ok {
			return "", fmt.Errorf("missing value for placeholder %q", ph)
		}
		if strings.Contains(v, "/") {
			return "", fmt.Errorf("value for %q contains '/'", ph)
		}
		res = strings.Replace(res, "{"+ph+"}", v, 1)
	}
	return res, nil
}

// Parse extracts placeholders from a resource string.
func (t *Template) Parse(resource string) (map[string]string, error) {
	m := t.regex.FindStringSubmatch(resource)
	if m == nil {
		return nil, fmt.Errorf("resource %q does not match template %q", resource, t.template)
	}
	if len(m)-1 != len(t.placeholders) {
		return nil, fmt.Errorf("unexpected number of matches")
	}
	out := make(map[string]string, len(t.placeholders))
	for i, ph := range t.placeholders {
		out[ph] = m[i+1]
	}
	return out, nil
}
