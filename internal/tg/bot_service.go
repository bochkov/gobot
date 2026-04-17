package tg

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"reflect"

	"log/slog"

	"github.com/bochkov/gobot/internal/lib/db"
	"github.com/bochkov/gobot/internal/repo"
	"github.com/bochkov/gobot/internal/util"
	"github.com/carlmjohnson/requests"
)

type service struct {
	tmpMsgDao    repo.TmpMsgDao
	adapters     []TgAnswerAdapter
	push         TgPushAdapter
	temporaryMsg bool
}

func NewAnswerService(adapters ...TgAnswerAdapter) Service {
	return &service{adapters: adapters, temporaryMsg: false}
}

func NewPushService(push TgPushAdapter) Service {
	return &service{push: push, temporaryMsg: false}
}

func NewPushTmpService(push TgPushAdapter, tmpMsgDao repo.TmpMsgDao) Service {
	return &service{push: push, temporaryMsg: true, tmpMsgDao: tmpMsgDao}
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

	apiUrl := fmt.Sprintf("https://api.telegram.org/bot%s/%s", token, methodName)
	reqBuilder := requests.
		URL(apiUrl).
		Method(http.MethodPost)

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

	res := &TypedResult[any]{Result: methodResponse}
	if err := reqBuilder.
		ToJSON(res).
		Fetch(context.Background()); err != nil {
		return nil, err
	}
	return res, nil
}

func (s service) Exec(method Method, token string) (bool, error) {
	methodName, _ := method.Describe()
	slog.Info("request", "method", methodName, "body", util.ToJson(method))
	apiUrl := fmt.Sprintf("https://api.telegram.org/bot%s/%s", token, methodName)
	if err := requests.
		URL(apiUrl).
		Method(http.MethodPost).
		BodyJSON(&method).
		Fetch(context.Background()); err != nil {
		return false, err
	}
	return true, nil
}

func (s service) Push(recepients []string) {
	token := db.GetProp(db.TgBotTokenKey, "")
	if token == "" {
		slog.Warn("no token specified")
		return
	}

	sm, err := s.push.PushData(recepients)
	if err != nil {
		slog.Warn("cannot invoke push", "err", err.Error())
		return
	}
	for _, m := range sm {
		res, exec := s.Execute(m, token)
		if exec != nil {
			slog.Warn(exec.Error())
		}
		if res.Ok && s.temporaryMsg {
			msg := res.Result.(*Message)
			s.scheduleToRemove(msg.Chat.Id, msg.Id)
		}
	}
}

func (s service) scheduleToRemove(chatId int64, msgId int64) {
	slog.Info("scheduled to remove", "chatId", chatId, "msgId", msgId)
	s.tmpMsgDao.Save(context.Background(), chatId, msgId)
}
