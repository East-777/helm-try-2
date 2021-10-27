package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"path"

	//"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	d, _ := os.Getwd()
	config, err := clientcmd.BuildConfigFromFlags("https://120.25.216.113:6443", path.Join(d, "config"))
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	version, err := clientset.Discovery().ServerVersion()
	fmt.Printf("%v", version)
	fmt.Printf("%v", err)
	namespaceList, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	for _, namespace := range namespaceList.Items {
		fmt.Println(namespace.Name)
	}
}
