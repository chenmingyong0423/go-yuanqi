package yuanqi

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

const (
	url                   = "https://open.hunyuan.tencent.com/openapi/v1/agent/chat/completions"
	HeaderXSource         = "X-Source"
	HeaderXSourceValue    = "openapi"
	HeaderAuthorization   = "Authorization"
	HeaderContentType     = "Content-Type"
	HeaderApplicationJson = "application/json"
)

type Session struct {
	AssistantId      string        `json:"assistant_id"`      // 助手ID
	UserId           string        `json:"user_id"`           // 用户ID，调用者业务侧的用户ID，会影响智能体的数据统计，建议按实际情况填写
	Token            string        `json:"-"`                 // Token
	AssistantVersion string        `json:"version,omitempty"` // 助手版本 (仅对内部开放)，可选
	Timeout          time.Duration `json:"-"`                 // 会话请求超时时间，可选

	// 是否以流式接口的形式返回数据，默认false，可选
	Stream bool `json:"stream"`
	// 默认为published，传preview时，表示使用草稿态智能体 (仅对内部开放)，可选
	ChatType string `json:"chat_type"`
	// 会话内容, 长度最多为40, 按对话时间从旧到新在数组中排列
	Messages []Message `json:"messages"`

	// 腾讯元器 API URL
	url string
}

func NewSession(completion *Chat) *Session {
	return &Session{
		AssistantId:      completion.AssistantId,
		UserId:           completion.UserId,
		Token:            completion.Token,
		AssistantVersion: completion.AssistantVersion,
		Timeout:          completion.Timeout,
		url:              url,
	}
}

func (c *Session) Request(ctx context.Context) (*SessionResponse, error) {
	if c.Stream {
		return nil, fmt.Errorf("stream request not supported")
	}
	body, err := jsoniter.Marshal(c)
	if err != nil {
		return nil, err
	}
	resp := new(SessionResponse)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add(HeaderXSource, HeaderXSourceValue)
	req.Header.Add(HeaderContentType, HeaderApplicationJson)
	req.Header.Add(HeaderAuthorization, fmt.Sprintf("Bearer %s", c.Token))
	client := http.DefaultClient
	if c.Timeout != 0 {
		client.Timeout = c.Timeout
	}
	httpResp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	data, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return resp, err
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, &HttpErrorResponse{
			Status:     httpResp.Status,
			StatusCode: httpResp.StatusCode,
			Body:       data,
		}
	}
	if err = jsoniter.Unmarshal(data, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}

func (c *Session) StreamRequest(ctx context.Context) (<-chan *SessionResponse, <-chan error) {
	// 创建返回的通道
	respChan := make(chan *SessionResponse)
	errChan := make(chan error)

	go func() {
		defer close(respChan)
		defer close(errChan)

		body, err := jsoniter.Marshal(c)
		if err != nil {
			errChan <- err
			return
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(body))
		if err != nil {
			errChan <- err
			return
		}
		req.Header.Add(HeaderXSource, HeaderXSourceValue)
		req.Header.Add(HeaderContentType, HeaderApplicationJson)
		req.Header.Add(HeaderAuthorization, fmt.Sprintf("Bearer %s", c.Token))

		client := http.DefaultClient
		if c.Timeout != 0 {
			client.Timeout = c.Timeout
		}

		httpResp, err := client.Do(req)
		if err != nil {
			errChan <- err
			return
		}
		defer httpResp.Body.Close()

		if httpResp.StatusCode != http.StatusOK {
			data, _ := io.ReadAll(httpResp.Body)
			errChan <- &HttpErrorResponse{
				Status:     httpResp.Status,
				StatusCode: httpResp.StatusCode,
				Body:       data,
			}
			return
		}

		scanner := bufio.NewScanner(httpResp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "data:") {
				data := strings.TrimSpace(line[len("data:"):])
				if data == "[DONE]" {
					break
				}
				var resp SessionResponse
				if err = jsoniter.UnmarshalFromString(data, &resp); err != nil {
					errChan <- err
					return
				}
				respChan <- &resp
			}
		}

		if err = scanner.Err(); err != nil {
			errChan <- err
			return
		}
	}()

	return respChan, errChan
}

func (c *Session) WithStream(stream bool) *Session {
	c.Stream = stream
	return c
}

func (c *Session) WithChatType(chatType string) *Session {
	c.ChatType = chatType
	return c
}

func (c *Session) AddMessages(messages ...Message) *Session {
	c.Messages = append(c.Messages, messages...)
	return c
}

func (c *Session) WithUrl(stream bool) *Session {
	c.Stream = stream
	return c
}
