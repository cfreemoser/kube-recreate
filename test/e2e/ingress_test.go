package e2e

import (
	"context"
	"fmt"
	"kube-recreate/cmd"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	v1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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
		TypeMeta: v1.TypeMeta{},
	}
}

func lsIngress(t *testing.T, ns string) []v1beta1.Ingress {
	temp, err := clientset.NetworkingV1beta1().Ingresses(ns).List(context.Background(), v1.ListOptions{})
	assert.NoError(t, err)
	return temp.Items
}

func RefreshCommand() *cobra.Command {
	return cmd.NewRefreshCommand(mockStreams(), "", "", "")
}

func TestDeletionOfAllIngressesInAllNamespaces(t *testing.T) {
	iCmd := RefreshCommand()
	iCmd.SetArgs([]string{
		"ingress",
		"--all-namespaces",
	})

	var beforeList []v1beta1.Ingress

	for i := 0; i < 10; i++ {
		beforeList = append(beforeList, lsIngress(t, fmt.Sprintf("test-%d", i))...)
	}

	err := iCmd.Execute()
	assert.NoError(t, err)

	for _, beforeIngress := range beforeList {
		afterIngress, err := clientset.NetworkingV1beta1().Ingresses(beforeIngress.Namespace).Get(context.Background(), beforeIngress.Name, v1.GetOptions{})
		assert.NoError(t, err)
		assert.NotEqual(t, beforeIngress.ResourceVersion, afterIngress.ResourceVersion)
	}
}

func TestDeletionOfOneIngress(t *testing.T) {
	iCmd := RefreshCommand()
	iCmd.SetArgs([]string{
		"ingress",
		"test-ingress-0",
		"-n",
		"test-0",
	})

	beforeList := lsIngress(t, "test-0")

	err := iCmd.Execute()
	assert.NoError(t, err)

	afterList := lsIngress(t, "test-0")
	assert.NoError(t, err)

	assert.NotEqual(t, beforeList[0].ResourceVersion, afterList[0].ResourceVersion)
	for i := 1; i < 10; i++ {
		assert.Equal(t, beforeList[i].ResourceVersion, afterList[i].ResourceVersion)
	}
}

func TestDeletionOfAllIngress(t *testing.T) {
	iCmd := RefreshCommand()
	iCmd.SetArgs([]string{
		"ingress",
		"-n",
		"test-1",
		"-a",
	})

	beforeList := lsIngress(t, "test-1")
	err := iCmd.Execute()
	assert.NoError(t, err)

	afterList := lsIngress(t, "test-1")

	for i := 0; i < 10; i++ {
		assert.NotEqual(t, beforeList[i].ResourceVersion, afterList[i].ResourceVersion)
	}
}
