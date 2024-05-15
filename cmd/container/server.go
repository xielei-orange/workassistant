package container

import (
	"fmt"
	"os"
	"workassistant/common"
	"workassistant/common/container"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	actionType    string
	namespace     string
	namespaceList []string
	file          string
	StartCmd      = &cobra.Command{
		Use:     "container",
		Short:   "容器管理工具",
		Long:    "可以对容器进行CURD操作",
		Example: "workassistant container update -f 2023.excel",
		PreRun: func(cmd *cobra.Command, args []string) {
			fmt.Printf("正在读取flag参数，稍等............\n")
			if file != "" {
				fmt.Printf("正在读取flag参数，稍等............\n")
				checkFields := []string{"命名空间", "服务", "版本号"}
				rows, err := common.ReadExcel("config/tmp/20240401.xlsx", checkFields)
				if err != nil {
					panic(err)
				}
				for index := range rows {
					//fmt.Println(index)
					//fmt.Println(value)
					index += 1
				}
			}
			fmt.Println(namespaceList)
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Run方法")
		},
	}
)

func init() {
	StartCmd.Flags().StringVarP(&actionType, "actionType", "t", "", "操作类型,CURD")
	StartCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "命名空间")
	StartCmd.Flags().StringVarP(&file, "file", "f", "", "deployment文件")
	config := common.GetConfig()
	fmt.Println(config.GetString("settings.database.source"))
	namespaceListArray, err := container.GetNamespaceList()
	if err != nil {
		panic(err)
	}
	for _, n := range namespaceListArray.Items {
		namespaceList = append(namespaceList, n.ObjectMeta.Name)
	}
	container.CreateDeploymentByNamespace("demo", "cloud", "nginx:latest", 3, map[string]string{"app": "nginx"})
}

func VerifyFlag(args []string) {
	// 本地命令执行的类型，CURD，create,update,read,delete
	if actionType == "" {
		color.Red("请输入flag,详情可使用 command container -h[--help]")
		os.Exit(-1)
	}
	//如果参数传入-f，忽略其他参数
	if file != "" {
		checkFields := []string{"命名空间", "服务", "版本号"}
		rows, err := common.ReadExcel("./20240401.xlsx", checkFields)
		if err != nil {
			panic(err)
		}
		for index, value := range rows {
			fmt.Println(index)
			fmt.Println(value)
		}
	} else {
		if namespace == "" {
			color.Red("请输入flag,详情可使用 command container -h[--help]")
			os.Exit(-1)
		}
		if len(args) == 0 {
			color.Red("请提供deployment名称")
			os.Exit(-1)
		}
	}

}
