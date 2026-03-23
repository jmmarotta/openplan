package cli

import (
	"openplan/internal/store"

	"github.com/spf13/cobra"
)

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
