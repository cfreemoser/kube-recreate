package e2e

import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	v1beta1 "k8s.io/api/networking/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/kind/pkg/cluster"
)

var (
	clientset *kubernetes.Clientset
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
	populateClusterWithIngresses(clientset)
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

func mockStreams() genericclioptions.IOStreams {
	return genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
}

func populateClusterWithIngresses(clientset *kubernetes.Clientset) []*v1beta1.Ingress {
	var ingresses []*v1beta1.Ingress
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		ing, err := clientset.NetworkingV1beta1().Ingresses("default").Create(ctx, createIngress(fmt.Sprintf("test-ingress-%d", i)), v1.CreateOptions{})
		if err != nil {
			panic(err)
		}
		ingresses = append(ingresses, ing)
	}

	return ingresses
}
