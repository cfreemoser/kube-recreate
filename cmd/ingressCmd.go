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

func NewIngressCommand(rCmd *ingressCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "ingress [name]",
		Short:        "Deletes and recreates all ingress resources",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			rCmd.ns = getNamespace(genericclioptions.NewConfigFlags(true), c)

			if getAllFlag(c) {
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

	var ingresses []v1beta1.Ingress
	if len(name) == 0 {
		l, err := client.LsIngress(ir.ns)
		if err != nil {
			return err
		}

		ingresses = append(ingresses, l...)
	} else {
		i, err := client.GetIngress(ir.ns, name)
		if err != nil {
			return err
		}

		ingresses = append(ingresses, i)
	}

	deleteAndRecreate(ingresses, ir.reporter, client)

	ir.reporter.PrintReport()
	return nil
}

func deleteAndRecreate(ingresses []v1beta1.Ingress, reporter *util.Reporter, client *k8s.K8sClient) {
	for _, ingress := range ingresses {
		err := client.DeleteIngress(&ingress)
		if err != nil {
			reporter.Append(ingress.Name, "Ingress", "FAILED", ingress.CreationTimestamp.String())
		}
		reporter.Append(ingress.Name, "Ingress", "DELETED", ingress.CreationTimestamp.String())
	}

	reporter.AddSeperator()

	for _, ingress := range ingresses {
		ingress.ResourceVersion = ""

		i, err := client.CreateIngress(&ingress)
		if err != nil {
			reporter.Append(ingress.Name, "Ingress", "FAILED", ingress.CreationTimestamp.String())
		}

		reporter.Append(i.Name, "Ingress", "CREATED", i.CreationTimestamp.String())
	}
}
