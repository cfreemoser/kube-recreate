package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type VersionCmd struct {
	settings *CmdSetting
}

// NewVersionCommand prints the version of this kubectl plugin.
func NewVersionCommand(settings *CmdSetting) *cobra.Command {
	versionCmd := &VersionCmd{settings: settings}

	cmd := &cobra.Command{
		Use:          "version",
		Short:        "Prints the version of kubectl-recreate",
		SilenceUsage: true,
		Run:          versionCmd.run,
	}

	return cmd
}

func (vc *VersionCmd) run(c *cobra.Command, args []string) {
	if len(vc.settings.CodeProperties.version) == 0 {
		fmt.Fprintf(vc.settings.Out, "[DEV BUILD] %s: %s\n", vc.settings.CodeProperties.branch, vc.settings.CodeProperties.commit)
		return
	}

	fmt.Fprintf(vc.settings.Out, "kubectl-recreate version: %s\n", vc.settings.CodeProperties.version)
	fmt.Fprintf(vc.settings.Out, "build commit: %s\n", vc.settings.CodeProperties.commit)
}
