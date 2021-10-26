package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os/user"
)

func main() {
	//// 1.创建路由
	//r := gin.Default()
	//// 2.绑定路由规则，执行的函数
	//// gin.Context，封装了request和response
	//r.GET("/", func(c *gin.Context) {
	//	c.String(http.StatusOK, "hello World!")
	//})
	//// 3.监听端口，默认在8080
	//// Run("里面不指定端口号默认为8080")
	//r.Run(":7000")

	homePath := GetHomePath()

	config, err := clientcmd.BuildConfigFromFlags("", fmt.Sprintf("%s/.kube/config", homePath))
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	nodeList, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	for _, node := range nodeList.Items {
		fmt.Println(node.Name)
	}
}

func GetHomePath() string {
	current, err := user.Current()
	if err == nil {
		return current.HomeDir
	}

	return ""
}
