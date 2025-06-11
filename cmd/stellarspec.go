package main

import (
	"context"
	"fmt"
	"os"
	"stellarspec/internal/reviewer"

	"github.com/spf13/cobra"
)

var (
	// flag 变量
	apiServer     string
	model         string
	key           string
	confPath      string
	maxPool       int
	commitID      string
	promptFile    string
	thinkingChain bool
)

var rootCmd = &cobra.Command{
	Use:   "stellarspec",
	Short: "a code reviewer tool base on llm",
}

var reviewCmd = &cobra.Command{
	Use:   "review [file/directory]",
	Short: "do code review",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var reviewPath string
		if len(args) > 0 {
			reviewPath = args[0]
		} else {
			reviewPath = "."
		}
		engine := reviewer.NewReviewEngine(context.Background(), reviewPath)
		engine.Run()
	},
}

func init() {
	// 全局 flags (对所有命令生效)
	rootCmd.PersistentFlags().StringVar(&apiServer, "set-apiserver", "", "设置API服务器地址")
	rootCmd.PersistentFlags().StringVar(&model, "set-model", "", "设置LLM模型")
	rootCmd.PersistentFlags().StringVar(&key, "set-key", "", "设置API密钥")
	rootCmd.PersistentFlags().StringVar(&confPath, "conf", "", "指定配置文件路径")

	// 本地 flags (只对特定命令生效)
	reviewCmd.Flags().IntVar(&maxPool, "max-pool", 10, "并发操作上限")
	reviewCmd.Flags().StringVar(&commitID, "commit-id", "", "指定commit ID")
	reviewCmd.Flags().StringVar(&promptFile, "prompt-file", "", "自定义prompt文件路径")
	reviewCmd.Flags().BoolVar(&thinkingChain, "thinking-chain", false, "输出模型思考链")

	// 添加子命令
	rootCmd.AddCommand(reviewCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
