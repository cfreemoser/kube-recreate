package k8s

import (
	"context"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClient struct {
	Clientset *kubernetes.Clientset
	Context   context.Context
	Ingress   *Ingress
	Namespace *Namespace
}

func NewK8sClient() (*K8sClient, error) {
	client := &K8sClient{
		Context: context.Background(),
	}

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

	client.Clientset = clientset
	client.Ingress = &Ingress{baseClient: *client}
	client.Namespace = &Namespace{baseClient: *client}

	return client, nil
}
