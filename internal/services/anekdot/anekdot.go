package anekdot

import (
	"github.com/bochkov/gobot/internal/services"
)

type Anekdot struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
}

type Service interface {
	services.Service
	GetRandom() (*Anekdot, error)
	GetById(id int) (*Anekdot, error)
}
