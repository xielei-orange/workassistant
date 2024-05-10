package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"tools/cmd/container"
	"tools/cmd/migrate"
)

var cobraCmd = &cobra.Command{
	Use:   "tools",
	Short: "tools",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			// 提示语
			tip()
			return errors.New("至少需要输入一位参数")
		}
		return nil
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		tip()
	}}

func tip() {
	fmt.Println("欢迎使用tools工具，该工具主要解决日常重复工作")
	fmt.Println("可使用-h或--help获取命令使用详情")
}

func init() {
	cobraCmd.AddCommand(container.StartCmd)
	cobraCmd.AddCommand(migrate.StartCmd)
}
func Execute() {
	if err := cobraCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
