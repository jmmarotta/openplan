package plan

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestParseBytesValidateAndTemplate(t *testing.T) {
	path := filepath.Join(t.TempDir(), FilenameForID("OPN-1_ABCDEFGH"))
	content := Template(Frontmatter{
		ID:     "OPN-1_ABCDEFGH",
		Title:  "Draft README",
		Status: StatusPlan,
		Tags:   []string{"docs", "readme"},
		Parent: "",
	})

	doc, err := ParseBytes(path, []byte(content))
	if err != nil {
		t.Fatalf("ParseBytes returned error: %v", err)
	}

	if issues := Validate(doc); len(issues) != 0 {
		t.Fatalf("Validate() returned issues: %#v", issues)
	}

	for _, heading := range []string{"## Objective", "## Verification", "## Notes"} {
		if !strings.Contains(content, heading) {
			t.Fatalf("template missing section %q", heading)
		}
	}
}

func TestValidateInvalidFields(t *testing.T) {
	doc := Document{
		Path: filepath.Join(t.TempDir(), FilenameForID("OPN-2_ABCDEFGH")),
		Frontmatter: Frontmatter{
			ID:     "wrong-id",
			Title:  "",
			Status: Status("bad"),
			Tags:   []string{"Needs Review", "Needs Review"},
			Parent: "123",
		},
	}

	issues := Validate(doc)
	if len(issues) < 5 {
		t.Fatalf("Validate() returned %d issues, want at least 5", len(issues))
	}
}
