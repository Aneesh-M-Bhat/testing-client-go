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

func createCronJob(clientset *kubernetes.Clientset) {
	cronjobsClient := clientset.BatchV1().CronJobs(corev1.NamespaceDefault)
	cronjob := &batchv1.CronJob{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "batch/v1",
			Kind:       "CronJob",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-cronjob",
		},
		Spec: batchv1.CronJobSpec{
			Schedule: "* * * * *",
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							RestartPolicy: corev1.RestartPolicyOnFailure,
							Containers: []corev1.Container{
								{
									Name:            "hello",
									Image:           "busybox:1.28",
									ImagePullPolicy: "IfNotPresent",
									Command:         strings.Split("/bin/sh,-c,date; echo Hello from the Kubernetes cluster", ","),
								},
							},
						},
					},
				},
			},
		},
	}

	fmt.Println("Creating CronJob...")
	result, err := cronjobsClient.Create(context.TODO(), cronjob, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created CronJob %q.\n", result.GetObjectMeta().GetName())

	prompt()
	fmt.Println("Updating CronJob...")
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := cronjobsClient.Get(context.TODO(), "demo-cronjob", metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("failed to get latest version of cronjob: %v", getErr))
		}
		result.Spec.JobTemplate.Spec.Template.Spec.RestartPolicy = corev1.RestartPolicyNever
		_, updateErr := cronjobsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})

	if retryErr != nil {
		panic(fmt.Errorf("update failed: %v", retryErr))
	}
	fmt.Println("Updated cronjob")

	prompt()
	fmt.Printf("Listing cronjobs in namespace %q:\n", corev1.NamespaceDefault)
	list, err := cronjobsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s \n", d.Name)
	}

	prompt()
	fmt.Printf("Deleting cronjob...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := cronjobsClient.Delete(context.TODO(), "demo-cronjob", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted cronjob.")
	prompt()
}
