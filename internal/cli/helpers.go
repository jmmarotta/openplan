package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/jmmarotta/openplan/internal/plan"
	"github.com/jmmarotta/openplan/internal/repo"
)

var stdout io.Writer = os.Stdout

// EditorUnsetError reports that commands requiring an editor cannot continue
// because `$EDITOR` is unset.
type EditorUnsetError struct{}

func (EditorUnsetError) Error() string {
	return "$EDITOR is not set"
}

// requireRepo discovers the nearest OpenPlan repository from the current
// working directory.
func requireRepo() (repo.Context, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return repo.Context{}, err
	}
	return repo.Discover(cwd)
}

// openEditor delegates to the user-selected shell command so editor arguments
// and wrappers continue to work as configured in `$EDITOR`.
func openEditor(path string) error {
	if _, err := editorCommand(); err != nil {
		return err
	}

	cmd := exec.Command("/bin/sh", "-c", `$EDITOR "$1"`, "openplan", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	return cmd.Run()
}

// writeStdout centralizes command output so tests can replace the destination.
func writeStdout(s string) error {
	_, err := io.WriteString(stdout, s)
	return err
}

// writeJSON keeps machine output consistently indented and newline-terminated.
func writeJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = io.Copy(stdout, bytes.NewReader(data))
	return err
}

// editorCommand validates the editor contract once so command handlers can fail
// before mutating repository state.
func editorCommand() (string, error) {
	editor := strings.TrimSpace(os.Getenv("EDITOR"))
	if editor == "" {
		return "", EditorUnsetError{}
	}
	return editor, nil
}

// formatIssues turns validation issues into a single user-facing error string
// so CLI presentation stays outside the lower-level packages.
func formatIssues(path string, issues []plan.ValidationIssue) error {
	var b strings.Builder
	fmt.Fprintf(&b, "invalid plan: %s\n", path)
	for _, issue := range issues {
		fmt.Fprintf(&b, "- %s: %s\n", issue.Field, issue.Message)
	}
	return errors.New(strings.TrimSpace(b.String()))
}

// formatValidationIssues renders repository-wide validation failures as one
// deterministic error so callers get a single non-zero result with full detail.
func formatValidationIssues(issues []plan.ValidationIssue) error {
	var b strings.Builder
	b.WriteString("validation failed\n")
	for _, issue := range issues {
		fmt.Fprintf(&b, "- %s: %s: %s\n", issue.Path, issue.Field, issue.Message)
	}
	return errors.New(strings.TrimSpace(b.String()))
}
