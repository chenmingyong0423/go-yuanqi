// Copyright 2024 chenmingyong0423

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package yuanqi

type (
	Role        string
	ContentType string
)

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"

	ContentTypeText    ContentType = "text"
	ContentTypeFileURL ContentType = "file_url"
)

type Message struct {
	// 角色, 'user' 或者 'assistant', 在 message 中必须是 user 与 assistant 交替(一问一答)
	Role Role `json:"role"`
	// 可以传入多种类型的内容，如图片、文件或文本
	Content []Content `json:"content"`
}

type MessageBuilder struct {
	message Message
}

func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{}
}

func (b *MessageBuilder) Role(role Role) *MessageBuilder {
	b.message.Role = role
	return b
}

func (b *MessageBuilder) Content(content ...Content) *MessageBuilder {
	b.message.Content = append(b.message.Content, content...)
	return b
}

func (b *MessageBuilder) Build() Message {
	return b.message
}

type Content struct {
	// 内容的类型，可选参数为 'text' 或 'file_url'，可选
	Type ContentType `json:"type,omitempty"`
	// 当 type 为 text 时使用，表示具体的文本内容，可选
	Text string `json:"text,omitempty"`
	// 当 type 为 file_url 时使用，表示具体的文件内容，可选
	FileUrl FileUrl `json:"file_url,omitempty"`
}

type ContentBuilder struct {
	content Content
}

func NewContentBuilder() *ContentBuilder {
	return &ContentBuilder{}
}

func (b *ContentBuilder) Text(text string) *ContentBuilder {
	b.content.Type = ContentTypeText
	b.content.Text = text
	return b
}

func (b *ContentBuilder) FileUrl(fileUrl FileUrl) *ContentBuilder {
	b.content.Type = ContentTypeFileURL
	b.content.FileUrl = fileUrl
	return b
}

func (b *ContentBuilder) Build() Content {
	return b.content
}

type FileUrl struct {
	// 文件的类型，例如image/video/audio/pdf/doc/txt等，可选
	Type string `json:"type,omitempty"`
	// 文件的 url，可选
	Url string `json:"url"`
}

type FileBuilder struct {
	fileUrl FileUrl
}

func NewFileBuilder() *FileBuilder {
	return &FileBuilder{}
}

func (b *FileBuilder) Type(fileType string) *FileBuilder {
	b.fileUrl.Type = fileType
	return b
}

func (b *FileBuilder) Url(url string) *FileBuilder {
	b.fileUrl.Url = url
	return b
}

func (b *FileBuilder) Build() FileUrl {
	return b.fileUrl
}
