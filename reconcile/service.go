package reconcile

import (
	"context"
	"os"
	"sigs.k8s.io/yaml"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func Service(ctx context.Context, cs *kubernetes.Clientset, service *corev1.Service) error {
	if ctx.Value("dryrun") != nil {
		b, err := yaml.Marshal(service)
		if err != nil {
			return err
		}
		_, err = os.Stdout.Write(b)
		return err
	}
	client := cs.CoreV1().Services(service.Namespace)
	existing, err := client.Get(ctx, service.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			_, err = client.Create(ctx, service, metav1.CreateOptions{})
			return err
		}
		return err
	}
	clusterIP := existing.Spec.ClusterIP
	existing.Labels = service.Labels
	existing.Annotations = service.Annotations
	existing.Spec = service.Spec
	existing.Spec.ClusterIP = clusterIP
	_, err = client.Update(ctx, existing, metav1.UpdateOptions{})
	return err
}

func ServiceAbsence(ctx context.Context, cs *kubernetes.Clientset, service *corev1.Service) error {
	return Absence(func() error {
		return cs.CoreV1().Services(service.Namespace).Delete(ctx, service.Name, metav1.DeleteOptions{})
	})
}
