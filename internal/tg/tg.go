package tg

import (
	"github.com/bochkov/gobot/internal/push"
)

type Service interface {
	push.Service
	Execute(method Method, token string) (*TypedResult[any], error)
	GetAnswers(chatId int64, txt string) []Method
}

// интерфейс для отправки сообщений через телеграм-бот
type TgPushAdapter interface {
	PushData(receivers []string) ([]Method, error)
}

// интерфейс для ответов пользователю через телеграм-бот
type TgAnswerAdapter interface {
	Description() string
	IsMatch(text string) bool
	Answer(chatId int64, txt string) []Method
}
