package reconcile

import (
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"os"
	"sigs.k8s.io/yaml"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ConfigMap(ctx context.Context, cs *kubernetes.Clientset, configMap *corev1.ConfigMap) error {
	if os.Getenv("DRYRUN") != "" {
		os.Stdout.WriteString("---\n")
		configMap.SetGroupVersionKind(schema.FromAPIVersionAndKind("v1", "ConfigMap"))
		b, err := yaml.Marshal(configMap)
		if err != nil {
			return err
		}
		_, err = os.Stdout.Write(b)
		return err
	}
	client := cs.CoreV1().ConfigMaps(configMap.Namespace)
	existing, err := client.Get(ctx, configMap.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			_, err = client.Create(ctx, configMap, metav1.CreateOptions{})
			return err
		}
		return err
	}
	existing.Labels = configMap.Labels
	existing.Annotations = configMap.Annotations
	existing.Data = configMap.Data
	existing.BinaryData = configMap.BinaryData
	existing.Immutable = configMap.Immutable
	_, err = client.Update(ctx, existing, metav1.UpdateOptions{})
	return err
}

func ConfigMapAbsence(ctx context.Context, cs *kubernetes.Clientset, configMap *corev1.ConfigMap) error {
	return Absence(func() error {
		return cs.CoreV1().ConfigMaps(configMap.Namespace).Delete(ctx, configMap.Name, metav1.DeleteOptions{})
	})
}
