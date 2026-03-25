package cli

import (
	"github.com/spf13/cobra"
)

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
