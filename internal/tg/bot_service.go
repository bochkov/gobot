package tg

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"log/slog"

	"github.com/bochkov/gobot/internal/lib/db"
	"github.com/bochkov/gobot/internal/util"
	"github.com/carlmjohnson/requests"
)

type service struct {
	adapters []TgAnswerAdapter
	push     TgPushAdapter
}

func NewAnswerService(adapters ...TgAnswerAdapter) Service {
	return &service{adapters: adapters}
}

func NewPushService(push TgPushAdapter) Service {
	return &service{push: push}
}

func (s service) GetAnswers(chatId int64, txt string) []Method {
	for _, serv := range s.adapters {
		if serv.IsMatch(txt) {
			slog.Info(fmt.Sprintf("choosed = %s", serv.Description()))
			return serv.Answer(chatId, txt)
		}
	}
	return nil
}

func (s service) Execute(method Method, token string) (*TypedResult[any], error) {
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

func (s service) Push() {
	token := db.GetProp(db.TgBotTokenKey, "")
	if token == "" {
		slog.Warn("no token specified")
		return
	}
	chatId := db.GetProp(db.ChatAutoSend, "")
	if chatId == "" {
		slog.Warn("no chat.id specified")
		return
	}
	receivers := strings.Split(chatId, ";")
	sm, err := s.push.PushData(receivers)
	if err != nil {
		slog.Warn("cannot invoke push", "err", err.Error())
		return
	}
	for _, m := range sm {
		if _, exec := s.Execute(m, token); exec != nil {
			slog.Warn(exec.Error())
		}
	}
}
