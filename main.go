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
	refreshCmd := cmd.NewRecreateCommand(streams, VERSION, COMMIT, BRANCH)
	if err := refreshCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
