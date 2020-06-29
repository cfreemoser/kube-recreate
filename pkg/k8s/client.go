package k8s

import (
	"context"
	"fmt"

	kv1 "k8s.io/api/core/v1"
	v1beta1 "k8s.io/api/networking/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClient struct {
	clientset *kubernetes.Clientset
}

func NewK8sClient() (*K8sClient, error) {
	client := &K8sClient{}

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return client, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return client, err
	}

	client.clientset = clientset

	return client, nil
}

func (client *K8sClient) LsIngress(namespace string) ([]v1beta1.Ingress, error) {
	ctx := context.Background()
	iclient := client.clientset.NetworkingV1beta1().Ingresses(namespace)

	ingressesList, err := iclient.List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return ingressesList.Items, nil
}

func (client *K8sClient) GetIngress(namespace, name string) (v1beta1.Ingress, error) {
	ctx := context.Background()
	iclient := client.clientset.NetworkingV1beta1().Ingresses(namespace)

	ingress, err := iclient.Get(ctx, name, v1.GetOptions{})
	test, _ := client.LsIngress(namespace)
	for _, t := range test {
		fmt.Println(t.Name)
	}
	if err != nil {
		return v1beta1.Ingress{}, err
	}

	return *ingress, nil
}

func (client *K8sClient) DeleteIngress(ingress *v1beta1.Ingress) error {
	ctx := context.Background()
	iclient := client.clientset.NetworkingV1beta1().Ingresses(ingress.Namespace)
	return iclient.Delete(ctx, ingress.Name, v1.DeleteOptions{})
}

func (client *K8sClient) CreateIngress(ingress *v1beta1.Ingress) (*v1beta1.Ingress, error) {
	ctx := context.Background()
	iclient := client.clientset.NetworkingV1beta1().Ingresses(ingress.Namespace)
	return iclient.Create(ctx, ingress, v1.CreateOptions{})
}

func (client *K8sClient) LsNamespaces() ([]kv1.Namespace, error) {
	ctx := context.Background()
	list, err := client.clientset.CoreV1().Namespaces().List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}
