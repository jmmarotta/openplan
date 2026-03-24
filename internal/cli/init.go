package cli

import (
	"fmt"
	"os"

	"openplan/internal/repo"

	"github.com/spf13/cobra"
)

// newInitCmd creates the repository-local `.plans/` directory and config file
func newInitCmd() *cobra.Command {
	var prefix string

	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Initialize OpenPlan in the current directory",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			ctx, err := repo.Init(cwd, repo.DefaultConfig(prefix))
			if err != nil {
				return err
			}

			return writeStdout(fmt.Sprintf("Initialized OpenPlan in %s with prefix %s\n", ctx.Root, ctx.Config.Prefix))
		},
	}

	cmd.Flags().StringVar(&prefix, "prefix", "OPN", "Ticket prefix for new plans")
	return cmd
}
