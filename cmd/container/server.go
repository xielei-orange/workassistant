package container

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"tools/common"
)

var (
	actionType string
	namespace  string
	file       string
	StartCmd   = &cobra.Command{
		Use:     "container",
		Short:   "容器管理工具",
		Long:    "可以对容器进行CURD操作",
		Example: "tools container update -f 2023.excel",
		PreRun: func(cmd *cobra.Command, args []string) {
			fmt.Printf("正在读取flag参数，稍等............\n")
			if file != "" {
				fmt.Printf("正在读取flag参数，稍等............\n")
				checkFields := []string{"命名空间", "服务", "版本号"}
				rows, err := common.ReadExcel("./20240401版本信息(1).xlsx", checkFields)
				if err != nil {
					panic(err)
				}
				for index, _ := range rows {
					//fmt.Println(index)
					//fmt.Println(value)
					index += 1
				}
			}
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
}
func verifyFlag(args []string) {
	// 本地命令执行的类型，CURD，create,update,read,delete
	if actionType == "" {
		color.Red("请输入flag,详情可使用 command container -h[--help]")
		os.Exit(-1)
	}
	//如果参数传入-f，忽略其他参数
	if file != "" {
		checkFields := []string{"命名空间", "服务", "版本号"}
		rows, err := common.ReadExcel("./20240401版本信息(1).xlsx", checkFields)
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
