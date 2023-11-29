package tg

import (
	"fmt"
	"io"
	"mime/multipart"

	"github.com/bochkov/gobot/internal/util"
)

type Method interface {
	Describe() (string, any)
}

type MultipartMethod interface {
	Method
	InputFileInfo() (string, InputFile)
}

type IntOrString interface {
	string | int64
}

type ParseMode string

const (
	HTML       ParseMode = "HTML"
	MARKDOWNV2 ParseMode = "MarkdownV2"
	MARKDOWN   ParseMode = "Markdown"
)

type SendOptions struct {
	MessageThreadId       string    `json:"message_thread_id,omitempty"`
	ParseMode             ParseMode `json:"parse_mode,omitempty"`
	DisableWebPagePreview bool      `json:"disable_web_page_preview,omitempty"`
	DisableNotification   bool      `json:"disable_notification,omitempty"`
}

type InputFile struct {
	bytes    []byte
	filename string
}

type SendMessage[T IntOrString] struct {
	ChatId T      `json:"chat_id"`
	Text   string `json:"text"`
	SendOptions
}

func (sm *SendMessage[T]) Describe() (string, any) {
	return "sendMessage", &Message{}
}

type SendDocument[T IntOrString] struct {
	ChatId                      T         `json:"chat_id,string"`
	InputFile                   InputFile `json:"-"`
	Caption                     string    `json:"caption,omitempty"`
	DisableContentTypeDetection bool      `json:"disable_content_type_detection,omitempty"`
	SendOptions
}

func (sd *SendDocument[T]) Describe() (string, any) {
	return "sendDocument", &Message{}
}

func (sd *SendDocument[T]) InputFileInfo() (string, InputFile) {
	return "document", sd.InputFile
}

func ToMultipart(method MultipartMethod, w io.Writer) (string, error) {
	// back and forth
	js := util.ToJson(method)
	data, err := util.FromJsonString(js)
	if err != nil {
		return "", err
	}

	mp := multipart.NewWriter(w)
	// write text fields
	for key, val := range data {
		if er := mp.WriteField(key, fmt.Sprintf("%v", val)); er != nil {
			return "", er
		}
	}
	// write binary field
	name, field := method.InputFileInfo()
	x, _ := mp.CreateFormFile(name, field.filename)
	if _, er := x.Write(field.bytes); er != nil {
		return "", er
	}
	if er := mp.Close(); er != nil {
		return "", err
	}
	return mp.FormDataContentType(), nil
}
