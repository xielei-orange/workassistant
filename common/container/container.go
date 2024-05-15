package container

import (
	"context"
	"fmt"

	appapi "k8s.io/api/apps/v1"
	coreapi "k8s.io/api/core/v1"
	metaapi "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Manager interface {
	read() appapi.DeploymentList
}

var (
	config, _ = clientcmd.BuildConfigFromFlags("", "./config/tmp/kubeconfig")
)

// func Read() (appapi.DeploymentList, error) {
// 	// 创建 Kubernetes 客户端
// 	clientset, err := kubernetes.NewForConfig(config)
// 	if err != nil {
// 		fmt.Println("Error creating Kubernetes client:", err)
// 		return appapi.DeploymentList{}, err
// 	}
// 	deploymentsClient := clientset.AppsV1().Deployments(coreapi.NamespaceAll)
// 	deploymentList, err := deploymentsClient.List(context.TODO(), metaapi.ListOptions{})
// 	if err != nil {
// 		fmt.Println("获取deployment列表失败")
// 	}
// 	return *deploymentList, nil
// }

func GetNamespaceList() (*coreapi.NamespaceList, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("Error creating Kubernetes client:", err)
		return &coreapi.NamespaceList{}, err
	}
	namespaceList, err := clientset.CoreV1().Namespaces().List(context.TODO(), metaapi.ListOptions{})
	if err != nil {
		fmt.Println("Error creating Kubernetes client:", err)
		return &coreapi.NamespaceList{}, err
	}
	return namespaceList, nil
}

// default
func GetDeploymentByNamespace(namespaces ...string) ([]*appapi.DeploymentList, error) {
	if len(namespaces) == 0 {
		namespaces = []string{coreapi.NamespaceAll}
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return []*appapi.DeploymentList{}, err
	}
	deploymentList := []*appapi.DeploymentList{}
	for _, namespace := range namespaces {
		d, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metaapi.ListOptions{})
		if err != nil {
			panic(err)
		}
		for _, d := range d.Items {
			fmt.Println(d.ObjectMeta.Name)
		}
		deploymentList = append(deploymentList, d)

	}
	return deploymentList, nil
}

func CreateDeploymentByNamespace(name, namespace, image string, replica int32, podLabels map[string]string) (*appapi.Deployment, error) {
	// DeleteDeploymentByNamespace(namespace, name)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return &appapi.Deployment{}, err
	}
	yaml := appapi.Deployment{
		ObjectMeta: metaapi.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		}, Spec: appapi.DeploymentSpec{
			Replicas: &replica,
			Selector: &metaapi.LabelSelector{
				MatchLabels: podLabels},
			Template: coreapi.PodTemplateSpec{
				ObjectMeta: metaapi.ObjectMeta{
					Labels: podLabels,
				},
				Spec: coreapi.PodSpec{
					Containers: []coreapi.Container{{
						Name:  name,
						Image: image,
						Ports: []coreapi.ContainerPort{{ContainerPort: 80}},
						ReadinessProbe: &coreapi.Probe{
							ProbeHandler: coreapi.ProbeHandler{HTTPGet: &coreapi.HTTPGetAction{
								Path: "/",
								Port: intstr.FromInt(80),
							},
							}, SuccessThreshold: 10,
						},
					}},
				},
			},
		},
	}

	createDeploymentObject, err := clientset.AppsV1().Deployments(namespace).Create(context.TODO(), &yaml, metaapi.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Create Deployment %s/%s success\n", namespace, namespace)

	deploymenInterface := clientset.AppsV1().Deployments(namespace)
	// timeout := int64(10)
	watcher, err := deploymenInterface.Watch(context.TODO(), metaapi.ListOptions{FieldSelector: "metadata.name=" + name, Watch: true})
	// LabelSelector: labels.FormatLabels(podLabels)})

	if err != nil {
		panic(err)
	}
	defer watcher.Stop()
	// 创建chan接收watch
	eventChan := watcher.ResultChan()
	// e := <-watcher.ResultChan()
	// fmt.Println(e)
	// for true {
	// 	e := <-watcher.ResultChan()
	// 	fmt.Println(e.Type, e.Object)

	// }
	// go routine处理watch
	for event := range eventChan {
		deploy, ok := event.Object.(*appapi.Deployment)
		if !ok {
			fmt.Println("Unexpected object type")
			continue
		}
		// 打印 Deployment 的副本数量和可用副本数量
		fmt.Printf("deployment更新中 %d/%d\n", deploy.Status.AvailableReplicas, replica)

		// 检查是否所有副本都已经可用
		if replica == deploy.Status.AvailableReplicas {
			// fmt.Printf("deployment更新中 %d/%d\n", deploy.Status.AvailableReplicas, replica)
			fmt.Println("All replicas are available.")
			break
		}
	}
	return createDeploymentObject, nil
}

func DeleteDeploymentByNamespace(namespace, name string) int {
	fmt.Printf("执行deployment删除指令，删除%s/%s\n", namespace, name)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	resp := clientset.AppsV1().Deployments(namespace).Delete(context.TODO(), name, metaapi.DeleteOptions{})
	fmt.Println(resp)
	if resp != nil {
		fmt.Println("deployment已删除")
		return 200
	}
	fmt.Println("deployment不存在")
	return 500
}
