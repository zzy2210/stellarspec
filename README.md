# StellarSpec
星鉴是一款基于 LLM 大模型的本地code review 工具，使用go语言开发，基于eino框架

## 使用说明

### 配置
`stellarspec --set-apiserver https://api.siliconflow.cn/v1/ `
设置调用apiserver地址，这将会在 $HOME/.stellarspec/ 下创建conf文件 
`stellarspec --set-model xxx `
设置llm 模型
`stellarspec --set-key sk-xxxxx `
设置密钥

当然，你可以保持多个conf文件，在需要的时候直接调用
`stellarspec --conf your_path review`

### 使用

`stellarspec review `
这将调用配置文件，获取当前目录当前缓冲区中的所有变更，提交给llm大模型。
并且将创建文件：code_reviw.md

`stellarspec review a.go`
这将比较该文件的缓冲区变更，并且提交

`stellarspec review work`
这将获取work文件夹下全部变更，并且提交

`stellarspec review --commit-id 1bacd3`
获取指定commit的变更

`stellarspec review . --max-pool 20`
这将启用并发操作，并发上线设置为20 （默认为10）

`stellarspec review --thinking-chain`
输出模型的思考链，显示详细的分析过程和推理步骤

`stellarspec review a.go --prompt-file my_prompt.txt --thinking-chain`
组合使用：使用自定义 prompt 并输出思考链