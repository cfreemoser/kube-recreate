package k8s

import (
	v1beta1 "k8s.io/api/networking/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Ingress struct {
	baseClient K8sClient
}

func (client *Ingress) Ls(namespace string) ([]v1beta1.Ingress, error) {
	iclient := client.baseClient.Clientset.NetworkingV1beta1().Ingresses(namespace)

	ingressesList, err := iclient.List(client.baseClient.Context, v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return ingressesList.Items, nil
}

func (client *Ingress) Get(name, namespace string) (v1beta1.Ingress, error) {
	iclient := client.baseClient.Clientset.NetworkingV1beta1().Ingresses(namespace)

	ingress, err := iclient.Get(client.baseClient.Context, name, v1.GetOptions{})

	if err != nil {
		return v1beta1.Ingress{}, err
	}

	return *ingress, nil
}

func (client *Ingress) Delete(ingress *v1beta1.Ingress) (*v1beta1.Ingress, error) {
	iclient := client.baseClient.Clientset.NetworkingV1beta1().Ingresses(ingress.Namespace)
	return ingress, iclient.Delete(client.baseClient.Context, ingress.Name, v1.DeleteOptions{})
}

func (client *Ingress) Create(ingress *v1beta1.Ingress) (*v1beta1.Ingress, error) {
	ingress.ResourceVersion = ""
	iclient := client.baseClient.Clientset.NetworkingV1beta1().Ingresses(ingress.Namespace)
	return iclient.Create(client.baseClient.Context, ingress, v1.CreateOptions{})
}
