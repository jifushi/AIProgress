package aichat

import "context"

func NewAIChat(ctx context.Context, userName string, sessionID string, modelType string, config map[string]interface{}) (*AIHelper, error) {
	// 创建新的AIHelper
	factory := GetGlobalFactory()
	helper, err := factory.CreateAIHelper(ctx, modelType, sessionID, config)
	if err != nil {
		return nil, err
	}
	return helper, nil
}
