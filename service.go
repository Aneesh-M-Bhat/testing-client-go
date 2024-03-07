package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

func createService(clientset *kubernetes.Clientset) {
	servicesClient := clientset.CoreV1().Services(corev1.NamespaceDefault)
	selectors := make(map[string]string)
	selectors["app.kubernetes.io/name"] = "MyApp"
	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-service",
		},
		Spec: corev1.ServiceSpec{
			Selector: selectors,
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       80,
					TargetPort: intstr.FromInt(9376),
				},
			},
		},
	}

	fmt.Println("Creating Service...")
	result, err := servicesClient.Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created Service %q.\n", result.GetObjectMeta().GetName())
	prompt()

	fmt.Println("Updating Service...")
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := servicesClient.Get(context.TODO(), "demo-service", metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("failed to get latest version of Service: %v", getErr))
		}

		result.Spec.Ports[0].TargetPort = intstr.FromInt(8000)
		_, updateErr := servicesClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})

	if retryErr != nil {
		panic(fmt.Errorf("update failed: %v", retryErr))
	}
	fmt.Println("Updated service")
	prompt()

	fmt.Printf("Listing services in namespace %q:\n", corev1.NamespaceDefault)
	list, err := servicesClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s \n", d.Name)
	}
	prompt()

	fmt.Printf("Deleting service...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := servicesClient.Delete(context.TODO(), "demo-service", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted service")
	prompt()
}
