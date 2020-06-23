package cmd

import (
	"errors"
	"io"
	"refresher/pkg/k8s"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type refreshCmd struct {
	out          io.Writer
	ns           string
	resourceType ResourceType
	reporter     *Reporter
}

type ResourceType string

const (
	Ingress ResourceType = "Ingress"
)

// NewRefreshCommand creates the command for rendering the Kubernetes server version.
func NewRefreshCommand(streams genericclioptions.IOStreams) *cobra.Command {
	rCmd := &refreshCmd{
		out:      streams.Out,
		ns:       getNamespace(genericclioptions.NewConfigFlags(true)),
		reporter: NewReporter(streams.Out),
	}

	cmd := &cobra.Command{
		Use:          "refresh",
		Short:        "Deletes and recreates all ingress resources",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("this command does not accept arguments")
			}
			err := rCmd.parseArgs(args)
			if err != nil {
				return err
			}
			return rCmd.run()
		},
	}

	return cmd
}

func (ir *refreshCmd) run() error {
	client, err := k8s.NewK8sClient()
	if err != nil {
		return err
	}

	l, err := client.LsIngress(ir.ns)

	for _, ingress := range l {
		client.DeleteIngress(&ingress)
		ir.reporter.Append(ingress.Name, "Ingress", "DELETED", ingress.CreationTimestamp.String())

	}

	ir.reporter.AddSeperator()

	for _, ingress := range l {
		ingress.ResourceVersion = ""

		i, err := client.CreateIngress(&ingress)
		if err != nil {
			ir.reporter.Append(ingress.Name, "Ingress", "FAILED", ingress.CreationTimestamp.String())
		}

		ir.reporter.Append(i.Name, "Ingress", "CREATED", i.CreationTimestamp.String())
	}

	ir.reporter.PrintReport()
	return nil
}

func (ir *refreshCmd) parseArgs(args []string) error {
	rArg := args[0]
	switch rArg {
	case "ingress":
		ir.resourceType = ResourceType("Ingress")
		return nil
	default:
		return errors.New("Unkown resource type")
	}
}

// getNamespace takes a set of kubectl flag values and returns the namespace we should be operating in
func getNamespace(flags *genericclioptions.ConfigFlags) string {
	namespace, _, err := flags.ToRawKubeConfigLoader().Namespace()
	if err != nil || len(namespace) == 0 {
		namespace = "default"
	}
	return namespace
}
