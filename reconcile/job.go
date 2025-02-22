package reconcile

import (
	"context"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"os"
	"sigs.k8s.io/yaml"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func Job(ctx context.Context, cs *kubernetes.Clientset, job *batchv1.Job) error {
	if os.Getenv("DRYRUN") != "" {
		os.Stdout.WriteString("---\n")
		job.SetGroupVersionKind(schema.FromAPIVersionAndKind("batch/v1", "Job"))
		b, err := yaml.Marshal(job)
		if err != nil {
			return err
		}
		_, err = os.Stdout.Write(b)
		return err
	}
	client := cs.BatchV1().Jobs(job.Namespace)
	existing, err := client.Get(ctx, job.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			_, err = client.Create(ctx, job, metav1.CreateOptions{})
			return err
		}
		return err
	}
	existing.Labels = job.Labels
	existing.Annotations = job.Annotations
	_, err = client.Update(ctx, existing, metav1.UpdateOptions{})
	return err
}

func JobAbsence(ctx context.Context, cs *kubernetes.Clientset, job *batchv1.Job) error {
	return Absence(func() error {
		propagationPolicy := metav1.DeletePropagationBackground
		return cs.BatchV1().Jobs(job.Namespace).Delete(ctx, job.Name, metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})
	})
}
