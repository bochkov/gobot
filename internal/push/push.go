package push

type Push interface {
	PushText() string
}

type PushService interface {
	Push(text string)
}
