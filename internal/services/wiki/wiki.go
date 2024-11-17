package wiki

import (
	"github.com/bochkov/gobot/internal/services"
)

type ThisDay struct {
	Date     string `json:"date"`
	WorldDay string `json:"worldday,omitempty"`
	Text     string `json:"text"`
	ImgSrc   string `json:"img,omitempty"`
}

type Service interface {
	services.Service
	Today() (*ThisDay, error)
}
