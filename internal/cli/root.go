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
	"text/tabwriter"

	"github.com/jmmarotta/openplan/internal/plan"
	"github.com/jmmarotta/openplan/internal/repo"
	"github.com/jmmarotta/openplan/internal/store"

	"github.com/spf13/cobra"
)

var stdout io.Writer = os.Stdout

// EditorUnsetError reports that commands requiring an editor cannot continue
// because `$EDITOR` is unset
type EditorUnsetError struct{}

func (EditorUnsetError) Error() string {
	return "$EDITOR is not set"
}

// Execute builds and runs the root command so `main` stays a thin adapter
func Execute() error {
	return newRootCmd().Execute()
}

// newRootCmd wires the full CLI surface in one place so the package-level
// command layout stays easy to audit
func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "openplan",
		Short:         "Filesystem-native planning for technical work",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(
		newInitCmd(),
		newNewCmd(),
		newEditCmd(),
		newListCmd(),
		newShowCmd(),
		newValidateCmd(),
		newSkillCmd(),
	)

	return cmd
}

// requireRepo discovers the nearest OpenPlan repository from the current
// working directory
func requireRepo() (repo.Context, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return repo.Context{}, err
	}
	return repo.Discover(cwd)
}

// formatListText keeps the human-readable list output deterministic and compact
// while still surfacing invalid files separately
func formatListText(result store.ListResult) (string, error) {
	var b strings.Builder

	if len(result.Documents) > 0 {
		tw := tabwriter.NewWriter(&b, 0, 4, 2, ' ', 0)
		_, _ = fmt.Fprintln(tw, "ID\tSTATUS\tTITLE\tTAGS\tPARENT")
		for _, doc := range result.Documents {
			parent := doc.Frontmatter.Parent
			if parent == "" {
				parent = "-"
			}
			_, _ = fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n",
				doc.Frontmatter.ID,
				doc.Frontmatter.Status,
				doc.Frontmatter.Title,
				strings.Join(doc.Frontmatter.Tags, ","),
				parent,
			)
		}
		if err := tw.Flush(); err != nil {
			return "", err
		}
	} else {
		b.WriteString("No plans found.\n")
	}

	if len(result.Issues) > 0 {
		if len(result.Documents) > 0 {
			b.WriteString("\n")
		}
		b.WriteString("Issues:\n")
		for _, issue := range result.Issues {
			fmt.Fprintf(&b, "- %s: %s: %s\n", issue.Path, issue.Field, issue.Message)
		}
	}

	return b.String(), nil
}

// formatShowText renders the metadata view without leaking full plan bodies
// into terminal output
func formatShowText(doc plan.Document) (string, error) {
	var b strings.Builder
	fmt.Fprintf(&b, "ID: %s\n", doc.Frontmatter.ID)
	fmt.Fprintf(&b, "Title: %s\n", doc.Frontmatter.Title)
	fmt.Fprintf(&b, "Status: %s\n", doc.Frontmatter.Status)
	fmt.Fprintf(&b, "Tags: %s\n", strings.Join(doc.Frontmatter.Tags, ","))
	if doc.Frontmatter.Parent == "" {
		b.WriteString("Parent: -\n")
	} else {
		fmt.Fprintf(&b, "Parent: %s\n", doc.Frontmatter.Parent)
	}
	fmt.Fprintf(&b, "Path: %s\n", doc.Path)

	return b.String(), nil
}

// openEditor delegates to the user-selected shell command so editor arguments
// and wrappers continue to work as configured in `$EDITOR`
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

// writeStdout centralizes command output so tests can replace the destination
func writeStdout(s string) error {
	_, err := io.WriteString(stdout, s)
	return err
}

// writeJSON keeps machine output consistently indented and newline-terminated
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
// before mutating repository state
func editorCommand() (string, error) {
	editor := strings.TrimSpace(os.Getenv("EDITOR"))
	if editor == "" {
		return "", EditorUnsetError{}
	}
	return editor, nil
}

// formatIssues turns validation issues into a single user-facing error string
// so CLI presentation stays outside the lower-level packages
func formatIssues(path string, issues []plan.ValidationIssue) error {
	var b strings.Builder
	fmt.Fprintf(&b, "invalid plan: %s\n", path)
	for _, issue := range issues {
		fmt.Fprintf(&b, "- %s: %s\n", issue.Field, issue.Message)
	}
	return errors.New(strings.TrimSpace(b.String()))
}

// formatValidationIssues renders repository-wide validation failures as one
// deterministic error so callers get a single non-zero result with full detail
func formatValidationIssues(issues []plan.ValidationIssue) error {
	var b strings.Builder
	b.WriteString("validation failed\n")
	for _, issue := range issues {
		fmt.Fprintf(&b, "- %s: %s: %s\n", issue.Path, issue.Field, issue.Message)
	}
	return errors.New(strings.TrimSpace(b.String()))
}

// listRow is the stable JSON shape for each `list` result row
type listRow struct {
	ID     string      `json:"id"`
	Title  string      `json:"title"`
	Status plan.Status `json:"status"`
	Tags   []string    `json:"tags"`
	Parent string      `json:"parent,omitempty"`
	Path   string      `json:"path"`
}

// listOutput is the machine-readable contract for `openplan list --json`
type listOutput struct {
	Plans  []listRow              `json:"plans"`
	Issues []plan.ValidationIssue `json:"issues,omitempty"`
}

// showOutput is the machine-readable contract for `openplan show --json`
type showOutput struct {
	ID     string      `json:"id"`
	Title  string      `json:"title"`
	Status plan.Status `json:"status"`
	Tags   []string    `json:"tags"`
	Parent string      `json:"parent,omitempty"`
	Path   string      `json:"path"`
}
