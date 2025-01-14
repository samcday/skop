package reconcile

import (
	"context"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"os"
	"sigs.k8s.io/yaml"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func DaemonSet(ctx context.Context, cs *kubernetes.Clientset, daemonSet *appsv1.DaemonSet) error {
	if os.Getenv("DRYRUN") != "" {
		os.Stdout.WriteString("---\n")
		daemonSet.SetGroupVersionKind(schema.FromAPIVersionAndKind("apps/v1", "DaemonSet"))
		b, err := yaml.Marshal(daemonSet)
		if err != nil {
			return err
		}
		_, err = os.Stdout.Write(b)
		return err
	}
	client := cs.AppsV1().DaemonSets(daemonSet.Namespace)
	existing, err := client.Get(ctx, daemonSet.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			_, err = client.Create(ctx, daemonSet, metav1.CreateOptions{})
			return err
		}
		return err
	}
	existing.Labels = daemonSet.Labels
	existing.Annotations = daemonSet.Annotations
	existing.Spec = daemonSet.Spec
	_, err = client.Update(ctx, existing, metav1.UpdateOptions{})
	return err
}

func DaemonSetAbsence(ctx context.Context, cs *kubernetes.Clientset, daemonSet *appsv1.DaemonSet) error {
	return Absence(func() error {
		propagationPolicy := metav1.DeletePropagationBackground
		return cs.AppsV1().DaemonSets(daemonSet.Namespace).Delete(ctx, daemonSet.Name, metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})
	})
}
