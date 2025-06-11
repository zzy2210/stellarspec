# StellarSpec

🌟 **星鉴** - 智能代码审查助手

一款基于 LLM 大模型的智能本地代码审查工具，采用 Go 语言开发，基于 [Eino 框架](https://github.com/cloudwego/eino) 构建。StellarSpec 能够自动检测 Git 仓库中的代码变更，并利用 AI 大模型进行深度代码审查，帮助开发者发现潜在问题、优化代码质量。

## ✨ 特性

- 🔍 **智能代码分析**: 基于 LLM 大模型进行深度代码审查
- 🚀 **高性能并发**: 支持多文件并发处理，提升审查效率
- 📊 **Git 集成**: 自动检测工作区和暂存区的代码变更
- 🎯 **多语言支持**: 支持 Go、JavaScript、TypeScript、Python 等 20+ 种编程语言
- 🛠️ **灵活配置**: 支持多种 API 服务器和模型配置
- 📝 **专业报告**: 生成详细的 Markdown 格式审查报告
- 🧠 **思维链模式**: 可选输出模型的详细分析推理过程

## 🚀 快速开始

### 安装

```bash
# 从源码编译
git clone https://github.com/your-username/stellarspec.git
cd stellarspec
go build -o stellarspec cmd/stellarspec.go
```

### 配置

首次使用需要配置 API 服务器、模型和密钥：

```bash
# 设置 API 服务器地址
stellarspec --set-apiserver https://api.siliconflow.cn/v1/

# 设置 LLM 模型
stellarspec --set-model deepseek-chat

# 设置 API 密钥
stellarspec --set-key sk-xxxxxxxxxxxxx
```

配置文件将自动保存到 `$HOME/.stellarspec/cnf`

### 基础使用

```bash
# 审查当前目录的所有变更
stellarspec review

# 审查指定文件
stellarspec review main.go

# 审查指定目录
stellarspec review ./src

# 查看帮助
stellarspec --help
```

审查完成后，将在当前目录生成 `code-review.md` 报告文件。

## 📖 详细使用说明

### 配置管理

```bash
# 使用自定义配置文件
stellarspec --conf /path/to/custom/config review

# 一次性设置多个配置项
stellarspec --set-apiserver https://api.openai.com/v1 \
           --set-model gpt-4 \
           --set-key sk-xxxxxx
```

### 审查选项

```bash
# 指定并发数量（默认 10）
stellarspec review . --max-pool 20

# 审查特定 commit 的变更
stellarspec review --commit-id 1bacd3f

# 启用思维链模式，查看详细分析过程
stellarspec review --thinking-chain

# 使用自定义 prompt 模板
stellarspec review --prompt-file custom_prompt.txt

# 组合使用多个选项
stellarspec review main.go --thinking-chain --max-pool 5
```

### 支持的文件类型

StellarSpec 支持以下编程语言的代码审查：

| 语言 | 扩展名 | 语言 | 扩展名 |
|------|--------|------|--------|
| Go | `.go` | Python | `.py` |
| JavaScript | `.js`, `.jsx` | Java | `.java` |
| TypeScript | `.ts`, `.tsx` | C/C++ | `.c`, `.cpp`, `.cc` |
| Rust | `.rs` | PHP | `.php` |
| Ruby | `.rb` | Swift | `.swift` |
| Kotlin | `.kt` | Scala | `.scala` |
| Shell | `.sh` | SQL | `.sql` |
| HTML | `.html`, `.htm` | CSS | `.css` |
| YAML | `.yaml`, `.yml` | JSON | `.json` |
| XML | `.xml` | Markdown | `.md` |

## 🔧 高级功能

### 并发处理

StellarSpec 支持多文件并发审查，显著提升大型项目的处理速度：

```bash
# 设置并发数为 20（适合大型项目）
stellarspec review . --max-pool 20
```

### 思维链分析

启用思维链模式可以查看 AI 模型的详细分析过程：

```bash
stellarspec review --thinking-chain
```

### 自定义 Prompt

可以使用自定义 prompt 文件来定制审查重点：

```bash
# 创建自定义 prompt 文件
echo "重点关注安全性和性能问题" > security_prompt.txt

# 使用自定义 prompt
stellarspec review --prompt-file security_prompt.txt
```

## 📊 审查报告

StellarSpec 生成的审查报告包含以下内容：

- **文件信息**: 文件路径、类型、审查时间
- **问题分析**: 发现的代码问题和潜在风险
- **改进建议**: 具体的优化建议
- **修改方案**: 详细的代码改进方案

报告示例：

```markdown
## 文件审查报告

**文件路径**: internal/reviewer/reviewer.go  
**文件类型**: Go  
**审查时间**: 2025-06-11 19:05:11

### 审查结果

### 发现的问题
1. **错误处理不完整**: 某些函数缺少适当的错误处理
2. **性能优化**: 可以使用更高效的数据结构

### 改进建议
1. 添加完整的错误检查和处理逻辑
2. 考虑使用 sync.Pool 优化内存分配

### 修改方案
建议在关键函数中添加 defer 语句确保资源正确释放...
```

## 🛠️ 技术架构

### 核心组件

- **Eino 框架**: 提供 LLM 模型集成和工作流编排
- **Git 集成**: 基于 go-git 库实现版本控制集成
- **并发处理**: 使用 goroutine 和信号量实现高效并发
- **配置管理**: 基于 INI 格式的灵活配置系统

### 项目结构

```
stellarspec/
├── cmd/                    # 命令行入口
│   └── stellarspec.go     # 主程序和命令定义
├── internal/              # 内部模块
│   ├── model/             # 数据模型
│   │   └── conf/          # 配置管理
│   └── reviewer/          # 核心审查引擎
├── go.mod                 # Go 模块定义
├── go.sum                 # 依赖版本锁定
└── README.md             # 项目文档
```

## 🤝 贡献指南

欢迎贡献代码！请遵循以下步骤：

1. Fork 本仓库
2. 创建特性分支: `git checkout -b feature/amazing-feature`
3. 提交更改: `git commit -m 'Add amazing feature'`
4. 推送分支: `git push origin feature/amazing-feature`
5. 提交 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 🙏 致谢

- [Eino](https://github.com/cloudwego/eino) - 强大的 LLM 应用框架
- [Cobra](https://github.com/spf13/cobra) - Go CLI 库
- [go-git](https://github.com/go-git/go-git) - Git 实现
- [go-diff](https://github.com/sergi/go-diff) - 差异比较库

## 📞 支持

如果你遇到问题或有建议，请：

- 提交 [Issue](https://github.com/your-username/stellarspec/issues)
- 发起 [Discussion](https://github.com/your-username/stellarspec/discussions)

---

⭐ 如果这个项目对你有帮助，请给我们一个 Star！