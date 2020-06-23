package e2e

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/kind/pkg/cluster"
)

func TestMain(m *testing.M) {
	provider := cluster.NewProvider()
	// TODO bootstrepping und cleanup not working?
	setup(provider)
	code := m.Run()
	shutdown(provider)
	os.Exit(code)
}

func setup(provider *cluster.Provider) {
	provider.Create(
		"e2eTesting",
		cluster.CreateWithNodeImage("kindest/node:v1.18.0"),
		cluster.CreateWithWaitForReady(5*time.Minute),
	)
}

func shutdown(provider *cluster.Provider) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Could not perform cleanup")
	}

	kubeconfigPath := path.Join(homeDir, ".kube", "config")

	provider.Delete("e2eTesting", kubeconfigPath)
}

func createClientset() (*kubernetes.Clientset, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	kubeconfigPath := path.Join(homeDir, ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func mustClientset(t *testing.T) *kubernetes.Clientset {
	clientset, err := createClientset()
	assert.NoError(t, err)
	return clientset
}

func mockStreams() genericclioptions.IOStreams {
	return genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}

}
