package cmd

import (
	"kube-recreate/pkg/k8s"
	"kube-recreate/pkg/util"

	kv1 "k8s.io/api/core/v1"
	v1beta1 "k8s.io/api/networking/v1beta1"

	"github.com/spf13/cobra"
)

type IngressCmd struct {
	settings *CmdSetting
}

func NewIngressCommand(settings *CmdSetting) *cobra.Command {

	iCmd := &IngressCmd{settings: settings}

	cmd := &cobra.Command{
		Use:          "ingress [name]",
		Short:        "Deletes and recreates all ingress resources",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) == 1 {
				settings.ObjectNameProvided = true
				settings.ObjectName = args[0]
			}

			return iCmd.run()
		},
	}

	return cmd
}

func (ir *IngressCmd) run() error {
	client, err := k8s.NewK8sClient()
	if err != nil {
		return err
	}

	var ingresses []v1beta1.Ingress
	var namespaces []string
	if ir.settings.ObjectNameProvided == false {
		if ir.settings.AllNamespacesFlag() {
			nsList, err := client.LsNamespaces()
			if err != nil {
				return err
			}
			namespaces = append(namespaces, mapNamespacesToNames(nsList)...)
		} else {
			namespaces = append(namespaces, ir.settings.Namespace())
		}

		for _, ns := range namespaces {
			l, err := client.LsIngress(ns)
			if err != nil {
				return err
			}

			ingresses = append(ingresses, l...)
		}
	} else {
		i, err := client.GetIngress(ir.settings.Namespace(), ir.settings.ObjectName)
		if err != nil {
			return err
		}

		ingresses = append(ingresses, i)
	}

	deleteAndRecreate(ingresses, ir.settings.Reporter, client)

	ir.settings.Reporter.PrintReport()
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
