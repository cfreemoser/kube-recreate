package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type VersionCmd struct {
	out     io.Writer
	version string
	commit  string
	branch  string
}

// NewVersionCommand prints the version of this kubectl plugin.
func NewVersionCommand(streams genericclioptions.IOStreams, version, commit, branch string) *cobra.Command {
	vCmd := &VersionCmd{
		out:     streams.Out,
		commit:  commit,
		version: version,
		branch:  branch,
	}

	cmd := &cobra.Command{
		Use:          "version",
		Short:        "Prints the version of kubectl-recreate",
		SilenceUsage: true,
		Run:          vCmd.run,
	}

	return cmd
}

func (vc *VersionCmd) run(c *cobra.Command, args []string) {
	if len(vc.version) == 0 {
		fmt.Fprintf(vc.out, "[DEV BUILD] %s: %s\n", vc.branch, vc.commit)
		return
	}

	fmt.Fprintf(vc.out, "kubectl-recreate version: %s\n", vc.version)
}
