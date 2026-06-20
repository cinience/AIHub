package service

import (
	"github.com/godeps/newapi/service/openaicompat"
	"github.com/godeps/newapi/setting/model_setting"
)

func ShouldChatCompletionsUseResponsesPolicy(policy model_setting.ChatCompletionsToResponsesPolicy, channelID int, channelType int, model string) bool {
	return openaicompat.ShouldChatCompletionsUseResponsesPolicy(policy, channelID, channelType, model)
}

func ShouldChatCompletionsUseResponsesGlobal(channelID int, channelType int, model string) bool {
	return openaicompat.ShouldChatCompletionsUseResponsesGlobal(channelID, channelType, model)
}
