package tg

import "github.com/bochkov/gobot/internal/push"

type MethodCustomize func(message *SendMessage[string])

type Worker interface {
	Description() string
	IsMatch(text string) bool
	Answer(msg *Message) []Method
}

type Service interface {
	push.PushService
	Execute(method Method, token string) (*TypedResult[any], error)
	GetAnswers(msg *Message) []Method
}
