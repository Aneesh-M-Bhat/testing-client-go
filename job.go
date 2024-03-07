package main

import (
	"context"
	"fmt"
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

func createJob(clientset *kubernetes.Clientset) {
	jobsClient := clientset.BatchV1().Jobs(corev1.NamespaceDefault)
	backOffLimit := int32(4)
	job := &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "batch/v1",
			Kind:       "Job",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-job",
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "pi",
							Image:   "perl:5.34.0",
							Command: strings.Split("perl,-Mbignum=bpi,-wle,print bpi(2000)", ","),
						},
					},
					RestartPolicy: "Never",
				},
			},
			BackoffLimit: &backOffLimit,
		},
	}

	fmt.Println("Creating Job...")
	result, err := jobsClient.Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created Job %q.\n", result.GetObjectMeta().GetName())

	prompt()
	fmt.Println("Updating Job...")
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := jobsClient.Get(context.TODO(), "demo-job", metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("failed to get latest version of Job: %v", getErr))
		}

		// result.Spec.Replicas = int32Ptr(3)
		backOffLimit = 6
		result.Spec.BackoffLimit = &backOffLimit
		_, updateErr := jobsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})

	if retryErr != nil {
		panic(fmt.Errorf("update failed: %v", retryErr))
	}
	fmt.Println("Updated job")

	prompt()
	fmt.Printf("Listing jobs in namespace %q:\n", corev1.NamespaceDefault)
	list, err := jobsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s \n", d.Name)
	}

	prompt()
	fmt.Printf("Deleting job...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := jobsClient.Delete(context.TODO(), "demo-job", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted job.")
	prompt()
}
