package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"openplan/internal/plan"
	"openplan/internal/repo"
)

func TestSkillCommandMatchesGolden(t *testing.T) {
	output, err := runCommand(t, t.TempDir(), "skill")
	if err != nil {
		t.Fatalf("skill command returned error: %v", err)
	}

	golden, err := os.ReadFile(filepath.Join("testdata", "skill.golden.md"))
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}

	if output != string(golden) {
		t.Fatalf("skill output mismatch\nwant:\n%s\ngot:\n%s", string(golden), output)
	}
}

func TestInitListAndShowJSON(t *testing.T) {
	root := t.TempDir()
	output, err := runCommand(t, root, "init", "--prefix", "OPN")
	if err != nil {
		t.Fatalf("init command returned error: %v", err)
	}
	if !strings.Contains(output, "Initialized OpenPlan") {
		t.Fatalf("unexpected init output: %q", output)
	}

	ctx, err := repo.Load(root)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	content := plan.Template(plan.Frontmatter{
		ID:     "OPN-1_ABCDEFGH",
		Title:  "Draft README",
		Status: plan.StatusPlan,
		Tags:   []string{"docs"},
	})
	if err := os.WriteFile(ctx.PlanPath("OPN-1_ABCDEFGH"), []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	listOutput, err := runCommand(t, root, "list", "--json")
	if err != nil {
		t.Fatalf("list command returned error: %v", err)
	}
	var listed struct {
		Plans []struct {
			ID string `json:"id"`
		} `json:"plans"`
	}
	if err := json.Unmarshal([]byte(listOutput), &listed); err != nil {
		t.Fatalf("json.Unmarshal list output returned error: %v", err)
	}
	if len(listed.Plans) != 1 || listed.Plans[0].ID != "OPN-1_ABCDEFGH" {
		t.Fatalf("unexpected list json: %s", listOutput)
	}

	showOutput, err := runCommand(t, root, "show", "OPN-1_ABCDEFGH", "--json")
	if err != nil {
		t.Fatalf("show command returned error: %v", err)
	}
	var shown struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal([]byte(showOutput), &shown); err != nil {
		t.Fatalf("json.Unmarshal show output returned error: %v", err)
	}
	if shown.ID != "OPN-1_ABCDEFGH" {
		t.Fatalf("unexpected show json: %s", showOutput)
	}
}

func TestValidateReportsInvalidPlans(t *testing.T) {
	root := t.TempDir()
	if _, err := runCommand(t, root, "init", "--prefix", "OPN"); err != nil {
		t.Fatalf("init command returned error: %v", err)
	}

	ctx, err := repo.Load(root)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if err := os.WriteFile(ctx.PlanPath("OPN-1_ABCDEFGH"), []byte("---\nid: wrong\n---\n"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	output, err := runCommand(t, root, "validate")
	if err == nil {
		t.Fatalf("validate command returned nil error, want failure")
	}
	if output != "" {
		t.Fatalf("validate stdout = %q, want empty string on failure", output)
	}
	if !strings.Contains(err.Error(), "validation failed") {
		t.Fatalf("validate error = %q, want validation header", err.Error())
	}

	output, err = runCommand(t, root, "list")
	if err != nil {
		t.Fatalf("list command returned error: %v", err)
	}
	if !strings.Contains(output, "Issues:") {
		t.Fatalf("expected list output to include issues, got %q", output)
	}
}

// runCommand executes the CLI against a temporary working directory while
// capturing stdout for assertions
func runCommand(t *testing.T, cwd string, args ...string) (string, error) {
	t.Helper()

	oldStdout := stdout
	defer func() { stdout = oldStdout }()

	var out bytes.Buffer
	stdout = &out

	oldWD, err := os.Getwd()
	if err != nil {
		return "", err
	}
	defer func() {
		_ = os.Chdir(oldWD)
	}()

	if err := os.Chdir(cwd); err != nil {
		return "", err
	}

	cmd := newRootCmd()
	cmd.SetArgs(args)
	err = cmd.Execute()
	return out.String(), err
}
