package push

type Push interface {
	PushText() string
}

type Service interface {
	Push(text string)
}
