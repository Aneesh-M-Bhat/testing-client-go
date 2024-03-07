package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

func createPod(clientset *kubernetes.Clientset) {
	podsClient := clientset.CoreV1().Pods(corev1.NamespaceDefault)
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-pod",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx",
					Image: "nginx",
				},
			},
		},
	}

	res, err := podsClient.Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created Pod %q.\n", res.GetObjectMeta().GetName())

	prompt()
	fmt.Println("Updating Pod...")
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := podsClient.Get(context.TODO(), "demo-pod", metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("failed to get latest version of Deployment: %v", getErr))
		}
		result.Spec.Containers[0].Image = "busybox"
		_, updateErr := podsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})

	if retryErr != nil {
		panic(fmt.Errorf("update failed: %v", retryErr))
	}
	fmt.Println("Updated pod")

	prompt()
	fmt.Printf("Listing pods in namespace %q:\n", corev1.NamespaceDefault)
	list, err := podsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s \n", d.Name)
	}

	prompt()
	fmt.Printf("Deleting pod...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := podsClient.Delete(context.TODO(), "demo-pod", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted pod.")
	prompt()
}
