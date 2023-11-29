package anekdot

import (
	"github.com/bochkov/gobot/internal/push"
)

type Anekdot struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
}

type Service interface {
	push.Push
	GetRandom() (*Anekdot, error)
	GetById(id int) (*Anekdot, error)
	Description() string
}
