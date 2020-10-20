package e2e

import (
	"kube-recreate/pkg/k8s"
	"os"
	"path"
	"testing"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/kind/pkg/cluster"
)

var (
	clientset *kubernetes.Clientset
	client    *k8s.K8sClient
)

func TestMain(m *testing.M) {
	provider := cluster.NewProvider()
	setup(provider)
	code := m.Run()
	shutdown(provider)
	os.Exit(code)
}

func createCluster(provider *cluster.Provider) {
	err := provider.Create(
		"e2e",
		cluster.CreateWithNodeImage("kindest/node:v1.18.0"),
		cluster.CreateWithWaitForReady(5*time.Minute),
		cluster.CreateWithDisplayUsage(true),
		cluster.CreateWithDisplaySalutation(true),
	)
	if err != nil {
		panic(err)
	}
}

func setup(provider *cluster.Provider) {
	createCluster(provider)
	clientset = mustClientset()
	client = mustClient()
	populateClusterWithIngresses(clientset, createNamespace(clientset))
}

func shutdown(provider *cluster.Provider) {
	kubeconfigPath := mustKubeconfigPath()

	err := provider.Delete("e2eTesting", kubeconfigPath)
	if err != nil {
		panic(err)
	}
}

func mustKubeconfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return path.Join(homeDir, ".kube", "config")
}

func createClientset() (*kubernetes.Clientset, error) {
	kubeconfigPath := mustKubeconfigPath()

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func mustClientset() *kubernetes.Clientset {
	clientset, err := createClientset()
	if err != nil {
		panic(err)
	}
	return clientset
}

func mustClient() *k8s.K8sClient {
	client, err := k8s.NewK8sClient()
	if err != nil {
		panic(err)
	}
	return client
}
