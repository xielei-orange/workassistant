package container

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Manager interface {
	read() appsv1.DeploymentList
}

var (
	config, _ = clientcmd.BuildConfigFromFlags("", "./config/kubeconfig")
)

func Read() (appsv1.DeploymentList, error) {
	// 创建 Kubernetes 客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("Error creating Kubernetes client:", err)
		return appsv1.DeploymentList{}, err
	}
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceAll)
	deploymentList, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("获取deployment列表失败")
	}
	return *deploymentList, nil
}
