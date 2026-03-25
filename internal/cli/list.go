package cli

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/jmmarotta/openplan/internal/plan"
	"github.com/jmmarotta/openplan/internal/store"

	"github.com/spf13/cobra"
)

// listRow is the stable JSON shape for each `list` result row.
type listRow struct {
	ID     string      `json:"id"`
	Title  string      `json:"title"`
	Status plan.Status `json:"status"`
	Tags   []string    `json:"tags"`
	Parent string      `json:"parent,omitempty"`
	Path   string      `json:"path"`
}

// listOutput is the machine-readable contract for `openplan list --json`.
type listOutput struct {
	Plans  []listRow              `json:"plans"`
	Issues []plan.ValidationIssue `json:"issues,omitempty"`
}

// newListCmd exposes repository browsing with optional filters while keeping
// invalid files visible as diagnostics instead of silently dropping them
func newListCmd() *cobra.Command {
	var includeClosed bool
	var statuses []string
	var tag string
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List plans in the repository",
		Args:    cobra.NoArgs,
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

// formatListText keeps the human-readable list output deterministic and compact
// while still surfacing invalid files separately.
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
