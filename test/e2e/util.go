package e2e

import (
	"context"
	"fmt"
	"kube-recreate/cmd"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	v1beta1 "k8s.io/api/networking/v1beta1"
	apiV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"

	"k8s.io/apimachinery/pkg/util/intstr"
)

func createIngress(name string) *v1beta1.Ingress {
	return &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Status: v1beta1.IngressStatus{},
		Spec: v1beta1.IngressSpec{
			Backend: &v1beta1.IngressBackend{
				ServiceName: "test",
				ServicePort: intstr.FromInt(9999),
			},
		},
		TypeMeta: metav1.TypeMeta{},
	}
}

func createNamespace(clientset *kubernetes.Clientset) []*v1.Namespace {
	var namespaces []*v1.Namespace

	ctx := context.Background()
	for i := 0; i < 10; i++ {
		ns, err := clientset.CoreV1().Namespaces().Create(ctx, &v1.Namespace{ObjectMeta: apiV1.ObjectMeta{Name: fmt.Sprintf("test-%d", i)}}, apiV1.CreateOptions{})
		if err != nil {
			panic(err)
		}
		namespaces = append(namespaces, ns)
	}
	return namespaces
}

func populateClusterWithIngresses(clientset *kubernetes.Clientset, namespaces []*v1.Namespace) []*v1beta1.Ingress {
	var ingresses []*v1beta1.Ingress
	ctx := context.Background()
	for _, ns := range namespaces {
		for i := 0; i < 10; i++ {
			ing, err := clientset.NetworkingV1beta1().Ingresses(ns.Name).Create(ctx, createIngress(fmt.Sprintf("test-ingress-%d", i)), apiV1.CreateOptions{})
			if err != nil {
				panic(err)
			}
			ingresses = append(ingresses, ing)
		}
	}

	return ingresses
}

func RecreateCommand(args ...string) *cobra.Command {
	cmd := cmd.NewRecreateCommand(mockStreams(), "", "", "")
	cmd.SetArgs(args)
	return cmd
}

func mustExecute(t *testing.T, iCmd *cobra.Command) {
	err := iCmd.Execute()
	assert.NoError(t, err)
}

func mockStreams() genericclioptions.IOStreams {
	return genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
}
