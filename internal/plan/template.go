package plan

import (
	"encoding/json"
	"fmt"
	"strings"
)

var templateSections = []string{
	"Objective",
	"Context",
	"Research",
	"Plan",
	"Outputs",
	"Verification",
	"Review",
	"Notes",
}

func Template(meta Frontmatter) string {
	var b strings.Builder
	// Render frontmatter explicitly to keep field order and empty-value handling
	// deterministic for git review and tests.
	b.WriteString("---\n")
	fmt.Fprintf(&b, "id: %s\n", meta.ID)
	fmt.Fprintf(&b, "title: %s\n", yamlStringLiteral(meta.Title))
	fmt.Fprintf(&b, "status: %s\n", meta.Status)
	if len(meta.Tags) == 0 {
		b.WriteString("tags: []\n")
	} else {
		b.WriteString("tags:\n")
		for _, tag := range NormalizeTags(meta.Tags) {
			fmt.Fprintf(&b, "  - %s\n", tag)
		}
	}
	fmt.Fprintf(&b, "parent: %s\n", yamlStringLiteral(meta.Parent))
	b.WriteString("---\n\n")

	for i, section := range templateSections {
		fmt.Fprintf(&b, "## %s\n\n", section)
		if i != len(templateSections)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func yamlStringLiteral(value string) string {
	encoded, err := json.Marshal(value)
	if err != nil {
		return `""`
	}
	return string(encoded)
}
