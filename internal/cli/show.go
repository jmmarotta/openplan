package cli

import (
	"fmt"
	"strings"

	"github.com/jmmarotta/openplan/internal/plan"
	"github.com/jmmarotta/openplan/internal/store"

	"github.com/spf13/cobra"
)

// showOutput is the machine-readable contract for `openplan show --json`.
type showOutput struct {
	ID     string      `json:"id"`
	Title  string      `json:"title"`
	Status plan.Status `json:"status"`
	Tags   []string    `json:"tags"`
	Parent string      `json:"parent,omitempty"`
	Path   string      `json:"path"`
}

// newShowCmd resolves one plan to metadata and path, reporting validation
// issues as repairable output instead of dumping the raw markdown body
func newShowCmd() *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "show <ID>",
		Short: "Show plan metadata and path",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := requireRepo()
			if err != nil {
				return err
			}

			result, err := store.Get(ctx, args[0])
			if err != nil {
				return err
			}
			if len(result.ValidationIssues) > 0 {
				return formatIssues(result.Path, result.ValidationIssues)
			}
			if result.Document == nil {
				return formatIssues(result.Path, nil)
			}
			doc := *result.Document

			if jsonOutput {
				return writeJSON(showOutput{
					ID:     doc.Frontmatter.ID,
					Title:  doc.Frontmatter.Title,
					Status: doc.Frontmatter.Status,
					Tags:   doc.Frontmatter.Tags,
					Parent: doc.Frontmatter.Parent,
					Path:   doc.Path,
				})
			}

			text, err := formatShowText(doc)
			if err != nil {
				return err
			}
			return writeStdout(text)
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Render machine-readable JSON")
	return cmd
}

// formatShowText renders the metadata view without leaking full plan bodies
// into terminal output.
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
