package kubernetes

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c KubernetesClient) GetNamespaces() ([]string, error) {
	namespaces, err := c.clientset.CoreV1().Namespaces().List(context.Background(), v1.ListOptions{})
	if err != nil {
		return []string{}, err
	}
	res := make([]string, 0, len(namespaces.Items))
	for _, ns := range namespaces.Items {
		res = append(res, ns.Name)
	}
	return res, nil
}
