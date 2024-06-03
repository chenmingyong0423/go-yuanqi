package yuanqi

import (
	"fmt"
)

type SessionResponse struct {
	ID          string   `json:"id"`           // 此次请求的 id
	Created     int64    `json:"created"`      // unix时间戳
	Choices     []Choice `json:"choices"`      // 返回的回复, 当前仅有一个
	AssistantID string   `json:"assistant_id"` // 实际使用的助手 id
	Usage       Usage    `json:"usage"`        // token 使用量
}

type Choice struct {
	Index        int           `json:"index"`         // 第几个回复
	FinishReason string        `json:"finish_reason"` // "stop" 表示正常结束, "sensitive" 表示审核不通过, "tool_fail" 表示调用工具失败
	Message      ChoiceMessage `json:"message"`       // 返回的内容
	Delta        Delta         `json:"delta"`         // 返回的内容（流式返回）
}

type ChoiceMessage struct {
	Role    string `json:"role"`    // 角色名称
	Content string `json:"content"` // 内容详情
	Steps   []Step `json:"steps"`   // 助手的执行步骤
}

type Step struct {
	Role       string     `json:"role"`         // 执行步骤中的角色名称，assistant 表示模型，tool 表示工具调用
	Content    string     `json:"content"`      // 执行步骤的结果，当角色为 assistant 时表示模型的输出内容，当角色为 tool 时表示工具的输出内容
	ToolCallID string     `json:"tool_call_id"` // 角色为 tool 时有效，内容为模型生成的工具调用中的唯一 ID
	ToolCalls  []ToolCall `json:"tool_calls"`   // 模型生成的工具调用
	Usage      Usage      `json:"usage"`        // 当前执行步骤的 token 使用量
	TimeCost   int        `json:"time_cost"`    // 当前执行步骤的耗时
}

type ToolCall struct {
	ID       string   `json:"id"`       // 工具调用的唯一 ID
	Type     string   `json:"type"`     // 调用的工具类型，当前只支持 function
	Function Function `json:"function"` // 具体调用的 function
}

type Function struct {
	Name      string `json:"name"`      // function 名称
	Desc      string `json:"desc"`      // function 描述
	Type      string `json:"type"`      // function 类型，当前支持 tool/knowledge/workflow
	Arguments string `json:"arguments"` // 调用 function 的参数，JSON 格式
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`     // 问题 token 使用量
	CompletionTokens int `json:"completion_tokens"` // 回答 token 使用量
	TotalTokens      int `json:"total_tokens"`      // token 总使用量
}

type Delta struct {
	Role       string     `json:"role"`         // 角色名称，assistant 表示模型，tool 表示工具调用（流式返回）
	Content    string     `json:"content"`      // 内容详情，当角色为 assistant 时表示模型的输出内容，当角色为 tool 时表示工具的输出内容（流式返回）
	ToolCallID string     `json:"tool_call_id"` // 角色为 tool 时有效，内容为模型生成的工具调用中对应的 tool_call ID （流式返回）
	ToolCalls  []ToolCall `json:"tool_calls"`   // 模型生成的工具调用（流式返回）
	TimeCost   int        `json:"time_cost"`    // 当前执行步骤的耗时（流式返回）
}

type HttpErrorResponse struct {
	Status     string // e.g. "400 Bad Request"
	StatusCode int    `json:"status_code"` // http 状态码
	Body       []byte `json:"body"`        // http 响应体
}

func (h *HttpErrorResponse) Error() string {
	return fmt.Sprintf("response error: statusCode: %d, status: %s", h.StatusCode, h.Status)
}
