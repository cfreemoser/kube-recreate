package cmd

import (
	"kube-recreate/pkg/k8s"

	v1beta1 "k8s.io/api/networking/v1beta1"

	"github.com/spf13/cobra"
)

type IngressCmd struct {
	settings  *CmdSetting
	client    *k8s.K8sClient
	ingresses []v1beta1.Ingress
}

func NewIngressCommand(settings *CmdSetting) *cobra.Command {

	ingressCmd := &IngressCmd{settings: settings}

	return &cobra.Command{
		Use:          "ingress [name]",
		Short:        "Deletes and recreates all ingress resources",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) == 1 {
				settings.ObjectNameProvided = true
				settings.ObjectName = args[0]
			}

			err := ingressCmd.init()
			if err != nil {
				return err
			}

			return ingressCmd.run()
		},
	}
}

func (ingressCmd *IngressCmd) run() error {
	if ingressCmd.settings.ObjectNameProvided {
		objectName := ingressCmd.settings.ObjectName
		namespace := ingressCmd.settings.Namespace()
		i, err := ingressCmd.client.Ingress.Get(objectName, namespace)
		if err != nil {
			return err
		}

		ingressCmd.ingresses = append(ingressCmd.ingresses, i)
	}

	if ingressCmd.settings.AllFlag() {
		err := ingressCmd.appendIngressesFromNamespace(ingressCmd.settings.Namespace())
		if err != nil {
			return err
		}
	}

	if ingressCmd.settings.AllNamespacesFlag() {
		namespaces, err := ingressCmd.client.Namespace.Ls()
		if err != nil {
			return err
		}

		for _, ns := range namespaces {
			err := ingressCmd.appendIngressesFromNamespace(ns.Name)
			if err != nil {
				return err
			}
		}

	}

	ingressCmd.deleteAndRecreate()

	ingressCmd.settings.Reporter.PrintReport()
	return nil
}

func (ir *IngressCmd) init() error {
	client, err := k8s.NewK8sClient()
	if err != nil {
		return err
	}

	ir.client = client
	return nil
}

func (ir *IngressCmd) deleteAndRecreate() {
	ir.ExecuteClientFunctionAndReport(ir.client.Ingress.Delete, "DELETE")
	ir.settings.Reporter.AddSeperator()
	ir.ExecuteClientFunctionAndReport(ir.client.Ingress.Create, "CREATE")
}

func (ir *IngressCmd) ExecuteClientFunctionAndReport(clientFunc func(ingress *v1beta1.Ingress) (*v1beta1.Ingress, error), verb string) {
	for _, ingress := range ir.ingresses {
		_, err := clientFunc(&ingress)
		if err != nil {
			ir.settings.Reporter.Append(ingress.Name, "INGRESS", "FAILED", ingress.CreationTimestamp.Time)
		}
		ir.settings.Reporter.Append(ingress.Name, "INGRESS", verb, ingress.CreationTimestamp.Time)
	}
}

func (ir *IngressCmd) appendIngressesFromNamespace(namespace string) error {
	objects, err := ir.client.Ingress.Ls(namespace)
	if err != nil {
		return err
	}

	ir.ingresses = append(ir.ingresses, objects...)
	return nil
}
