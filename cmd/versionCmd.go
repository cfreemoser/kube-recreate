package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

type VersionCmd struct {
	out     io.Writer
	version string
	commit  string
	branch  string
}

// NewVersionCommand prints the version of this kubectl plugin.
func NewVersionCommand(versionCmd *VersionCmd) *cobra.Command {

	cmd := &cobra.Command{
		Use:          "version",
		Short:        "Prints the version of kubectl-recreate",
		SilenceUsage: true,
		Run:          versionCmd.run,
	}

	return cmd
}

func (vc *VersionCmd) run(c *cobra.Command, args []string) {
	if len(vc.version) == 0 {
		fmt.Fprintf(vc.out, "[DEV BUILD] %s: %s\n", vc.branch, vc.commit)
		return
	}

	fmt.Fprintf(vc.out, "kubectl-recreate version: %s\n", vc.version)
	fmt.Fprintf(vc.out, "build commit: %s\n", vc.commit)
}
