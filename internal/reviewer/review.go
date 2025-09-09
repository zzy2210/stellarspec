package reviewer

import (
    "fmt"
    "path/filepath"
    "time"

    "github.com/cloudwego/eino/components/prompt"
    "github.com/cloudwego/eino/compose"
    "github.com/cloudwego/eino/schema"
    "github.com/fatih/color"
)

// reviewSingleFile 对单个文件变更进行审查并写入报告
func (e *Engine) reviewSingleFile(d gitDiff) error {
    if e.chatModel == nil {
        return fmt.Errorf("chat model is nil")
    }

    g := compose.NewGraph[map[string]any, *schema.Message]()
    ext := filepath.Ext(d.FilePath)

    // 打印审查开始
    color.Cyan("▶ review: %s\n", d.FilePath)
    start := time.Now()

    // 根据语言设置选择提示词
    var systemTpl string
    if e.cfg.Language == "en" {
        systemTpl = fmt.Sprintf("You are a %s development expert. You will provide code review conclusions for the code changes provided by the user. Please output the issues in the original code, your review suggestions, and modification proposals in your conclusion. Please keep the total output within 200 words", ext)
    } else {
        // 默认中文
        systemTpl = fmt.Sprintf("你是一位  %s 研发专家，现在你将对用户给出的代码变更内容给出对应的code reviewer 结论。我需要你在结论中输出原有代码相关问题，你的评审建议，与修改方案.请将整体输出控制在200字内", ext)
    }

    chatTpl := prompt.FromMessages(schema.FString,
        schema.SystemMessage(systemTpl),
        schema.MessagesPlaceholder("message_histories", true),
        schema.UserMessage("{user_query}"),
    )
    _ = g.AddChatTemplateNode(nodeOfPrompt, chatTpl)
    _ = g.AddChatModelNode(nodeOfModel, e.chatModel)
    _ = g.AddEdge(compose.START, nodeOfPrompt)
    _ = g.AddEdge(nodeOfPrompt, nodeOfModel)
    _ = g.AddEdge(nodeOfModel, compose.END)
    r, err := g.Compile(e.ctx, compose.WithMaxRunSteps(10))
    if err != nil {
        // 不再 panic，返回错误
        return fmt.Errorf("compile graph failed: %w", err)
    }

    ret, err := r.Invoke(e.ctx, map[string]any{
        "message_histories": []*schema.Message{},
        "user_query":        d.Content,
    })
    if err != nil {
        return fmt.Errorf("invoke failed: %w", err)
    }

    e.mutex.Lock()
    defer e.mutex.Unlock()

    lang := getFileLanguage(d.FilePath)
    if err := e.writeReviewToFile(d.FilePath, ret, lang); err != nil {
        return fmt.Errorf("write review failed: %w", err)
    }
    duration := time.Since(start)
    color.Green("✔ reviewed: %s in %v\n", d.FilePath, duration)
    return nil
}

const (
    nodeOfModel  = "model"
    nodeOfPrompt = "prompt"
)
