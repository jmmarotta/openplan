package cli

import (
	"fmt"

	"github.com/jmmarotta/openplan/internal/store"

	"github.com/spf13/cobra"
)

// newOutput is the machine-readable contract for `openplan new --json`.
type newOutput struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

// newNewCmd creates a templated plan file and hands control to the user's
// editor immediately so the CLI does not own plan-body mutation
func newNewCmd() *cobra.Command {
	var tags []string
	var parent string
	var noEdit bool
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:     "new [title]",
		Aliases: []string{"n"},
		Short:   "Create a new plan",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := requireRepo()
			if err != nil {
				return err
			}
			if !noEdit {
				if _, err := editorCommand(); err != nil {
					return err
				}
			}

			title := ""
			if len(args) == 1 {
				title = args[0]
			}

			doc, err := store.Create(ctx, store.NewPlanInput{
				Title:  title,
				Tags:   tags,
				Parent: parent,
			})
			if err != nil {
				return err
			}

			if jsonOutput {
				if err := writeJSON(newOutput{ID: doc.Frontmatter.ID, Path: doc.Path}); err != nil {
					return err
				}
			} else {
				if err := writeStdout(fmt.Sprintf("Created %s at %s\n", doc.Frontmatter.ID, doc.Path)); err != nil {
					return err
				}
			}

			if noEdit {
				return nil
			}

			return openEditor(doc.Path)
		},
	}

	cmd.Flags().StringSliceVar(&tags, "tag", nil, "Tag to attach to the new plan")
	cmd.Flags().StringVar(&parent, "parent", "", "Parent plan ID")
	cmd.Flags().BoolVar(&noEdit, "no-edit", false, "Create the plan without opening $EDITOR")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Render machine-readable JSON")
	return cmd
}
