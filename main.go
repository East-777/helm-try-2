package main

import (
	"context"
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"log"
	"os"
	"path"
	//"helm.sh/helm/v3/pkg/action"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var settings = cli.New()

func init() {
	d, _ := os.Getwd()
	kubeConfig := path.Join(d, "config")

	actionConfig := new(action.Configuration)
	getter := &genericclioptions.ConfigFlags{
		KubeConfig: &kubeConfig,
	}
	kc := kube.New(getter)
	actionConfig.RESTClientGetter = getter
	actionConfig.KubeClient = kc

	settings.KubeConfig = kubeConfig

}

func main() {
	d, _ := os.Getwd()
	config, err := clientcmd.BuildConfigFromFlags("", path.Join(d, "config"))
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	namespaceList, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	for _, namespace := range namespaceList.Items {
		fmt.Println(namespace.Name)
	}
}
