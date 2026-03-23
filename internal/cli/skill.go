package cli

import (
	_ "embed"

	"github.com/spf13/cobra"
)

//go:embed skill.md
var skillDoc string

// newSkillCmd prints the embedded AI-facing usage guide verbatim so it can be
// consumed by tools without runtime file lookups
func newSkillCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "skill",
		Short: "Print OpenPlan guidance for AI agents",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return writeStdout(skillDoc)
		},
	}
}
