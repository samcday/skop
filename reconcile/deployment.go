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

func Deployment(ctx context.Context, cs *kubernetes.Clientset, deployment *appsv1.Deployment) error {
	if os.Getenv("DRYRUN") != "" {
		os.Stdout.WriteString("---\n")
		deployment.SetGroupVersionKind(schema.FromAPIVersionAndKind("apps/v1", "Deployment"))
		b, err := yaml.Marshal(deployment)
		if err != nil {
			return err
		}
		_, err = os.Stdout.Write(b)
		return err
	}
	client := cs.AppsV1().Deployments(deployment.Namespace)
	existing, err := client.Get(ctx, deployment.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			_, err = client.Create(ctx, deployment, metav1.CreateOptions{})
			return err
		}
		return err
	}
	existing.Labels = deployment.Labels
	existing.Annotations = deployment.Annotations
	existing.Spec = deployment.Spec
	_, err = client.Update(ctx, existing, metav1.UpdateOptions{})
	return err
}

func DeploymentAbsence(ctx context.Context, cs *kubernetes.Clientset, deployment *appsv1.Deployment) error {
	return Absence(func() error {
		propagationPolicy := metav1.DeletePropagationBackground
		return cs.AppsV1().Deployments(deployment.Namespace).Delete(ctx, deployment.Name, metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})
	})
}
