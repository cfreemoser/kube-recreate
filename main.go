package main

import (
	"os"

	"refresher/cmd"

	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func main() {
	// cmd.SetVersion(version)
	flags := pflag.NewFlagSet("kubectl-ns", pflag.ExitOnError)
	pflag.CommandLine = flags

	refreshCmd := cmd.NewRefreshCommand(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := refreshCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
