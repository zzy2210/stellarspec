# StellarSpec

[简体中文](README.md) | [English](README.en.md)

🌟 **星鉴** - 智能代码审查助手

一款基于 LLM 大模型的智能本地代码审查工具，采用 Go 语言开发，基于 [Eino 框架](https://github.com/cloudwego/eino) 构建。StellarSpec 能够自动检测 Git 仓库中的代码变更，并利用 AI 大模型进行深度代码审查，帮助开发者发现潜在问题、优化代码质量。

## ✨ 特性概览

已实现（当前可用）
- 🔍 智能代码分析：基于 LLM 的代码审查
- 🚀 并发处理：固定 10 并发执行
- 📊 Git 集成：自动检测工作区与暂存区变更（相对 HEAD）
- 🎯 多语言识别：按文件扩展名识别 20+ 语言类型
- 🌐 国际化：支持中文/英文报告与提示词
- 🛠️ 配置管理：API Server/模型/密钥/语言持久化到本地配置
- 📝 报告输出：按文件生成 Markdown 追加式报告 `code-review.md`

规划中（待完善/接线中）
- ⏱ 并发上限可配：`--max-pool`（目前参数保留，未生效，默认 10）
- 🔖 指定提交审查：`--commit-id`（目前未生效）
- 🧠 思维链输出：`--thinking-chain`（目前未生效）
- 📝 自定义 Prompt：`--prompt-file`（目前未生效）
- 🎯 范围过滤：直接审查指定文件/子目录（当前仅按仓库变更进行审查）

## 🚀 快速开始

### 安装

推荐使用 `Makefile` 进行构建：

```bash
# 克隆仓库
git clone https://github.com/your-username/stellarspec.git
cd stellarspec

# 构建（产物位于 build/stellar）
make build
```

构建完成后，二进制位于 `build/stellar`。如需全局使用，可复制到 `PATH` 目录，例如：

```bash
sudo cp build/stellar /usr/local/bin/stellar
```

### 配置

首次使用需要配置 API 服务器、模型和密钥。请将 `stellar` 替换为你实际的二进制文件名（如果已更改）。

```bash
# 设置 API 服务器地址
stellar --set-apiserver https://api.siliconflow.cn/v1/

# 设置 LLM 模型
stellar --set-model deepseek-chat

# 设置 API 密钥
stellar --set-key sk-xxxxxxxxxxxxx

# 设置语言（可选，默认为中文）
stellar --set-lang zh  # 中文
stellar --set-lang en  # 英文
```

配置文件将自动保存到 `$HOME/.stellarspec/cnf`

### 基础使用

```bash
# 审查当前 Git 仓库中的变更（相对 HEAD）
./build/stellar review

# 使用已安装的二进制
stellar review

# 查看帮助（或使用 make run）
stellar --help
```

说明
- 目前仅对“仓库变更”进行审查；传入 `review [path]` 将作为仓库路径尝试打开 Git 仓库。
- 直接对“单个文件/子目录”进行精确审查尚未接线，规划中。
- 审查完成后将在工作目录生成 `code-review.md` 报告文件。

## 📖 详细使用说明

### 配置管理

```bash
# 使用自定义配置文件
stellar --conf /path/to/custom/config review

# 一次性设置多个配置项
stellar --set-apiserver https://api.openai.com/v1 \
           --set-model gpt-4 \
           --set-key sk-xxxxxx \
           --set-lang en

# 切换语言设置
stellar --set-lang zh  # 切换为中文
stellar --set-lang en  # 切换为英文
```

### 审查选项（占位，规划中）

以下选项已在 CLI 中预留，但暂未在引擎内生效，接线后方可使用：

```bash
# 指定并发数量（默认 10）
stellar review . --max-pool 20

# 审查特定 commit 的变更
stellar review --commit-id 1bacd3f

# 启用思维链模式，查看详细分析过程
stellar review --thinking-chain

# 使用自定义 prompt 模板
stellar review --prompt-file custom_prompt.txt
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

当前并发度固定为 10。`--max-pool` 选项已预留，接线后支持自定义并发上限。

### 思维链分析

`--thinking-chain` 预留中，未来可查看模型更详细的分析过程（默认关闭）。

### 自定义 Prompt

`--prompt-file` 预留中，后续将支持通过文件自定义审查提示词与重点。

## 📊 审查报告

报告文件为 `code-review.md`，按文件输出包含：

- 文件信息：路径、识别的语言类型、时间戳
- 审查结果：模型给出的结论文本（按语言切换模板）

说明：示例中的“发现的问题/改进建议/修改方案”等分节由模型生成文本中承载，模板暂未强制分节。

## 🛠️ 技术架构

### 核心组件

- **Eino 框架**: 提供 LLM 模型集成和工作流编排
- **Git 集成**: 基于 go-git 库实现版本控制集成
- **并发处理**: 使用 goroutine 和信号量实现高效并发
- **配置管理**: 基于 INI 格式的灵活配置系统

### 项目结构

```
stellarspec/
├── cmd/                    # Cobra CLI 入口
│   └── stellarspec.go
├── internal/
│   ├── model/
│   │   └── conf/          # INI 配置读写
│   └── reviewer/          # 变更收集 / 并发执行 / 报告输出
├── build/                 # 构建产物（git 忽略）
├── go.mod
├── go.sum
└── README.md
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

- [ethereal14](https://github.com/ethereal14) - my good friend

## 📞 支持

如果你遇到问题或有建议，请：

- 提交 Issue
- 发起 Discussion

---

⭐ 如果这个项目对你有帮助，请给我们一个 Star！
