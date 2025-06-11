package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	config "stellarspec/internal/model/conf"
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
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		handleConfigFlags()
	},
	Run: func(cmd *cobra.Command, args []string) {
		// 如果只是设置配置，不需要额外操作
		// 配置已经在 PersistentPreRun 中处理了
		if apiServer != "" || model != "" || key != "" {
			fmt.Println("配置设置完成")
			return
		}

		// 如果没有提供任何参数，显示帮助信息
		cmd.Help()
	},
}

// 获取默认路径
func getDefaultConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// 如果获取失败，使用当前目录
		return "config.ini"
	}

	configDir := filepath.Join(homeDir, ".stellarspec")
	return filepath.Join(configDir, "cnf")
}

// 确保配置目录存在
func ensureConfigDir(configPath string) error {
	configDir := filepath.Dir(configPath)

	// 检查目录是否存在
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		// 目录不存在，创建目录
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("create config directory failed: %v", err)
		}
	}

	return nil
}

// 添加处理配置的函数
func handleConfigFlags() {
	// 默认配置文件路径
	configPath := getDefaultConfigPath()
	if confPath != "" {
		configPath = confPath
	}

	// 确保配置目录存在
	if err := ensureConfigDir(configPath); err != nil {
		fmt.Printf("创建配置目录失败: %v\n", err)
		os.Exit(1)
	}

	// 如果有设置 API 服务器
	if apiServer != "" {
		if err := config.SaveAPIServer(apiServer, configPath); err != nil {
			fmt.Printf("保存 API 服务器配置失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("API 服务器已设置为: %s\n", apiServer)
	}

	// 如果有设置模型
	if model != "" {
		if err := config.SaveModel(model, configPath); err != nil {
			fmt.Printf("保存模型配置失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("模型已设置为: %s\n", model)
	}

	// 如果有设置密钥
	if key != "" {
		if err := config.SaveKey(key, configPath); err != nil {
			fmt.Printf("保存密钥配置失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("API 密钥已设置\n")
	}
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
		configPath := getDefaultConfigPath()
		if confPath != "" {
			configPath = confPath
		}

		engine := reviewer.NewReviewEngine(context.Background(), reviewPath)
		baseConf, err := config.LoadFile(configPath)
		if err != nil {
			fmt.Println("load config file failed: err= %v", err)
		}
		engine.CreateModel(baseConf)
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
