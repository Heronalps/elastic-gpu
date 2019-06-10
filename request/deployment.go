package request

import (
	"flag"
	"fmt"
	"path/filepath"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/retry"

	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

var clientset *kubernetes.Clientset
var err error

func init() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	// This config is credential information of kubernetes
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}

/*
QueryGPU queries number of GPU
Parameters:
	namespace
	deployment
*/
func QueryGPU(namespace string, deployment string) {
	deploymentsClient := clientset.AppsV1().Deployments(namespace)

	fmt.Println("Querying deployment...")

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := deploymentsClient.Get(deployment, metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get Deployment %v ", getErr))
		}
		numGpu := result.Spec.Template.Spec.Containers[0].Resources.Requests["nvidia.com/gpu"]
		fmt.Printf("Current Number of GPU is %v \n", numGpu.Value())
		return getErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Query failed: %v", retryErr))
	}
}

/*
Update number of GPU in the deployment
Parameters:
	namespace
	deployment - deployment name
	diff - difference of GPU number. Positive of negative integer
*/
func Update(namespace string, deployment string, diff int64) {
	deploymentsClient := clientset.AppsV1().Deployments(namespace)

	fmt.Println("Updating deployment...")

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := deploymentsClient.Get(deployment, metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of Deployment %v", getErr))
		}

		numGpu := result.Spec.Template.Spec.Containers[0].Resources.Requests["nvidia.com/gpu"]
		fmt.Printf("Current Number of GPU is %v \n", numGpu.Value())
		quant := resource.NewQuantity(diff, resource.DecimalSI)
		numGpu.Add(*quant)
		result.Spec.Template.Spec.Containers[0].Resources.Limits["nvidia.com/gpu"] = numGpu
		result.Spec.Template.Spec.Containers[0].Resources.Requests["nvidia.com/gpu"] = numGpu
		_, updateErr := deploymentsClient.Update(result)

		fmt.Printf("Updated Number of GPU is %v \n", numGpu.Value())
		return updateErr
	})

	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
}
