package e2e

import (
	"context"
	"fmt"
	"kube-recreate/cmd"
	"testing"

	"github.com/stretchr/testify/assert"
	v1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
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

func populateClusterWithIngresses(t *testing.T, clientset *kubernetes.Clientset) []*v1beta1.Ingress {
	var ingresses []*v1beta1.Ingress
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		ing, err := clientset.NetworkingV1beta1().Ingresses("default").Create(ctx, createIngress(fmt.Sprintf("test-ingress-%d", i)), v1.CreateOptions{})
		assert.NoError(t, err)
		ingresses = append(ingresses, ing)
	}

	return ingresses
}

func TestDeletionOfOneIngress(t *testing.T) {
	clientset := mustClientset(t)
	iCmd := cmd.NewIngressCommand(mockStreams())
	iCmd.SetArgs([]string{
		"test-ingress-0",
	})

	beforeList := populateClusterWithIngresses(t, clientset)

	err := iCmd.Execute()
	assert.NoError(t, err)

	afterList, _ := clientset.NetworkingV1beta1().Ingresses("default").List(context.Background(), v1.ListOptions{})

	assert.NotEqual(t, beforeList[0].ResourceVersion, afterList.Items[0].ResourceVersion)
	for i := 1; i < 10; i++ {
		assert.Equal(t, beforeList[i].ResourceVersion, afterList.Items[i].ResourceVersion)
	}
}
