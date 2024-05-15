package migrate

import (
	"fmt"
	"workassistant/common/container"
	"workassistant/common/database"
	"workassistant/models"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	coreapi "k8s.io/api/core/v1"
)

var (
	// actionType string
	// namespace  string
	StartCmd = &cobra.Command{
		Use:     "migrate",
		Example: "workassistant migrate",
		Run: func(cmd *cobra.Command, args []string) {
			if err := database.RunMigrations(); err != nil {
				color.Yellow("数据迁移失败")
			}
			color.Green("数据迁移成功")
			deploymentList, err := container.GetDeploymentByNamespace(coreapi.NamespaceAll)
			if err != nil {
				panic(err)
			}
			var dataList []models.Container
			for _, deploymentListNamespace := range deploymentList {
				for _, d := range deploymentListNamespace.Items {
					container := d.Spec.Template.Spec.Containers[0]
					cpuReq := resourceCondition(container.Resources, "cpu", "request")
					cpuLimit := resourceCondition(container.Resources, "cpu", "limit")
					memReq := resourceCondition(container.Resources, "memory", "request")
					memLimit := resourceCondition(container.Resources, "memory", "limit")
					data := models.Container{
						Name:         d.Name,
						Namespace:    d.Namespace,
						CpuRequest:   cpuReq,
						CpuLimit:     cpuLimit,
						MemReq:       memReq,
						MemLimit:     memLimit,
						Replicas:     *d.Spec.Replicas,
						ImageVersion: container.Image,
						CreateAt:     d.CreationTimestamp.Time,
					}
					dataList = append(dataList, data)
				}
				resp := database.DB.Create(&dataList)
				if resp.Error != nil {
					fmt.Println("数据插入失败")
				}
			}
		},
	}
)

func resourceCondition(ResourceList coreapi.ResourceRequirements, ResourceName coreapi.ResourceName, resourceKind string) string {
	switch resourceKind {
	case "request":
		if value, exist := ResourceList.Requests[ResourceName]; exist {
			return value.String()
		}
		return "无限制"
	case "limit":
		if value, exist := ResourceList.Limits[ResourceName]; exist {
			return value.String()
		}
		return "无限制"
	}
	return "unknown"
}
