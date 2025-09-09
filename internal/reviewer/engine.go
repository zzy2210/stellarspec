package reviewer

import (
    "context"
    "fmt"
    config "stellarspec/internal/model/conf"
    "sync"

    "github.com/cloudwego/eino-ext/components/model/openai"
    "github.com/fatih/color"
)

// EngineConfig 承载从 CLI 映射的参数（部分暂不启用）
type EngineConfig struct {
    ReviewPath    string
    MaxWorkers    int
    CommitID      string
    PromptPath    string
    ThinkingChain bool
    OutputFile    string
    Language      string
}

// Engine 负责编排：拉取变更 -> 并发审查 -> 写报告
type Engine struct {
    ctx context.Context

    cfg       EngineConfig
    chatModel *openai.ChatModel // 模型客户端

    // 文件写入互斥
    mutex sync.Mutex
}

func NewEngine(ctx context.Context, cfg EngineConfig) *Engine {
    return &Engine{ctx: ctx, cfg: cfg}
}

// CreateModel 根据基础配置创建模型客户端
func (e *Engine) CreateModel(conf *config.BaseConfig) error {
    if conf == nil {
        return fmt.Errorf("nil model config")
    }
    cm, err := newChatModel(e.ctx, conf)
    if err != nil {
        return err
    }
    e.chatModel = cm
    return nil
}

// Run 执行审查流程（返回错误而非 panic）
func (e *Engine) Run() error {
    diffs, err := e.gitDiff()
    if err != nil {
        return fmt.Errorf("get git diff failed: %w", err)
    }

    // 为保持行为一致，仍使用默认 10 并发；暂不启用 MaxWorkers
    maxWorkers := 10
    semaphore := make(chan struct{}, maxWorkers)

    var wg sync.WaitGroup
    for _, diff := range diffs {
        wg.Add(1)
        d := diff
        go func() {
            defer wg.Done()
            semaphore <- struct{}{}
            defer func() { <-semaphore }()
            if err := e.reviewSingleFile(d); err != nil {
                // 彩色错误输出，但不中断其他任务
                color.Red("✖ review failed: %s, err=%v\n", d.FilePath, err)
            }
        }()
    }
    wg.Wait()
    return nil
}
