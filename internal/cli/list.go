package cli

import (
	"fmt"
	"strings"

	"openplan/internal/plan"
	"openplan/internal/store"

	"github.com/spf13/cobra"
)

// newListCmd exposes repository browsing with optional filters while keeping
// invalid files visible as diagnostics instead of silently dropping them
func newListCmd() *cobra.Command {
	var includeClosed bool
	var statuses []string
	var tag string
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List plans in the repository",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := requireRepo()
			if err != nil {
				return err
			}

			parsedStatuses := make([]plan.Status, 0, len(statuses))
			for _, raw := range statuses {
				status := plan.Status(strings.ToLower(strings.TrimSpace(raw)))
				if !status.Valid() {
					return fmt.Errorf("invalid status %q", raw)
				}
				parsedStatuses = append(parsedStatuses, status)
			}

			result, err := store.List(ctx, store.Query{
				IncludeClosed: includeClosed,
				Statuses:      parsedStatuses,
				Tag:           strings.ToLower(strings.TrimSpace(tag)),
			})
			if err != nil {
				return err
			}

			if jsonOutput {
				rows := make([]listRow, 0, len(result.Documents))
				for _, doc := range result.Documents {
					rows = append(rows, listRow{
						ID:     doc.Frontmatter.ID,
						Title:  doc.Frontmatter.Title,
						Status: doc.Frontmatter.Status,
						Tags:   doc.Frontmatter.Tags,
						Parent: doc.Frontmatter.Parent,
						Path:   doc.Path,
					})
				}
				return writeJSON(listOutput{Plans: rows, Issues: result.Issues})
			}

			text, err := formatListText(result)
			if err != nil {
				return err
			}
			return writeStdout(text)
		},
	}

	cmd.Flags().BoolVar(&includeClosed, "all", false, "Include done and cancelled plans")
	cmd.Flags().StringSliceVar(&statuses, "status", nil, "Filter by status")
	cmd.Flags().StringVar(&tag, "tag", "", "Filter by tag")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Render machine-readable JSON")
	return cmd
}
