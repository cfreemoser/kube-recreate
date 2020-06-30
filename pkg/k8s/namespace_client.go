package k8s

import (
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Namespace struct {
	baseClient K8sClient
}

func (ns *Namespace) Ls() ([]v1.Namespace, error) {
	list, err := ns.baseClient.Clientset.CoreV1().Namespaces().List(ns.baseClient.Context, metaV1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}
