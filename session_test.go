package yuanqi

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSession_Request(t *testing.T) {
	testCases := []struct {
		name        string
		ctx         context.Context
		assistantId string
		userId      string
		token       string

		stream  bool
		role    Role
		content Content
		text    string

		wantError require.ErrorAssertionFunc
	}{
		{
			name:        "invalid request params for assistantId",
			ctx:         context.Background(),
			assistantId: "",
			userId:      os.Getenv("YUANQI_USER_ID"),
			token:       os.Getenv("YUANQI_TOKEN"),
			wantError: func(t require.TestingT, err error, i ...interface{}) {
				resp := &HttpErrorResponse{}
				require.True(t, errors.As(err, &resp))
			},
		},
		{
			name:        "invalid request params for userId",
			ctx:         context.Background(),
			assistantId: os.Getenv("YUANQI_USER_ID"),
			userId:      "",
			token:       os.Getenv("YUANQI_TOKEN"),
			wantError: func(t require.TestingT, err error, i ...interface{}) {
				resp := &HttpErrorResponse{}
				require.True(t, errors.As(err, &resp))
			},
		},
		{
			name:        "invalid request params for token",
			ctx:         context.Background(),
			assistantId: os.Getenv("YUANQI_USER_ID"),
			userId:      os.Getenv("YUANQI_USER_ID"),
			token:       "",
			wantError: func(t require.TestingT, err error, i ...interface{}) {
				resp := &HttpErrorResponse{}
				require.True(t, errors.As(err, &resp))
			},
		},
		//{
		//	name:        "invalid request params for message(role)",
		//	ctx:         context.Background(),
		//	assistantId: os.Getenv("YUANQI_ASSISTANT_ID"),
		//	userId:      os.Getenv("YUANQI_USER_ID"),
		//	token:       os.Getenv("YUANQI_TOKEN"),
		//	stream:      false,
		//	content:     NewContentBuilder().Text("你好").Build(),
		//	wantError: func(t require.TestingT, err error, i ...interface{}) {
		//		resp := &HttpErrorResponse{}
		//		require.True(t, errors.As(err, &resp))
		//	},
		//},
		//{
		//	name:        "invalid request params for content",
		//	ctx:         context.Background(),
		//	assistantId: os.Getenv("YUANQI_ASSISTANT_ID"),
		//	userId:      os.Getenv("YUANQI_USER_ID"),
		//	token:       os.Getenv("YUANQI_TOKEN"),
		//	stream:      false,
		//	role:        RoleUser,
		//	wantError: func(t require.TestingT, err error, i ...interface{}) {
		//		resp := &HttpErrorResponse{}
		//		require.True(t, errors.As(err, &resp))
		//	},
		//},
		{
			name:        "success",
			ctx:         context.Background(),
			assistantId: os.Getenv("YUANQI_ASSISTANT_ID"),
			userId:      os.Getenv("YUANQI_USER_ID"),
			token:       os.Getenv("YUANQI_TOKEN"),
			stream:      false,
			role:        RoleUser,
			content:     NewContentBuilder().Text("你好").Build(),
			wantError:   require.NoError,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建一个聊天对象
			chat := NewChat(tc.assistantId, tc.userId, tc.token)

			// 创建新的会话对象并设置会话流和类型
			session := chat.Session().WithStream(tc.stream)

			// 创建消息
			messageBuilder := NewMessageBuilder()
			if tc.role != "" {
				messageBuilder.Role(tc.role)
			}
			if tc.content.Type != "" {
				messageBuilder.Content(tc.content)
			}
			// 添加消息并发送以及处理错误
			resp, err := session.AddMessages(messageBuilder.Build()).Request(context.Background())
			tc.wantError(t, err)
			if err == nil {
				require.NotZero(t, resp.ID)
			}
		})
	}
}

func TestSession_StreamRequest(t *testing.T) {
	{
		// 400 参数错误
		// 创建一个聊天对象
		chat := NewChat("", "", "")
		session := chat.Session().WithStream(true)
		respChan, errChan := session.StreamRequest(context.Background())
		for {
			select {
			case resp, ok := <-respChan:
				require.False(t, ok)
				require.Nil(t, resp)
				respChan = nil
			case err, ok := <-errChan:
				if !ok {
					errChan = nil
				} else {
					resp := &HttpErrorResponse{}
					require.True(t, errors.As(err, &resp))
				}
			}
			if respChan == nil && errChan == nil {
				break
			}
		}
	}
	{
		// 正常请求
		assistantId := os.Getenv("YUANQI_ASSISTANT_ID")

		userId := os.Getenv("YUANQI_USER_ID")

		token := os.Getenv("YUANQI_TOKEN")

		// 创建一个聊天对象
		chat := NewChat(assistantId, userId, token)

		// 创建新的会话对象并设置会话流和类型
		session := chat.Session().WithStream(true)

		// 创建消息内容
		textContent := NewContentBuilder().Text("text").Build()
		//imageContent := NewContentBuilder().FileUrl(NewFileBuilder().Type("image").Url("https://domain/1.jpg").Build()).Build()
		// 创建消息
		message := NewMessageBuilder().
			Role("user").
			Content(textContent).Build()
		// 添加消息并发送以及处理错误
		respChan, errChan := session.AddMessages(message).StreamRequest(context.Background())
		for {
			select {
			case resp, ok := <-respChan:
				if !ok {
					respChan = nil
				} else {
					require.NotNil(t, resp)
					require.NotZero(t, resp.ID)
				}
			case err, ok := <-errChan:
				require.False(t, ok)
				require.Nil(t, err)
				errChan = nil
			}
			if respChan == nil && errChan == nil {
				break
			}
		}
	}
}
