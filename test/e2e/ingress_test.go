package e2e

import (
	"context"
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

// func createNs(name string) *core.Namespace {
// 	return &core.Namespace{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name: name,
// 		},
// 		Spec:     core.NamespaceSpec{},
// 		Status:   core.NamespaceStatus{},
// 		TypeMeta: v1.TypeMeta{},
// 	}
// }

func lsIngress(t *testing.T) []v1beta1.Ingress {
	temp, err := clientset.NetworkingV1beta1().Ingresses("default").List(context.Background(), v1.ListOptions{})
	assert.NoError(t, err)
	return temp.Items
}

func RefreshCommand() *cobra.Command {
	return cmd.NewRefreshCommand(mockStreams(), "", "", "")
}

func TestDeletionOfOneIngress(t *testing.T) {
	iCmd := RefreshCommand()
	iCmd.SetArgs([]string{
		"ingress",
		"test-ingress-0",
	})

	beforeList := lsIngress(t)

	err := iCmd.Execute()
	assert.NoError(t, err)

	afterList := lsIngress(t)
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
		"-a",
	})

	beforeList := lsIngress(t)

	err := iCmd.Execute()
	assert.NoError(t, err)

	afterList := lsIngress(t)

	for i := 0; i < 10; i++ {
		assert.NotEqual(t, beforeList[i].ResourceVersion, afterList[i].ResourceVersion)
	}
}
