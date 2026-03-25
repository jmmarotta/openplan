package cli

import (
	"fmt"
	"os"

	"github.com/jmmarotta/openplan/internal/plan"

	"github.com/spf13/cobra"
)

// newEditCmd reopens a plan by full ID without adding any mutation behavior to
// the CLI itself
func newEditCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "edit <ID>",
		Aliases: []string{"e"},
		Short:   "Open an existing plan in $EDITOR",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := requireRepo()
			if err != nil {
				return err
			}
			if _, err := editorCommand(); err != nil {
				return err
			}
			if _, err := plan.ParseID(args[0]); err != nil {
				return err
			}

			path := ctx.PlanPath(args[0])
			if _, err := os.Stat(path); err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("plan not found: %s", args[0])
				}
				return err
			}

			return openEditor(path)
		},
	}
}
