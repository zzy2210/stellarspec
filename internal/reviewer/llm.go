package reviewer

import (
    "context"
    config "stellarspec/internal/model/conf"

    "github.com/cloudwego/eino-ext/components/model/openai"
)

// newChatModel 创建底层 OpenAI ChatModel
func newChatModel(ctx context.Context, conf *config.BaseConfig) (*openai.ChatModel, error) {
    modelConf := &openai.ChatModelConfig{
        APIKey:  conf.Key,
        BaseURL: conf.APIServer,
        Model:   conf.Model,
    }
    return openai.NewChatModel(ctx, modelConf)
}
