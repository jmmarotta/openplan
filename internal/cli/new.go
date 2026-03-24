package cli

import (
	"fmt"

	"openplan/internal/store"

	"github.com/spf13/cobra"
)

// newNewCmd creates a templated plan file and hands control to the user's
// editor immediately so the CLI does not own plan-body mutation
func newNewCmd() *cobra.Command {
	var tags []string
	var parent string

	cmd := &cobra.Command{
		Use:     "new [title]",
		Aliases: []string{"n"},
		Short:   "Create a new plan and open it in $EDITOR",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := requireRepo()
			if err != nil {
				return err
			}
			if _, err := editorCommand(); err != nil {
				return err
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

			if err := writeStdout(fmt.Sprintf("Created %s at %s\n", doc.Frontmatter.ID, doc.Path)); err != nil {
				return err
			}
			return openEditor(doc.Path)
		},
	}

	cmd.Flags().StringSliceVar(&tags, "tag", nil, "Tag to attach to the new plan")
	cmd.Flags().StringVar(&parent, "parent", "", "Parent plan ID")
	return cmd
}
