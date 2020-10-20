package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewRefreshCommand creates the command for rendering the Kubernetes server version.
func NewRecreateCommand(streams genericclioptions.IOStreams, version, commit, branch string) *cobra.Command {

	cmd := &cobra.Command{
		Use:          "recreate",
		Short:        "Deletes and recreates all ingress resources",
		SilenceUsage: true,
	}

	cmd.PersistentFlags().StringP("namespace", "n", "", "Select the namespace")
	cmd.PersistentFlags().BoolP("all", "a", false, "If present, recreate all object(s) in namespace")
	cmd.PersistentFlags().Bool("all-namespaces", false, "If present, recreate the requested object(s) across all namespaces.")

	settings := NewCmdSettings(
		WithOutWriter(streams.Out),
		WithParentCmd(cmd),
		WithCodeBranch(branch),
		WithCodeCommit(commit),
		WithCodeVersion(version))

	cmd.AddCommand(NewIngressCommand(settings))
	cmd.AddCommand(NewVersionCommand(settings))

	return cmd
}
