package cmd

import (
	"kube-recreate/pkg/util"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewRefreshCommand creates the command for rendering the Kubernetes server version.
func NewRefreshCommand(streams genericclioptions.IOStreams, version, commit, branch string) *cobra.Command {

	cmd := &cobra.Command{
		Use:          "recreate",
		Short:        "Deletes and recreates all ingress resources",
		SilenceUsage: true,
	}

	cmd.PersistentFlags().StringP("namespace", "n", "", "Set the namespace")
	cmd.PersistentFlags().BoolP("all", "a", false, "All resources in namespace")

	vCmd := &VersionCmd{
		out:     streams.Out,
		commit:  commit,
		version: version,
		branch:  branch,
	}

	rCmd := &ingressCmd{
		out:      streams.Out,
		reporter: util.NewReporter(streams.Out),
	}

	cmd.AddCommand(NewIngressCommand(rCmd))
	cmd.AddCommand(NewVersionCommand(vCmd))

	return cmd
}

// getNamespace takes a set of kubectl flag values and returns the namespace we should be operating in
func getNamespace(flags *genericclioptions.ConfigFlags, cmd *cobra.Command) string {
	ns, err := cmd.Flags().GetString("namespace")
	if err != nil || len(ns) == 0 {
		namespace, _, err := flags.ToRawKubeConfigLoader().Namespace()
		if err != nil || len(namespace) == 0 {
			namespace = "default"
		}
		return namespace
	}

	return ns
}

func getAllFlag(cmd *cobra.Command) bool {
	all, err := cmd.Flags().GetBool("all")
	if err != nil {
		return false
	}
	return all
}
