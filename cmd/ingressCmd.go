package cmd

import (
	"errors"
	"io"
	"kube-recreate/pkg/k8s"
	"kube-recreate/pkg/util"

	kv1 "k8s.io/api/core/v1"
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
				return rCmd.run("", false)
			}
			if len(args) == 1 {
				return rCmd.run(args[0], false)
			}
			if getAllNamespacesFlag(c) {
				return rCmd.run("", true)

			}

			return errors.New("Define Resource or use --all")
		},
	}

	return cmd
}

func (ir *ingressCmd) run(name string, allNamespaces bool) error {
	client, err := k8s.NewK8sClient()
	if err != nil {
		return err
	}

	var ingresses []v1beta1.Ingress
	var namespaces []string
	if len(name) == 0 {
		if allNamespaces {
			nsList, err := client.LsNamespaces()
			if err != nil {
				return err
			}
			namespaces = append(namespaces, mapNamespacesToNames(nsList)...)
		} else {
			namespaces = append(namespaces, ir.ns)
		}

		for _, ns := range namespaces {
			l, err := client.LsIngress(ns)
			if err != nil {
				return err
			}

			ingresses = append(ingresses, l...)
		}
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

func mapNamespacesToNames(namespaces []kv1.Namespace) []string {
	var result []string
	for _, ns := range namespaces {
		result = append(result, ns.Name)
	}
	return result
}
