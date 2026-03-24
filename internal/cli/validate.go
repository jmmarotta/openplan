package cli

import (
	"openplan/internal/store"

	"github.com/spf13/cobra"
)

// newValidateCmd checks every plan file in the repository and fails if any
// frontmatter or validation issues are present
func newValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "validate",
		Aliases: []string{"v"},
		Short:   "Validate all plans in the repository",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := requireRepo()
			if err != nil {
				return err
			}

			result, err := store.List(ctx, store.Query{IncludeClosed: true})
			if err != nil {
				return err
			}
			if len(result.Issues) == 0 {
				return writeStdout("All plans valid.\n")
			}

			return formatValidationIssues(result.Issues)
		},
	}
}
