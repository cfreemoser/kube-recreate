package cmd

import (
	"errors"
	"io"
	"kube-recreate/pkg/k8s"
	"kube-recreate/pkg/util"

	v1beta1 "k8s.io/api/networking/v1beta1"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type ingressCmd struct {
	out      io.Writer
	ns       string
	reporter *util.Reporter
}

func NewIngressCommand(streams genericclioptions.IOStreams) *cobra.Command {
	rCmd := &ingressCmd{
		out:      streams.Out,
		reporter: util.NewReporter(streams.Out),
	}

	cmd := &cobra.Command{
		Use:          "ingress",
		Short:        "Deletes and recreates all ingress resources",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			rCmd.ns = getNamespace(genericclioptions.NewConfigFlags(true), c)

			if getAllFlag(c) == true {
				return rCmd.run("")
			}
			if len(args) == 1 {
				return rCmd.run(args[0])
			}

			return errors.New("Define Resource or use --all")
		},
	}

	return cmd
}

func (ir *ingressCmd) run(name string) error {
	client, err := k8s.NewK8sClient()
	if err != nil {
		return err
	}

	ingL := make([]v1beta1.Ingress, 0)

	if len(name) == 0 {
		l, err := client.LsIngress(ir.ns)
		if err != nil {
			return err
		}

		ingL = append(ingL, l...)
	} else {
		i, err := client.GetIngress(ir.ns, name)
		if err != nil {
			return err
		}

		ingL = append(ingL, i)
	}

	for _, ingress := range ingL {
		client.DeleteIngress(&ingress)
		ir.reporter.Append(ingress.Name, "Ingress", "DELETED", ingress.CreationTimestamp.String())

	}

	ir.reporter.AddSeperator()

	for _, ingress := range ingL {
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
