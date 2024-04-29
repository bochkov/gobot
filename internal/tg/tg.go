package tg

import "github.com/bochkov/gobot/internal/push"

type MethodCustomize func(message *SendMessage[string])

type Worker interface {
	Description() string
	IsMatch(text string) bool
	Answer(chatId int64, txt string) []Method
}

type Service interface {
	push.Service
	Execute(method Method, token string) (*TypedResult[any], error)
	GetAnswers(chatId int64, txt string) []Method
}
