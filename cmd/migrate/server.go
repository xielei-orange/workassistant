package migrate

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"tools/common/container"
	"tools/common/database"
	"tools/models"
)

var (
	actionType string
	namespace  string
	StartCmd   = &cobra.Command{
		Use:     "migrate",
		Example: "tools migrate",
		Run: func(cmd *cobra.Command, args []string) {
			if err := database.RunMigrations(); err != nil {
				color.Yellow("数据迁移失败")
			}
			color.Green("数据迁移成功")
			deploymentList, err := container.Read()
			if err != nil {
				panic(err)
				fmt.Println("获取deployment失败")
			}
			var dataList []models.Container
			for _, d := range deploymentList.Items {
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

		},
	}
)

func resourceCondition(ResourceList v1.ResourceRequirements, ResourceName v1.ResourceName, resourceKind string) string {
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
