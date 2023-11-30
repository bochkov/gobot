package tg

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/bochkov/gobot/internal/lib/db"
	"github.com/bochkov/gobot/internal/util"
	"github.com/carlmjohnson/requests"
	"log/slog"
)

type service struct {
	workers []Worker
}

func NewService(workers ...Worker) Service {
	return &service{workers: workers}
}

func (s *service) GetAnswers(msg *Message) []Method {
	for _, serv := range s.workers {
		if serv.IsMatch(msg.Text) {
			slog.Info(fmt.Sprintf("choosed = %s", serv.Description()))
			return serv.Answer(msg)
		}
	}
	return nil
}

func (s *service) Execute(method Method, token string) (*TypedResult[any], error) {
	methodName, methodResponse := method.Describe()
	slog.Info("request", "method", methodName, "body", util.ToJson(method))
	res := &TypedResult[any]{
		Result: methodResponse,
	}
	reqBuilder := requests.URL(fmt.Sprintf("https://api.telegram.org/bot%s/%s", token, methodName))
	if reflect.ValueOf(method).Elem().FieldByName("InputFile").IsValid() {
		buf := &bytes.Buffer{}
		contentType, err := ToMultipart(method.(MultipartMethod), buf)
		if err != nil {
			return nil, err
		}
		reqBuilder.BodyBytes(buf.Bytes()).ContentType(contentType)
	} else {
		reqBuilder.BodyJSON(&method)
	}
	if err := reqBuilder.
		Method(http.MethodPost).
		ToJSON(res).
		Fetch(context.Background()); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *service) Push(text string) {
	if text == "" {
		slog.Warn("empty text")
		return
	}
	token := db.GetProp(db.TgBotTokenKey, "")
	if token == "" {
		slog.Warn("no token specified")
		return
	}
	chatId := db.GetProp(db.ChatAutoSend, "") // TODO
	if chatId == "" {
		slog.Warn("no chat.id specified")
		return
	}
	for _, chat := range strings.Split(chatId, ";") {
		sm := new(SendMessage[string])
		sm.ChatId = chat
		sm.Text = text
		sm.SendOptions.DisableWebPagePreview = true
		sm.SendOptions.DisableNotification = true
		if _, exec := s.Execute(sm, token); exec != nil {
			slog.Warn(exec.Error())
		}
	}
}
