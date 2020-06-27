package main

import (
	"os"

	"kube-recreate/cmd"

	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var COMMIT string
var BRANCH string
var VERSION string

func main() {
	flags := pflag.NewFlagSet("kubectl-ns", pflag.ExitOnError)
	pflag.CommandLine = flags

	streams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	refreshCmd := cmd.NewRefreshCommand(streams)
	refreshCmd.PersistentFlags().StringP("namespace", "n", "", "Set the namespace")
	refreshCmd.PersistentFlags().BoolP("all", "a", false, "All resources in namespace")

	refreshCmd.AddCommand(cmd.NewIngressCommand(streams))
	refreshCmd.AddCommand(cmd.NewVersionCommand(streams, VERSION, COMMIT, BRANCH))
	if err := refreshCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
