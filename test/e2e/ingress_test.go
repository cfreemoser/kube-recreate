package e2e

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	v1beta1 "k8s.io/api/networking/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRecreateAllNamespaces(t *testing.T) {
	// Arrange
	iCmd := RecreateCommand("ingress",
		"--all-namespaces",
	)
	ingressesBefore := mustGetAllIngressesInCluster(t)

	// Act
	mustExecute(t, iCmd)

	// Assert
	for _, ingressesAfter := range ingressesBefore {
		afterIngress, err := clientset.NetworkingV1beta1().Ingresses(ingressesAfter.Namespace).Get(context.Background(), ingressesAfter.Name, v1.GetOptions{})
		assert.NoError(t, err)
		assert.NotEqual(t, ingressesAfter.ResourceVersion, afterIngress.ResourceVersion)
	}
}

func TestRecreateWithObjectName(t *testing.T) {
	// Arrange
	iCmd := RecreateCommand(
		"ingress",
		"test-ingress-0",
		"-n",
		"test-0",
	)
	ingressesBefore := mustLsIngress(t, "test-0")

	// Act
	mustExecute(t, iCmd)

	// Assert
	ingressesAfter := mustLsIngress(t, "test-0")
	assert.NotEqual(t, ingressesBefore[0].ResourceVersion, ingressesAfter[0].ResourceVersion)
	for i := 1; i < 10; i++ {
		assert.Equal(t, ingressesBefore[i].ResourceVersion, ingressesAfter[i].ResourceVersion)
	}
}

func TestRecreateAllInNamespaces(t *testing.T) {
	iCmd := RecreateCommand("ingress",
		"-n",
		"test-1",
		"-a")

	ingressesBefore := mustLsIngress(t, "test-1")
	mustExecute(t, iCmd)

	ingressesAfter := mustLsIngress(t, "test-1")

	for i := 0; i < 10; i++ {
		assert.NotEqual(t, ingressesBefore[i].ResourceVersion, ingressesAfter[i].ResourceVersion)
	}
}

func mustLsIngress(t *testing.T, ns string) []v1beta1.Ingress {
	result, err := client.Ingress.Ls(ns)
	assert.NoError(t, err)
	return result
}

func mustGetAllIngressesInCluster(t *testing.T) []v1beta1.Ingress {
	var stateBefore []v1beta1.Ingress
	namespaces, err := client.Namespace.Ls()
	assert.NoError(t, err)
	for _, ns := range namespaces {
		stateBefore = append(stateBefore, mustLsIngress(t, ns.Name)...)
	}
	return stateBefore
}
