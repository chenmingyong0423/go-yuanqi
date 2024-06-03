<h1 align="center">
  go-yuanqi
</h1>

[![GitHub Repo stars](https://img.shields.io/github/stars/chenmingyong0423/go-yuanqi)](https://github.com/chenmingyong0423/go-yuanqi/stargazers)
[![GitHub issues](https://img.shields.io/github/issues/chenmingyong0423/go-yuanqi)](https://github.com/chenmingyong0423/go-yuanqi/issues)
[![GitHub License](https://img.shields.io/github/license/chenmingyong0423/go-yuanqi)](https://github.com/chenmingyong0423/go-yuanqi/blob/main/LICENSE)
[![GitHub release (with filter)](https://img.shields.io/github/v/release/chenmingyong0423/go-yuanqi)](https://github.com/chenmingyong0423/go-yuanqi)
[![Go Report Card](https://goreportcard.com/badge/github.com/chenmingyong0423/go-yuanqi)](https://goreportcard.com/report/github.com/chenmingyong0423/go-yuanqi)
[![All Contributors](https://img.shields.io/badge/all_contributors-1-orange.svg?style=flat-square)](#contributors-)

`go-yuanqi` 库是一个用于简化腾讯元器 `API` 调用的库，通过这个库，开发者可以更高效，更简洁地与腾讯元器 `API` 交互，减少重复代码，提高开发效率。

## 功能
- **链式调用**
    - 通过链式调用的方式封装请求参数和调用接口，使代码更加简洁和可读。

- **非流式 API 交互**
    - 适用于一次性获取数据的场景。例如当 `stream` 参数指定为 `false` 的场景。

- **流式 API 交互**
    - 支持处理流式响应，例如当 `stream` 参数被指定为 `true` 的场景。

# 入门指南
## 安装
```shell
go get github.com/chenmingyong0423/go-mongox
```

## 使用
### 非流式 API 交互
```go
// 创建一个聊天对象
chat := yuanqi.NewChat("assistantId", "userId", "token", yuanqi.WithAssistantVersion(""), yuanqi.WithTimeOut(10*time.Second))

// 创建新的会话对象并设置会话流和类型
session := chat.Chat().WithStream(false).WithChatType("published")

// 创建消息内容
// - 文字消息
textContent := yuanqi.NewContentBuilder().Text("你好").Build()
// 图片消息
imageContent := yuanqi.NewContentBuilder().FileUrl(yuanqi.NewFileBuilder().Type("image").Url("https://domain/1.jpg").Build()).Build()
// 创建消息
message := yuanqi.NewMessageBuilder().
    Role("user").
    Content(textContent, imageContent).Build()

// 添加消息并发送以及处理错误
resp, err := session.AddMessages(message).Request(context.Background())
```
非流式 API 交互需要调用 `Request` 方法，该方法会返回一个 `SessionResponse` 对象和一个 `error` 对象。
### 流式 API 交互
```go
// 创建一个聊天对象
chat := yuanqi.NewChat("assistantId", "userId", "token", yuanqi.WithAssistantVersion(""), yuanqi.WithTimeOut(10*time.Second))

// 创建新的会话对象并设置会话流和类型
session := chat.Chat().WithStream(true).WithChatType("published")

// 创建消息内容
// - 文字消息
textContent := yuanqi.NewContentBuilder().Text("你好").Build()
// 图片消息
imageContent := yuanqi.NewContentBuilder().FileUrl(yuanqi.NewFileBuilder().Type("image").Url("https://domain/1.jpg").Build()).Build()
// 创建消息
message := yuanqi.NewMessageBuilder().
    Role("user").
    Content(textContent, imageContent).Build()

// 添加消息并发送以及处理错误
respChan, errChan := session.AddMessages(message).StreamRequest(context.Background())
for {
    select {
    case resp, ok := <-respChan:
        if !ok {
            respChan = nil
        } else {
            fmt.Println(resp)
        }
    case err, ok := <-errChan:
        if !ok {
            errChan = nil
        } else {
            panic(err)
        }
    }
    if respChan == nil && errChan == nil {
        break
    }
}
```
流式 API 交互需要调用 `StreamRequest` 方法，该方法会返回一个 `chan SessionResponse` 对象和一个 `chan error` 对象。